#!/bin/bash

go build -o=/mnt/storage/im-worker-volume/im-worker -v ./cmd/im-service/
cp -r ./config/env /mnt/storage/im-worker-volume/config
sudo docker build -t im-worker .
#sudo docker compose up -d im-worker --build