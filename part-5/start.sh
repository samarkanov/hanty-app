#!/bin/bash
app="khanty.master"

docker build -t ${app} --tag ${app}:0.0.1 .
docker run -d -p 10004:10004 --name=${app} ${app}
