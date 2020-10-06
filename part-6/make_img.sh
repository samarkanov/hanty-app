#!/bin/bash
app="khanty.client"

docker build -t ${app} --tag ${app}:0.0.1 .
