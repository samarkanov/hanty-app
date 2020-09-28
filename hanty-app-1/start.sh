#!/bin/bash
app="hanty.topic"
this_service="topic"
config_file="/tmp/app/config.json"

docker build -t ${app} --tag ${app}:0.0.1 .
docker run -d -p 10001:10001 \
              -e "THIS_SERVICE_NAME=${this_service}" \
              -e "HANTY_CONFIG_FILE=${config_file}" \
              --name=${app} ${app}
