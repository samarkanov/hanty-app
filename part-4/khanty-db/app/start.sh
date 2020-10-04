#!/bin/bash
app="khanty.db"

docker build -t ${app} --tag ${app}:0.0.1 .
docker run -d -p 10003:10003 --name=${app} ${app}
