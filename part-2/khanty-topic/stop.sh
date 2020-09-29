#!/bin/bash
app="khanty.topic"

docker stop ${app}
docker rm --force ${app}
