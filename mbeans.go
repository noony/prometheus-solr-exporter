package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func findMBeansData(mBeansData []json.RawMessage, query string) json.RawMessage {
	var decoded string
	for i := 0; i < len(mBeansData); i++ {
		err := json.Unmarshal(mBeansData[i], &decoded)
		if err == nil {
			if decoded == query || decoded == query+"HANDLER" {
				return mBeansData[i+1]
			}
		}
	}

	return nil
}

func processMbeans(e *Exporter, coreName string, data io.Reader) []error {
	mBeansData := &MBeansData{}
	errors := []error{}
	if err := json.NewDecoder(data).Decode(mBeansData); err != nil {
		errors = append(errors, fmt.Errorf("Failed to unmarshal mbeansdata JSON into struct: %v", err))
		return errors
	}

	var coreMetrics map[string]Core
	if err := json.Unmarshal(findMBeansData(mBeansData.SolrMbeans, "CORE"), &coreMetrics); err != nil {
		errors = append(errors, fmt.Errorf("Failed to unmarshal mbeans core metrics JSON into struct: %v", err))
		return errors
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
		errors = append(errors, fmt.Errorf("Failed to unmarshal mbeans query metrics JSON into struct: %v, json : %s", err, b))
		return errors
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
		errors = append(errors, fmt.Errorf("Failed to unmarshal mbeans update metrics JSON into struct: %v", err))
		return errors
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

	cacheData := findMBeansData(mBeansData.SolrMbeans, "CACHE")
	b = bytes.Replace(cacheData, []byte(":\"NaN\""), []byte(":0.0"), -1)
	b = bytes.Replace(b, []byte("CACHE.searcher.perSegFilter."), []byte(""), -1)
	b = bytes.Replace(b, []byte("CACHE.searcher.queryResultCache."), []byte(""), -1)
	b = bytes.Replace(b, []byte("CACHE.searcher.fieldValueCache."), []byte(""), -1)
	b = bytes.Replace(b, []byte("CACHE.searcher.filterCache."), []byte(""), -1)
	b = bytes.Replace(b, []byte("CACHE.searcher.documentCache."), []byte(""), -1)
	mbeanerrs := handleCacheMbeans(b, e, coreName)
	for _, e := range mbeanerrs {
		errors = append(errors, e)
	}
	return errors
}

func handleCacheMbeans(data []byte, e *Exporter, coreName string) []error {
	var cacheMetrics map[string]Cache
	var errors = []error{}
	if err := json.Unmarshal(data, &cacheMetrics); err != nil {
		errors = append(errors, fmt.Errorf("Failed to unmarshal mbeans cache metrics JSON into struct (core : %s): %v, json : %s", coreName, err, data))
	} else {
		for name, metrics := range cacheMetrics {
			if metrics.Class == "org.apache.solr.search.SolrFieldCacheMBean" || metrics.Class == "org.apache.solr.search.SolrFieldCacheBean" {
				continue
			}
			hitratio, err := strconv.ParseFloat(string(metrics.Stats.Hitratio), 64)
			if err != nil {
				errors = append(errors, fmt.Errorf("Fail to convert Hitratio in float: %v", err))
			}
			cumulativeHitratio, err := strconv.ParseFloat(string(metrics.Stats.CumulativeHitratio), 64)
			if err != nil {
				errors = append(errors, fmt.Errorf("Fail to convert Cumulative Hitratio in float: %v", err))
			}
			e.gaugeCache["cumulative_evictions"].WithLabelValues(coreName, name, metrics.Class).Set(float64(metrics.Stats.CumulativeEvictions))
			e.gaugeCache["cumulative_hitratio"].WithLabelValues(coreName, name, metrics.Class).Set(cumulativeHitratio)
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
	return errors
}
