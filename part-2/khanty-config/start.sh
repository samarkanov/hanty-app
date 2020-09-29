#!/bin/bash
app="khanty.config"
this_service="config"
config_file="/tmp/app/config.json"

docker build -t ${app} --tag ${app}:0.0.1 .
docker run -d -p 10002:10002 \
              -e "THIS_SERVICE_NAME=${this_service}" \
              -e "KHANTY_CONFIG_FILE=${config_file}" \
              --name=${app} ${app}
