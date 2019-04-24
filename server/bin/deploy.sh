#!/bin/bash

if [[ $# -ne 1 ]]; then
    echo "Usage: deploy.sh [FunctionName]"
    exit 0
fi

go build

goBuildRes=$?
if [[ ${goBuildRes} -ne 0 ]]; then
    echo "Not deploying"
    exit goBuildRes
fi

gcloud functions deploy $1 --runtime go111 --trigger-http
