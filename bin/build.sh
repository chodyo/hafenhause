#!/bin/bash

go mod tidy

docker build -t hafenhause .
