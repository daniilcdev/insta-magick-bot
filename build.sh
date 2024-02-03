#!/bin/bash

# build and up worker
cd ./workers/im-worker
./scripts/build-wrapper.sh
cd ../../

cd ./telegram-frontend-service
./scripts/build-wrapper.sh
cd ../

source .env
sudo docker compose up -d database --build
sudo docker compose up -d im-worker --build
sudo docker compose up -d telegram-service --build