#!/bin/bash

go mod tidy
docker-compose -f docker-compose.yml up -d
make dropdb
make createdb
make migratedbup

