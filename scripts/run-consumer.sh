#!/bin/bash

export GRPC_SERVER=localhost:8080
export RABBITMQ_HOST=amqp://guest:guest@localhost:5672/

go build -o main-consumer ./cmd/consumer/main.go && ./main-consumer
