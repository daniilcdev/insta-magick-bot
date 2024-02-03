#!/bin/bash

source ../.env

go build -o=$DOCKER_LOCAL_MOUNT/telegram-service -v ./cmd/

# copy env config
cp -r ./config/env/ $DOCKER_LOCAL_MOUNT/
