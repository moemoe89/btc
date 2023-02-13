#!/bin/bash

export APP_ENV=dev
export SERVER_PORT=8080
export POSTGRES_USER=test
export POSTGRES_PASSWORD=test
export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432
export POSTGRES_DB=test
export OTEL_AGENT=http://localhost:14268/api/traces

go build ./cmd/main.go && ./main
