#!/bin/bash

# build and up worker
cd ./workers/im-worker
./scripts/build-wrapper.sh
cd ../../
sudo docker compose up -d im-worker --build