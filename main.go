package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	namespace       = "solr"
	pidFileHelpText = `Path to Solr pid file.

	If provided, the standard process metrics get exported for the Solr
	process, prefixed with 'solr_process_...'. The solr_process exporter
	needs to have read access to files owned by the Solr process. Depends on
	the availability of /proc.

	https://prometheus.io/docs/instrumenting/writing_clientlibs/#process-metrics.`
	adminCoresPath = "/admin/cores?action=STATUS&wt=json"
)

var (
	listenAddress    = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9231").String()
	metricsPath      = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	solrURI          = kingpin.Flag("solr.address", "URI on which to scrape Solr.").Default("http://localhost:8983").String()
	solrContextPath  = kingpin.Flag("solr.context-path", "Solr webapp context path.").Default("/solr").String()
	solrExcludedCore = kingpin.Flag("solr.excluded-core", "Regex to exclude core from monitoring").Default("").String()
	solrTimeout      = kingpin.Flag("solr.timeout", "Timeout for trying to get stats from Solr.").Default("5s").Duration()
	solrPidFile      = kingpin.Flag("solr.pid-file", "").Default(pidFileHelpText).String()
)

func main() {
	log.AddFlags(kingpin.CommandLine)
	kingpin.Version(version.Print("solr_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	log.Infoln("Starting solr_exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, *solrTimeout)
				if err != nil {
					return nil, err
				}
				if err := c.SetDeadline(time.Now().Add(*solrTimeout)); err != nil {
					return nil, err
				}
				return c, nil
			},
		},
	}

	solrBaseURL := fmt.Sprintf("%s%s", *solrURI, *solrContextPath)

	pingExporter, err := NewPingCollector(*client, solrBaseURL)
	if err != nil {
		log.Errorf("Failed to create ping metrics collector: %v", err)
	}
	prometheus.MustRegister(pingExporter)

	metricsExporter, err := NewMetricsCollector(*client, solrBaseURL)
	if err != nil {
		log.Errorf("Failed to create ping metrics collector: %v", err)
	}
	prometheus.MustRegister(metricsExporter)

	if *solrPidFile != "" {
		procExporter := prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{
			PidFn: func() (int, error) {
				content, err := ioutil.ReadFile(*solrPidFile)
				if err != nil {
					return 0, fmt.Errorf("Can't read pid file: %s", err)
				}
				value, err := strconv.Atoi(strings.TrimSpace(string(content)))
				if err != nil {
					return 0, fmt.Errorf("Can't parse pid file: %s", err)
				}
				return value, nil
			},
			Namespace: "solr",
		})
		prometheus.MustRegister(procExporter)
	}

	log.Infoln("Listening on", *listenAddress)
	http.Handle(*metricsPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Solr Exporter</title></head>
             <body>
             <h1>Solr Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}

func fetchHTTP(client http.Client, url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return nil, fmt.Errorf("HTTP status %d url : %s", resp.StatusCode, url)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body: %v", err)
	}
	return body, nil
}

// Return list of cores from solr server
func getCoreList(client http.Client, solrBaseURL string) ([]string, error) {
	body, err := fetchHTTP(client, fmt.Sprintf("%s%s", solrBaseURL, adminCoresPath))
	if err != nil {
		return nil, fmt.Errorf("Fail to get solr core list JSON into struct: %v", err)
	}
	adminCoresName := &AdminCoresName{}
	err = json.Unmarshal(body, adminCoresName)
	if err != nil {
		return nil, fmt.Errorf("Fail to unmarshal solr core list JSON into struct: %v", err)
	}
	serverCores := []string{}
	for _, v := range adminCoresName.Status {
		serverCores = append(serverCores, v.Name)
	}
	return serverCores, nil
}
