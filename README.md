# Solr Exporter

[![Docker Pulls](https://img.shields.io/docker/pulls/noony/prometheus-solr-exporter.svg?maxAge=604800)](https://hub.docker.com/r/noony/prometheus-solr-exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/noony/prometheus-solr-exporter)](https://goreportcard.com/report/github.com/noony/prometheus-solr-exporter)

Prometheus exporter for various metrics about Solr, written in Go.

### Installation

For pre-built binaries please take a look at the releases.  
https://github.com/noony/prometheus-solr-exporter

#### Docker

```bash
docker pull noony/prometheus-solr-exporter
docker run noony/prometheus-solr-exporter --solr.address=http://url-to-solr:port
```

#### Configuration

Below is the command line options summary:

```bash
prometheus-solr-exporter --help
```

| Argument              | Description |
| --------              | ----------- |
| solr.address          | URI on which to scrape Solr. (default "http://localhost:8983") |
| solr.context-path     | Solr webapp context path. (default "/solr") |
| solr.pid-file         | Path to Solr pid file |
| solr.timeout          | Timeout for trying to get stats from Solr. (default 5s) |
| solr.excluded-core    | Regex to exclude core from monitoring|
| web.listen-address    | Address to listen on for web interface and telemetry. (default ":9231")|
| web.telemetry-path    | Path under which to expose metrics. (default "/metrics")|

### Building

Clone the repository and just launch this command
```bash
make build
```

### Testing

[![Build Status](https://travis-ci.org/noony/prometheus-solr-exporter.png?branch=master)][travisci]

```bash
make test
```

[travisci]: https://travis-ci.org/noony/prometheus-solr-exporter

### Grafana dashboard

See https://grafana.com/dashboards/2551

## License

Apache License 2.0, see [LICENSE](https://github.com/noony/prometheus-solr-exporter/blob/master/LICENSE).
