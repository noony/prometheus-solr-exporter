package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const (
	mbeansPath     = "/admin/mbeans?stats=true&wt=json&cat=CORE&cat=QUERYHANDLER&cat=UPDATEHANDLER&cat=CACHE"
	adminCoresPath = "/admin/cores?action=STATUS&wt=json"
)

var (
	gaugeAdminMetrics = map[string]string{
		"num_docs":      "num_docs",
		"size_in_bytes": "size_in_bytes",
		"deleted_docs":  "deleted_docs",
		"max_docs":      "max_docs",
	}
	gaugeCoreMetrics = map[string]string{
		"num_docs":     "num_docs",
		"deleted_docs": "deleted_docs",
		"max_docs":     "max_docs",
	}
	gaugeQueryMetrics = map[string]string{
		"15min_rate_reqs_per_second": "15min_rate_reqs_per_second",
		"5min_rate_reqs_per_second":  "5min_rate_reqs_per_second",
		"75th_pc_request_time":       "75th_pc_request_time",
		"95th_pc_request_time":       "95th_pc_request_time",
		"99th_pc_request_time":       "99th_pc_request_time",
		"999th_pc_request_time":      "999th_pc_request_time",
		"avg_requests_per_second":    "avg_requests_per_second",
		"avg_time_per_request":       "avg_time_per_request",
		"errors":                     "errors",
		"handler_start":              "handler_start",
		"median_request_time":        "median_request_time",
		"requests":                   "requests",
		"timeouts":                   "timeouts",
		"total_time":                 "total_time",
	}
	gaugeUpdateMetrics = map[string]string{
		"adds":                        "adds",
		"autocommit_max_docs":         "autocommit_max_docs",
		"autocommit_max_time":         "autocommit_max_time",
		"autocommits":                 "autocommits",
		"commits":                     "commits",
		"cumulative_adds":             "cumulative_adds",
		"cumulative_deletes_by_id":    "cumulative_deletes_by_id",
		"cumulative_deletes_by_query": "cumulative_deletes_by_query",
		"cumulative_errors":           "cumulative_errors",
		"deletes_by_id":               "deletes_by_id",
		"deletes_by_query":            "deletes_by_query",
		"docs_pending":                "docs_pending",
		"errors":                      "errors",
		"expunge_deletes":             "expunge_deletes",
		"optimizes":                   "optimizes",
		"rollbacks":                   "rollbacks",
		"soft_autocommits":            "soft_autocommits",
	}
	gaugeCacheMetrics = map[string]string{
		"cumulative_evictions": "cumulative_evictions",
		"cumulative_hitratio":  "cumulative_hitratio",
		"cumulative_hits":      "cumulative_hits",
		"cumulative_inserts":   "cumulative_inserts",
		"cumulative_lookups":   "cumulative_lookups",
		"evictions":            "evictions",
		"hitratio":             "hitratio",
		"hits":                 "hits",
		"inserts":              "inserts",
		"lookups":              "lookups",
		"size":                 "size",
		"warmup_time":          "warmup_time",
	}
)

// Return list of cores from solr server
func getCoresFromStatus(adminCoresStatus *AdminCoresStatus) []string {
	serverCores := []string{}
	for coreName := range adminCoresStatus.Status {
		serverCores = append(serverCores, coreName)
	}
	return serverCores
}

// Exporter collects Solr stats from the given server and exports
// them using the prometheus metrics package.
type Exporter struct {
	MbeansUrl    string
	AdminCoreUrl string
	mutex        sync.RWMutex

	up prometheus.Gauge

	gaugeAdmin  map[string]*prometheus.GaugeVec
	gaugeCore   map[string]*prometheus.GaugeVec
	gaugeQuery  map[string]*prometheus.GaugeVec
	gaugeUpdate map[string]*prometheus.GaugeVec
	gaugeCache  map[string]*prometheus.GaugeVec

	client *http.Client
}

// NewExporter returns an initialized Exporter.
func NewExporter(solrURI string, solrContextPath string, timeout time.Duration) *Exporter {
	gaugeAdmin := make(map[string]*prometheus.GaugeVec, len(gaugeAdminMetrics))
	gaugeCore := make(map[string]*prometheus.GaugeVec, len(gaugeCoreMetrics))
	gaugeQuery := make(map[string]*prometheus.GaugeVec, len(gaugeQueryMetrics))
	gaugeUpdate := make(map[string]*prometheus.GaugeVec, len(gaugeUpdateMetrics))
	gaugeCache := make(map[string]*prometheus.GaugeVec, len(gaugeCacheMetrics))

	for name, help := range gaugeAdminMetrics {
		gaugeAdmin[name] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace + "_admin",
			Name:      name,
			Help:      help,
		}, []string{"core"})
	}

	for name, help := range gaugeCoreMetrics {
		gaugeCore[name] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace + "_core",
			Name:      name,
			Help:      help,
		}, []string{"core", "handler", "class"})
	}

	for name, help := range gaugeQueryMetrics {
		gaugeQuery[name] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace + "_queryhandler",
			Name:      name,
			Help:      help,
		}, []string{"core", "handler", "class"})
	}

	for name, help := range gaugeUpdateMetrics {
		gaugeUpdate[name] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace + "_updatehandler",
			Name:      name,
			Help:      help,
		}, []string{"core", "handler", "class"})
	}
	for name, help := range gaugeCacheMetrics {
		gaugeCache[name] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace + "_cache",
			Name:      name,
			Help:      help,
		}, []string{"core", "handler", "class"})
	}

	mbeansUrl := fmt.Sprintf("%s%s/%s%s", solrURI, solrContextPath, "%s", mbeansPath)
	adminCoreUrl := fmt.Sprintf("%s%s%s", solrURI, solrContextPath, adminCoresPath)

	// Init our exporter.
	return &Exporter{
		MbeansUrl:    mbeansUrl,
		AdminCoreUrl: adminCoreUrl,

		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "up",
			Help:      "Was the Solr instance query successful?",
		}),

		gaugeAdmin:  gaugeAdmin,
		gaugeCore:   gaugeCore,
		gaugeQuery:  gaugeQuery,
		gaugeUpdate: gaugeUpdate,
		gaugeCache:  gaugeCache,

		client: &http.Client{
			Transport: &http.Transport{
				Dial: func(netw, addr string) (net.Conn, error) {
					c, err := net.DialTimeout(netw, addr, timeout)
					if err != nil {
						return nil, err
					}
					if err := c.SetDeadline(time.Now().Add(timeout)); err != nil {
						return nil, err
					}
					return c, nil
				},
			},
		},
	}
}

// Describe describes all the metrics ever exported by the solr
// exporter. It implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.up.Desc()

	for _, vec := range e.gaugeAdmin {
		vec.Describe(ch)
	}
	for _, vec := range e.gaugeCore {
		vec.Describe(ch)
	}
	for _, vec := range e.gaugeQuery {
		vec.Describe(ch)
	}
	for _, vec := range e.gaugeUpdate {
		vec.Describe(ch)
	}
	for _, vec := range e.gaugeCache {
		vec.Describe(ch)
	}
}

// Collect fetches the stats from configured solr location and
// delivers them as Prometheus metrics. It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock() // To protect metrics from concurrent collects.
	defer e.mutex.Unlock()

	// Reset metrics.
	for _, vec := range e.gaugeAdmin {
		vec.Reset()
	}
	for _, vec := range e.gaugeCore {
		vec.Reset()
	}
	for _, vec := range e.gaugeQuery {
		vec.Reset()
	}
	for _, vec := range e.gaugeUpdate {
		vec.Reset()
	}
	for _, vec := range e.gaugeCache {
		vec.Reset()
	}

	e.up.Set(0)
	defer func() { ch <- e.up }()

	resp, err := e.client.Get(e.AdminCoreUrl)
	if err != nil {
		log.Errorf("Error while querying Solr for admin stats: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Failed to read admin stats response body: %v", err)
		return
	}

	adminCoresStatus := &AdminCoresStatus{}
	err = json.Unmarshal(body, adminCoresStatus)
	if err != nil {
		log.Errorf("Failed to unmarshal solr admin JSON into struct: %v", err)
		return
	}

	for core, metrics := range adminCoresStatus.Status {
		e.gaugeAdmin["num_docs"].WithLabelValues(core).Set(float64(metrics.Index.NumDocs))
		e.gaugeAdmin["size_in_bytes"].WithLabelValues(core).Set(float64(metrics.Index.SizeInBytes))
		e.gaugeAdmin["deleted_docs"].WithLabelValues(core).Set(float64(metrics.Index.DeletedDocs))
		e.gaugeAdmin["max_docs"].WithLabelValues(core).Set(float64(metrics.Index.MaxDoc))
	}

	cores := getCoresFromStatus(adminCoresStatus)

	for _, coreName := range cores {
		mBeansUrl := fmt.Sprintf(e.MbeansUrl, coreName)
		resp, err := e.client.Get(mBeansUrl)
		if err != nil {
			log.Errorf("Error while querying Solr for mbeans stats: %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Errorf("solr: API responded with status-code %d, expected %d, url %s",
				resp.StatusCode, http.StatusOK, mBeansUrl)
			return
		}

		mBeansData := &MBeansData{}
		if err := json.NewDecoder(resp.Body).Decode(mBeansData); err != nil {
			log.Errorf("Failed to unmarshal mbeansdata JSON into struct: %v", err)
			return
		}

		var coreMetrics map[string]Core
		if err := json.Unmarshal(findMBeansData(mBeansData.SolrMbeans, "CORE"), &coreMetrics); err != nil {
			log.Errorf("Failed to unmarshal mbeans core metrics JSON into struct: %v", err)
			return
		}

		for name, metrics := range coreMetrics {
			if strings.Contains(name, "@") {
				continue
			}

			e.gaugeCore["deleted_docs"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.DeletedDocs))
			e.gaugeCore["max_docs"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.MaxDoc))
			e.gaugeCore["num_docs"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.NumDocs))
		}

		b := bytes.Replace(findMBeansData(mBeansData.SolrMbeans, "QUERY"), []byte(":\"NaN\""), []byte(":0.0"), -1)
		var queryMetrics map[string]QueryHandler
		if err := json.Unmarshal(b, &queryMetrics); err != nil {
			log.Errorf("Failed to unmarshal mbeans query metrics JSON into struct: %v", err)
			return
		}

		for name, metrics := range queryMetrics {
			if strings.Contains(name, "@") || strings.Contains(name, "/admin") || strings.Contains(name, "/debug/dump") || strings.Contains(name, "/schema") || strings.Contains(name, "org.apache.solr.handler.admin") {
				continue
			}

			var FiveminRateRequestsPerSecond, One5minRateRequestsPerSecond float64
			if metrics.Stats.One5minRateReqsPerSecond == nil && metrics.Stats.FiveMinRateReqsPerSecond == nil {
				FiveminRateRequestsPerSecond = float64(metrics.Stats.FiveminRateRequestsPerSecond)
				One5minRateRequestsPerSecond = float64(metrics.Stats.One5minRateRequestsPerSecond)
			} else {
				FiveminRateRequestsPerSecond = float64(*metrics.Stats.FiveMinRateReqsPerSecond)
				One5minRateRequestsPerSecond = float64(*metrics.Stats.One5minRateReqsPerSecond)
			}

			e.gaugeQuery["15min_rate_reqs_per_second"].WithLabelValues(coreName, name, metrics.Class).Set(One5minRateRequestsPerSecond)
			e.gaugeQuery["5min_rate_reqs_per_second"].WithLabelValues(coreName, name, metrics.Class).Set(FiveminRateRequestsPerSecond)
			e.gaugeQuery["75th_pc_request_time"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.Seven5thPcRequestTime))
			e.gaugeQuery["95th_pc_request_time"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.Nine5thPcRequestTime))
			e.gaugeQuery["99th_pc_request_time"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.Nine9thPcRequestTime))
			e.gaugeQuery["999th_pc_request_time"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.Nine99thPcRequestTime))
			e.gaugeQuery["avg_requests_per_second"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.AvgRequestsPerSecond))
			e.gaugeQuery["avg_time_per_request"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.AvgTimePerRequest))
			e.gaugeQuery["errors"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.Errors))
			e.gaugeQuery["handler_start"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.HandlerStart))
			e.gaugeQuery["median_request_time"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.MedianRequestTime))
			e.gaugeQuery["requests"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.Requests))
			e.gaugeQuery["timeouts"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.Timeouts))
			e.gaugeQuery["total_time"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.TotalTime))
		}

		var updateMetrics map[string]UpdateHandler
		if err := json.Unmarshal(findMBeansData(mBeansData.SolrMbeans, "UPDATE"), &updateMetrics); err != nil {
			log.Errorf("Failed to unmarshal mbeans update metrics JSON into struct: %v", err)
			return
		}

		for name, metrics := range updateMetrics {
			if strings.Contains(name, "@") || strings.HasPrefix(name, "/") {
				continue
			}
			var autoCommitMaxTime int
			if len(metrics.Stats.AutocommitMaxTime) > 2 {
				autoCommitMaxTime, _ = strconv.Atoi(metrics.Stats.AutocommitMaxTime[:len(metrics.Stats.AutocommitMaxTime)-2])
			}
			e.gaugeUpdate["adds"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.Adds))
			e.gaugeUpdate["autocommit_max_docs"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.AutocommitMaxDocs))
			e.gaugeUpdate["autocommit_max_time"].WithLabelValues(coreName, name, metrics.Class).Set(float64(autoCommitMaxTime))
			e.gaugeUpdate["autocommits"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.Autocommits))
			e.gaugeUpdate["commits"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.Commits))
			e.gaugeUpdate["cumulative_adds"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.CumulativeAdds))
			e.gaugeUpdate["cumulative_deletes_by_id"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.CumulativeDeletesByID))
			e.gaugeUpdate["cumulative_deletes_by_query"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.CumulativeDeletesByQuery))
			e.gaugeUpdate["cumulative_errors"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.CumulativeErrors))
			e.gaugeUpdate["deletes_by_id"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.DeletesByID))
			e.gaugeUpdate["deletes_by_query"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.DeletesByQuery))
			e.gaugeUpdate["docs_pending"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.DocsPending))
			e.gaugeUpdate["errors"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.Errors))
			e.gaugeUpdate["expunge_deletes"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.ExpungeDeletes))
			e.gaugeUpdate["optimizes"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.Optimizes))
			e.gaugeUpdate["rollbacks"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.Rollbacks))
			e.gaugeUpdate["soft_autocommits"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.SoftAutocommits))
		}

		// Try to decode solr > v5 cache metrics
		var cacheMetrics map[string]Cache
		if err := json.Unmarshal(findMBeansData(mBeansData.SolrMbeans, "CACHE"), &cacheMetrics); err != nil {
			var cacheMetricsSolrV4 map[string]CacheSolrV4
			// Try to decode solr v4 metrics
			if err := json.Unmarshal(findMBeansData(mBeansData.SolrMbeans, "CACHE"), &cacheMetricsSolrV4); err != nil {
				log.Errorf("Failed to unmarshal mbeans cache metrics JSON into struct (core : %s): %v", coreName, err)
				return
			} else {
				for name, metrics := range cacheMetricsSolrV4 {
					if metrics.Class == "org.apache.solr.search.SolrFieldCacheMBean" {
						continue
					}
					hitratio, err := strconv.ParseFloat(metrics.Stats.Hitratio, 64)
					if err != nil {
						log.Errorf("Fail to convert Hitratio in float")
					}
					cumulative_hitratio, err := strconv.ParseFloat(metrics.Stats.CumulativeHitratio, 64)
					if err != nil {
						log.Errorf("Fail to convert Cumulative Hitratio in float")
					}
					e.gaugeCache["cumulative_evictions"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.CumulativeEvictions))
					e.gaugeCache["cumulative_hitratio"].WithLabelValues(coreName, name, metrics.Class).Set(cumulative_hitratio)
					e.gaugeCache["cumulative_hits"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.CumulativeHits))
					e.gaugeCache["cumulative_inserts"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.CumulativeInserts))
					e.gaugeCache["cumulative_lookups"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.CumulativeLookups))
					e.gaugeCache["evictions"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.Evictions))
					e.gaugeCache["hitratio"].WithLabelValues(coreName, name, metrics.Class).Set(hitratio)
					e.gaugeCache["hits"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.Hits))
					e.gaugeCache["inserts"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.Inserts))
					e.gaugeCache["lookups"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.Lookups))
					e.gaugeCache["size"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.Size))
					e.gaugeCache["warmup_time"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.WarmupTime))
				}
			}
		} else {
			for name, metrics := range cacheMetrics {
				if metrics.Class == "org.apache.solr.search.SolrFieldCacheMBean" {
					continue
				}
				e.gaugeCache["cumulative_evictions"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.CumulativeEvictions))
				e.gaugeCache["cumulative_hitratio"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.CumulativeHitratio))
				e.gaugeCache["cumulative_hits"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.CumulativeHits))
				e.gaugeCache["cumulative_inserts"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.CumulativeInserts))
				e.gaugeCache["cumulative_lookups"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.CumulativeLookups))
				e.gaugeCache["evictions"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.Evictions))
				e.gaugeCache["hitratio"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.Hitratio))
				e.gaugeCache["hits"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.Hits))
				e.gaugeCache["inserts"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.Inserts))
				e.gaugeCache["lookups"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.Lookups))
				e.gaugeCache["size"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.Size))
				e.gaugeCache["warmup_time"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.WarmupTime))
			}
		}
	}

	// Report metrics.
	for _, vec := range e.gaugeAdmin {
		vec.Collect(ch)
	}
	for _, vec := range e.gaugeCore {
		vec.Collect(ch)
	}
	for _, vec := range e.gaugeQuery {
		vec.Collect(ch)
	}
	for _, vec := range e.gaugeUpdate {
		vec.Collect(ch)
	}
	for _, vec := range e.gaugeCache {
		vec.Collect(ch)
	}

	// Successfully processed stats.
	e.up.Set(1)
}

func findMBeansData(mBeansData []json.RawMessage, query string) json.RawMessage {
	var decoded string
	for i := 0; i < len(mBeansData); i += 1 {
		err := json.Unmarshal(mBeansData[i], &decoded)
		if err == nil {
			if decoded == query || decoded == query + "HANDLER" {
				return mBeansData[i+1]
			}
		}
	}

	return nil
}
