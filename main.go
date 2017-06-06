package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
)

const (
	namespace = "solr"
	indexHTML = `
    <html>
        <head>
            <title>Solr Exporter</title>
        </head>
        <body>
            <h1>Solr Exporter</h1>
            <p>
            <a href='%s'>Metrics</a>
            </p>
        </body>
    </html>`
)

func main() {
	const pidFileHelpText = `Path to Solr pid file.

    If provided, the standard process metrics get exported for the Solr
    process, prefixed with 'solr_process_...'. The solr_process exporter
    needs to have read access to files owned by the Solr process. Depends on
    the availability of /proc.

    https://prometheus.io/docs/instrumenting/writing_clientlibs/#process-metrics.`

	var (
		listenAddress = flag.String("web.listen-address", ":9231", "Address to listen on for web interface and telemetry.")
		metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
		solrURI       = flag.String("solr.address", "http://localhost:8983", "URI on which to scrape Solr.")
		solrTimeout   = flag.Duration("solr.timeout", 5*time.Second, "Timeout for trying to get stats from Solr.")
		solrPidFile   = flag.String("solr.pid-file", "", pidFileHelpText)
		showVersion   = flag.Bool("version", false, "Print version information.")
	)
	flag.Parse()

	if *showVersion {
		fmt.Fprintln(os.Stdout, version.Print("solr_exporter"))
		os.Exit(0)
	}

	log.Infoln("Starting solr_exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	exporter := NewExporter(*solrURI, *solrTimeout)
	prometheus.MustRegister(exporter)
	prometheus.MustRegister(version.NewCollector("solr_exporter"))

	if *solrPidFile != "" {
		procExporter := prometheus.NewProcessCollectorPIDFn(
			func() (int, error) {
				content, err := ioutil.ReadFile(*solrPidFile)
				if err != nil {
					return 0, fmt.Errorf("Can't read pid file: %s", err)
				}
				value, err := strconv.Atoi(strings.TrimSpace(string(content)))
				if err != nil {
					return 0, fmt.Errorf("Can't parse pid file: %s", err)
				}
				return value, nil
			}, namespace)
		prometheus.MustRegister(procExporter)
	}

	log.Infoln("Listening on", *listenAddress)
	http.Handle(*metricsPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Haproxy Exporter</title></head>
             <body>
             <h1>Solr Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
