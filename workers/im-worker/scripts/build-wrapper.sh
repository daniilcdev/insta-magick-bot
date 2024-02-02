#!/bin/bash

go build -o=./im-worker -v ./cmd/im-service/

sudo docker compose up -d im-worker --build