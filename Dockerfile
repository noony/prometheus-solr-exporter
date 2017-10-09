FROM golang:1.9.1-alpine as builder
MAINTAINER noony <noony@users.noreply.github.com>

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go
ENV CGO_ENABLED=0

COPY . $GOPATH/src/github.com/noony/prometheus-solr-exporter
WORKDIR $GOPATH/src/github.com/noony/prometheus-solr-exporter

RUN apk add --update --no-cache \
       make \
       git \
       ca-certificates
RUN go get ./... && make build

FROM quay.io/prometheus/busybox:latest
COPY --from=builder /go/bin/prometheus-solr-exporter /bin/prometheus-solr-exporter
ENTRYPOINT ["/bin/prometheus-solr-exporter"]
EXPOSE 9231
