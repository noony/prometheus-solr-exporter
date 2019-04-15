package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

var pingPath = "/admin/ping?wt=json"

//PingCollector collects ping  metrics from solr
type PingCollector struct {
	ping *prometheus.Desc

	client      http.Client
	pingURL     string
	solrBaseURL string
}

// NewPingCollector returns a new Collector exposing solr ping statistics.
func NewPingCollector(client http.Client, solrBaseURL string) (*PingCollector, error) {
	pingURL := fmt.Sprintf("%s/%%s%s", solrBaseURL, pingPath)
	return &PingCollector{
		ping: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "ping"),
			"See following URL: https://lucene.apache.org/solr/guide/ping.html",
			[]string{"core"}, nil,
		),
		client:      client,
		pingURL:     pingURL,
		solrBaseURL: solrBaseURL,
	}, nil
}

// Update exposes ping related metrics from solr.
func (c *PingCollector) Update(ch chan<- prometheus.Metric) error {
	coreList, err := getCoreList(c.client, c.solrBaseURL)
	if err != nil {
		return err
	}
	for _, core := range coreList {
		_, err = fetchHTTP(c.client, fmt.Sprintf(c.pingURL, core))
		if err != nil {
			ch <- prometheus.MustNewConstMetric(c.ping, prometheus.GaugeValue, float64(0), core)
		} else {
			ch <- prometheus.MustNewConstMetric(c.ping, prometheus.GaugeValue, float64(1.0), core)
		}
	}

	return nil
}

// Collect implements the prometheus.Collector interface.
func (c *PingCollector) Collect(ch chan<- prometheus.Metric) {
	if err := c.Update(ch); err != nil {
		log.Errorf("Failed to collect metrics: %v", err)
	}
}

// Describe implements the prometheus.Collector interface.
func (c *PingCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.ping
}
