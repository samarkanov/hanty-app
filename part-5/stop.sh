#!/bin/bash
app="khanty.master"

docker stop ${app}
docker rm --force ${app}
