# Solr Exporter for Prometheus

[![Docker Pulls](https://img.shields.io/docker/pulls/noony/prometheus-solr-exporter.svg?maxAge=604800)][hub]
[![Docker Build Status](https://img.shields.io/docker/build/noony/prometheus-solr-exporter.svg)][hub]

A [Solr](http://lucene.apache.org/solr/) exporter for prometheus.

## Building

The solr exporter exports metrics from a solr server for
consumption by prometheus.

Environment variable : REMOTE_HOSTPORT (ie : 10.10.10.10:18983)

By default the solr\_exporter serves on port `0.0.0.0:9231` at `/metrics`

[hub]: https://hub.docker.com/r/noony/prometheus-solr-exporter/