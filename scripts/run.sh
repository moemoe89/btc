#!/bin/bash

export SERVER_PORT=8080

go build ./cmd/main.go && ./main
