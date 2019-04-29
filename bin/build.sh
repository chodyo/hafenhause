#!/bin/bash

go mod tidy

docker build -t hafenhause .

docker-compose up -d --force-recreate

docker-compose logs -f
