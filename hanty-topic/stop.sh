#!/bin/bash
app="hanty.topic"

docker stop ${app}
docker rm --force ${app}
