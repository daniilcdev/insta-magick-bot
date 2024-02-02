#!/bin/bash
go version
echo $DOCKER_LOCAL_MOUNT
go build -o=$DOCKER_LOCAL_MOUNT/im-worker -v ./workers/im-worker/cmd/im-service/
cp -r ./workers/im-worker/config/env $DOCKER_LOCAL_MOUNT/config

sudo docker compose up -d im-worker --build