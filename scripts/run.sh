#!/bin/bash

export SERVER_PORT=8080
export POSTGRES_USER=test
export POSTGRES_PASSWORD=test
export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432
export POSTGRES_DB=test

go build ./cmd/main.go && ./main
