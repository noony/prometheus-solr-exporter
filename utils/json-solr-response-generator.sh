#!/bin/bash

solr_tags=$(curl -s https://hub.docker.com/v2/repositories/library/solr/tags/?page_size=100 | jq -r '.results|.[]|.name' | grep -v slim | grep -v alpine | grep -v latest)

for j in ${solr_tags}
do
    echo "Generate json for solr tag : $j"
    nohup docker run --name solr-$j -p8983:8983 solr:$j &
    until $(curl --output /dev/null --silent --head --fail "http://localhost:8983/solr/#/"); do
        echo "Solr is unavailable - sleeping"
        sleep 1
    done
    echo "Solr $j ready !"
    sleep 2
    docker exec -ti solr-$j bin/solr create_core -c gettingstarted
    sleep 2
    until $(curl --output /dev/null --silent --head --fail "http://localhost:8983/solr/gettingstarted/admin/mbeans?stats=true&wt=json&cat=CORE&cat=QUERYHANDLER&cat=UPDATEHANDLER&cat=CACHE"); do
        echo "Core gettingstated stats are unavailable - sleeping"
        sleep 1
    done
    mkdir generated-json/$j || true
    curl -o generated-json/$j/admin-cores.json --silent --fail "http://localhost:8983/solr/admin/cores?action=STATUS&wt=json"
    curl -o generated-json/$j/mbeans.json --silent --fail "http://localhost:8983/solr/gettingstarted/admin/mbeans?stats=true&wt=json&cat=CORE&cat=QUERYHANDLER&cat=UPDATEHANDLER&cat=CACHE"
    docker stop solr-$j
done