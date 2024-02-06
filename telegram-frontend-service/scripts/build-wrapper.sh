#!/bin/bash
source ../.env

go build -o=$DOCKER_LOCAL_MOUNT/telegram-service -v ./cmd/

# make config files
mkdir -p $DOCKER_LOCAL_MOUNT/env/private
echo DB_DRIVER=postgres > $DOCKER_LOCAL_MOUNT/env/private/db.env
echo DB_CONN="user=$DB_USER password=$DB_PASS dbname=$DB_NAME host=postgres sslmode=disable" \
    >> $DOCKER_LOCAL_MOUNT/env/private/db.env

echo TELEGRAM_BOT_TOKEN=$TELEGRAM_TOKEN > $DOCKER_LOCAL_MOUNT/env/private/telegram.env
echo PROCESSED_FILES_DIR=./res/processed/ >> $DOCKER_LOCAL_MOUNT/env/private/telegram.env