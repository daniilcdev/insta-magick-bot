#!/bin/bash

cd workers/im-worker
pwd
go build -o=./im-worker -v ./cmd/im-service/...

cd ../../
sudo docker compose up -d im-worker --build