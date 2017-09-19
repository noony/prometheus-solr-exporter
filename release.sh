#!/bin/bash

rm -rf .build .tarballs

VERSION=$(cat VERSION)

make build
promu crossbuild
promu crossbuild tarballs
promu checksum .tarballs
promu release .tarballs


rm solr_exporter
ln -s .build/linux-amd64/solr_exporter solr_exporter

make docker DOCKER_IMAGE_NAME=noony/prometheus-solr-exporter DOCKER_IMAGE_TAG=v$VERSION
docker login
docker tag "noony/prometheus-solr-exporter:v$VERSION" "noony/prometheus-solr-exporter:latest"
docker push noony/prometheus-solr-exporter