#!/bin/bash

if [ $# -ne 1 ]; then
    echo "Usage: deploy.sh [FunctionName]"
    exit 0
fi

gcloud functions deploy $1 --runtime go111 --trigger-http
