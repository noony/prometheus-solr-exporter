#!/bin/sh

set -e

env

sed -i.bak s/%host%:%port%/$REMOTE_HOSTPORT/g /opt/jmx_exporter/config.yml
sed -i.bak s/%ssl%/$JMX_SSL/g /opt/jmx_exporter/config.yml

cat /opt/jmx_exporter/config.yml

SERVICE_PORT=9231
JVM_OPTS="-Dcom.sun.management.jmxremote.ssl=false -Dcom.sun.management.jmxremote.authenticate=false -Dcom.sun.management.jmxremote.port=18983"

java $JVM_OPTS -jar /opt/jmx_exporter/jmx_prometheus_httpserver-$VERSION-jar-with-dependencies.jar $SERVICE_PORT /opt/jmx_exporter/config.yml
