#!/bin/bash

docker-compose -f docker-compose.yml up -d
make dropdb
make createdb
make migratedbup

