#!/bin/bash

cp .env back.env

echo PWD=`pwd` > .env
cat back.env >> .env

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

mv back.env .env