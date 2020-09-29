#!/bin/bash
app="khanty.db"
this_service="db"

docker build -t ${app} --tag ${app}:0.0.1 .
docker run -d -p 10003:10003 \
              -e "THIS_SERVICE_NAME=${this_service}" \
              --name=${app} ${app}
