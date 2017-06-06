# Solr Exporter

[![Docker Pulls](https://img.shields.io/docker/pulls/noony/prometheus-solr-exporter.svg?maxAge=604800)](https://hub.docker.com/r/noony/prometheus-solr-exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/noony/prometheus-solr-exporter)](https://goreportcard.com/report/github.com/noony/prometheus-solr-exporter)

Prometheus exporter for various metrics about Solr, written in Go.

### Installation

```bash
go get -u github.com/noony/prometheus-solr-exporter
```

### Configuration

```bash
prometheus-solr-exporter --help
```

| Argument              | Description |
| --------              | ----------- |
| solr.address          | URI on which to scrape Solr. (default "http://localhost:8983") |
| solr.pid-file         | Path to Solr pid file |
| solr.timeout          | Timeout for trying to get stats from Solr. (default 5s) |
| web.listen-address    | Address to listen on for web interface and telemetry. (default ":9231")|
| web.telemetry-path    | Path under which to expose metrics. (default "/metrics")|

### Building

```bash
make build
```

### Testing

[![Build Status](https://travis-ci.org/noony/prometheus-solr-exporter.png?branch=master)][travisci]

```bash
make test
```

[travisci]: https://travis-ci.org/noony/prometheus-solr-exporter

## License

Apache License 2.0, see [LICENSE](https://github.com/noony/prometheus-solr-exporter/blob/master/LICENSE).