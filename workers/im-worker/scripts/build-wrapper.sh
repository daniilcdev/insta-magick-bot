#!/bin/bash
source ../../.env

go build -o=$DOCKER_LOCAL_MOUNT/im-worker -v ./cmd/im-service/

# make resource folders
mkdir -p $DOCKER_LOCAL_MOUNT/res/tmp \
 $DOCKER_LOCAL_MOUNT/res/pending \
 $DOCKER_LOCAL_MOUNT/res/processed

# make config file
mkdir -p $DOCKER_LOCAL_MOUNT/config/env

echo IM_IN_DIR="./res/pending/" > $DOCKER_LOCAL_MOUNT/config/env/imagemagick.env
echo IM_OUT_DIR="./res/processed/" >> $DOCKER_LOCAL_MOUNT/config/env/imagemagick.env
echo IM_TEMP_DIR="./res/tmp/" >> $DOCKER_LOCAL_MOUNT/config/env/imagemagick.env
