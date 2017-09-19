FROM        quay.io/prometheus/busybox:latest
MAINTAINER  noony <noony@users.noreply.github.com>

COPY solr_exporter /bin/solr_exporter

ENTRYPOINT ["/bin/solr_exporter"]
EXPOSE     9231