#!/bin/bash

go build -o=$DOCKER_LOCAL_MOUNT/im-worker -v ./cmd/im-service/
cp -r ./config/env $DOCKER_LOCAL_MOUNT/config

sudo docker compose up -d im-worker --build