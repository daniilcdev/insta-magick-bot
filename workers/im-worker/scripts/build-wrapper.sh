#!/bin/bash

go build -o=./im-worker -v ./cmd/im-service/
sudo docker build --no-cache -t im-worker .