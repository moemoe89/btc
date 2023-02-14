#!/bin/bash

export SERVER_PORT=8080

go build -o main-consumer ./cmd/consumer/main.go && ./main-consumer
