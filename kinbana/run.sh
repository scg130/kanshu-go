#!/bin/bash
docker rm -f kibana
 docker run -d --log-driver json-file --log-opt max-size=100m --log-opt max-file=5 \
 --name kibana \
 -p 5601:5601 \
 -e TZ="Asia/Shanghai" \
 --privileged=true \
 -v ./kibana.yml:/usr/share/kibana/config/kibana.yml \
 kibana:7.17.18