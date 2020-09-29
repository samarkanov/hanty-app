#!/bin/bash
app="khanty.config"

docker stop ${app}
docker rm --force ${app}
