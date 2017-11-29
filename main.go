package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
)

var (
	listenAddress    = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9231").String()
	metricsPath      = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	solrURI          = kingpin.Flag("solr.address", "URI on which to scrape Solr.").Default("http://localhost:8983").String()
	solrContextPath  = kingpin.Flag("solr.context-path", "Solr webapp context path.").Default("/solr").String()
	solrExcludedCore = kingpin.Flag("solr.excluded-core", "Regex to exclude core from monitoring").Default("").String()
	solrTimeout      = kingpin.Flag("solr.timeout", "Timeout for trying to get stats from Solr.").Default("5s").Duration()
	solrPidFile      = kingpin.Flag("solr.pid-file", "").Default(pidFileHelpText).String()
	ignoreCore       = kingpin.Flag("solr.ignore-core", "comma-separated list of cores to ignore").Default("").String()
)

var landingPage = []byte(`<html>
<head><title>Solr exporter</title></head>
<body>
<h1>Solr exporter</h1>
<p><a href='` + *metricsPath + `'>Metrics</a></p>
</body>
</html>
`)

func main() {
	log.AddFlags(kingpin.CommandLine)
	kingpin.Version(version.Print("solr_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	log.Infoln("Starting solr_exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	exporter := NewExporter(*solrURI, *solrContextPath, *solrTimeout, *solrExcludedCore, *ignoreCore)
	registry := prometheus.NewRegistry()
	registry.Register(exporter)
	registry.Register(version.NewCollector("solr_exporter"))

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
			}, "solr")
		registry.Register(procExporter)
	}

	log.Infoln("Listening on", *listenAddress)
	http.Handle(*metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(landingPage)
	})
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
