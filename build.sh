#!/bin/bash

mv .env back.env

echo CWD=`pwd` > .env
cat back.env >> .env

# build and up worker
cd ./image-service-worker
./scripts/build-wrapper.sh
cd ../

cd ./telegram-frontend-service
./scripts/build-wrapper.sh
cd ../

source .env
sudo docker compose up -d database --build
sudo docker compose up -d image-service-worker --build
sudo docker compose up -d telegram-service --build

goose -dir=schemas postgres "user=$DB_USER password=$DB_PASS dbname=$DB_NAME $DB_EXTRA" up

mv back.env .env