# BTC Service

BTC Service handles BTC transaction and User balance related data.

## Table of Contents

<!-- @import "[TOC]" {cmd="toc" depthFrom=1 depthTo=6 orderedList=false} -->

<!-- code_chunk_output -->

- [Table of Contents](#table-of-contents)
- [Project Summary](#project-summary)
- [Installation](#installation)
    - [1. Set Up Golang Development Environment](#1-set-up-golang-development-environment)
    - [2. Install Development Utility Tools](#2-install-development-utility-tools)
- [Development workflow and guidelines](#development-workflow-and-guidelines)
    - [1. API](#1-api)
    - [2. TimescaleDB + GUI](#2-timescaledb--gui)
    - [3. Migration](#3-migration)
    - [4. Unit Test](#4-unit-test)
    - [5. Linter](#5-linter)
    - [6. Run the service](#6-run-the-service)
- [Documentation](#documentation)
  - [Visualize Code Diagram](#visualize-code-diagram)
  - [RPC Sequence Diagram](#rpc-sequence-diagram)

<!-- /code_chunk_output -->

## Project Summary

| Item                     | Description                                                                                                           |
|--------------------------|-----------------------------------------------------------------------------------------------------------------------|
| Golang Version           | [1.19](https://golang.org/doc/go1.19)                                                                                 |
| Database                 | [timescale](https://www.timescale.com)                                                                                |
| Migration                | [migrate](https://github.com/golang-migrate/migrate)                                                                  |
| moq                      | [mockgen](https://github.com/golang/mock)                                                                             |
| Linter                   | [GolangCI-Lint](https://github.com/golangci/golangci-lint)                                                            |
| Testing                  | [testing](https://golang.org/pkg/testing/) and [testify/assert](https://godoc.org/github.com/stretchr/testify/assert) |
| API                      | [gRPC](https://grpc.io/docs/tutorials/basic/go/)                                                                      |
| Application Architecture | [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)                    |
| Directory Structure      | [Standard Go Project Layout](https://github.com/golang-standards/project-layout)                                      |
| CI (Lint & Test)         | [GitHubActions](https://github.com/features/actions)                                                                  |
| Visualize Code Diagram   | [go-callviz](https://github.com/ofabry/go-callvis)                                                                    |
| Sequence Diagram         | [Mermaid](https://mermaid.js.org)                                                                                     |
| Protobuf Operations      | [buf](https://buf.build)                                                                                              |

## Installation

### 1. Set Up Golang Development Environment

See the following page to download and install Golang.

https://go.dev/doc/install

### 2. Install Development Utility Tools

You can install all tools for development and deployment for this service by running:

```sh
$ go mod download
```

```sh
$ make install
```

---

## Development workflow and guidelines

### 1. API

This project using gRPC and Protocol Buffers, thus all needed data like Service definition, RPC's list, Entities will store in [api/proto](api/proto) directory.

If you unfamiliar with Protocol Buffer, please visit this link for the detail:

* https://protobuf.dev

For generating the Proto files, make sure to have these libs installed on your system, please refer to this link:

* https://buf.build/
* https://grpc.io/docs/protoc-installation
* https://grpc.io/docs/languages/go/quickstart/

The validation for this API using `protoc-gen-validate`, for the detail please refer to this lib:

* https://github.com/bufbuild/protoc-gen-validate

Then, generating the Protobuf files can be done my this command:

```sh
$ make protoc
```

### 2. TimescaleDB + GUI

Run TimescaleDB locally with the GUI (pgAdmin) can be executed with the following docker-compose command:

```sh
$ docker-compose -f ./development/docker-compose.yml up timescaledb pgadmin
```

> NOTE:
> TimescaleDB will use port 5432 and pgAdmin will use 5050, please make sure those port are unused in yur system.
> If the port conflicted, you can change the port on the [development/docker-compose.yml](docker-compose.yml) file.

If you don't have a docker-compose installed, please refer to this page https://docs.docker.com/compose/

### 3. Migration

Make sure the database already running, after that we need some tables and dummy data in order to test the application, please run this command to do the migration:

```sh
$ docker-compose -f ./development/docker-compose.yml up migration
```

### 4. Unit Test

Make sure the database already running, then you can simply execute the following command to run all test cases in this service:

```sh
$ make test
```

### 5. Linter

For running the linter make sure these libs already installed in your system:

* https://github.com/golangci/golangci-lint
* https://github.com/yoheimuta/protolint

Then checks the Go and Proto code style using lint can be done with this command:

```sh
$ make lint
```

### 6. Run the service

For running the service, you need the database running and set up some env variables:

```
export APP_ENV=dev
export SERVER_PORT=8080
export POSTGRES_USER=test
export POSTGRES_PASSWORD=test
export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432
export POSTGRES_DB=test
```

Or you can just execute the sh file:

```sh
$ ./scripts/run.sh
```

# Documentation

## Visualize Code Diagram

To help give a better understanding about reading the code
such as relations with packages and types, here are some diagrams listed
generated automatically using [https://github.com/ofabry/go-callvis](https://github.com/ofabry/go-callvis)

<!-- start diagram doc -->
1. [main diagram](docs/diagrams/main.png)
2. [di diagram](docs/diagrams/di.png)
3. [handler diagram](docs/diagrams/handler.png)
4. [usecases diagram](docs/diagrams/usecases.png)
5. [datastore diagram](docs/diagrams/datastore.png)

<!-- end diagram doc -->

## RPC Sequence Diagram

To help give a better understanding about reading the RPC flow
such as relations with usecases and repositories, here are some sequence diagrams (generated automatically) listed in Markdown file and written in Mermaid JS [https://mermaid-js.github.io/mermaid/](https://mermaid-js.github.io/mermaid/) format.

To generate the RPC sequence diagram, there's a Makefile command that can be use:

1. Run this command to generate specific RPC `make sequence-diagram RPC=GetData`.
2. For generates multiple RPC's, just adding the other RPC by comma `make sequence-diagram RPC=GetData,GetList`.
3. For generates all RPC's, use wildcard * in the parameter `make sequence-diagram RPC=*`.

<!-- start rpc sequence diagram doc -->
1. [CreateTransaction RPC - Sequence Diagram](docs/sequence-diagrams/rpc/create-transaction.md)
2. [GetUserBalance RPC - Sequence Diagram](docs/sequence-diagrams/rpc/get-user-balance.md)
3. [ListTransaction RPC - Sequence Diagram](docs/sequence-diagrams/rpc/list-transaction.md)

<!-- end rpc sequence diagram doc -->
