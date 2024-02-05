#!/bin/bash
source ../../.env

go build -o=$DOCKER_LOCAL_MOUNT/im-worker -v ./cmd/im-service/

# copy env config
mkdir -p $DOCKER_LOCAL_MOUNT/config
cp -r ./config/env/ $DOCKER_LOCAL_MOUNT/config/

# make resource folders
mkdir -p $DOCKER_LOCAL_MOUNT/res/tmp \
 $DOCKER_LOCAL_MOUNT/res/pending \
 $DOCKER_LOCAL_MOUNT/res/processed
