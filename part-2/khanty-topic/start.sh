#!/bin/bash
app="khanty.topic"
this_service="topic"

docker build -t ${app} --tag ${app}:0.0.1 .
docker run -d -p 10001:10001 \
              -e "THIS_SERVICE_NAME=${this_service}" \
              --name=${app} ${app}
