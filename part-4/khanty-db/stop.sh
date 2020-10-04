#!/bin/bash
app="khanty.db"

docker stop ${app}
docker rm --force ${app}
