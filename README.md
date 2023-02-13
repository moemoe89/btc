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
    - [7. Test the service](#7-test-the-service)
    - [8. Loading Testing](#8-load-testing)
- [Documentation](#documentation)
  - [Visualize Code Diagram](#visualize-code-diagram)
  - [RPC Sequence Diagram](#rpc-sequence-diagram)

<!-- /code_chunk_output -->

## Project Summary

| Item                     | Description                                                                                                         |
|--------------------------|---------------------------------------------------------------------------------------------------------------------|
| Golang Version           | [1.19](https://golang.org/doc/go1.19)                                                                               |
| Database                 | [timescale](https://www.timescale.com)                                                                              |
| Migration                | [migrate](https://github.com/golang-migrate/migrate)                                                                |
| moq                      | [mockgen](https://github.com/golang/mock)                                                                           |
| Linter                   | [GolangCI-Lint](https://github.com/golangci/golangci-lint)                                                          |
| Testing                  | [testing](https://golang.org/pkg/testing/) and [testify/assert](https://godoc.org/github.com/stretchr/testify/assert) |
| Load Testing             | [ghz](https://ghz.sh)                                                 |
| API                      | [gRPC](https://grpc.io/docs/tutorials/basic/go/)                                                                    |
| Application Architecture | [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)                  |
| Directory Structure      | [Standard Go Project Layout](https://github.com/golang-standards/project-layout)                                    |
| CI (Lint & Test)         | [GitHubActions](https://github.com/features/actions)                                                                |
| Visualize Code Diagram   | [go-callviz](https://github.com/ofabry/go-callvis)                                                                  |
| Sequence Diagram         | [Mermaid](https://mermaid.js.org)                                                                                   |
| Protobuf Operations      | [buf](https://buf.build)                                                                                            |

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

The default email & password for pgAdmin are:
* email: `admin@admin.com` 
* password: `admin123`

With this following TimescaleDB info:
* host: `timescaledb` -> change this to `localhost` if you try to connect from outside Docker
* port: `5432`
* username: `test`
* password: `test`
* db: `test`

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

### 7. Test the service

The example how to call the gRPC service written in Golang can be seen on this [example-client](scripts/example-client) file.

If you want to test by GUI client, you can use either BloomRPC (although already no longer active) or Postman.
For the detail please visit these links:
* https://github.com/bloomrpc/bloomrpc
* https://www.postman.com

Basically you just need to import the [api/proto/service.proto](api/proto/service.proto) file if you want to test via BloomRPC / Postman.

> NOTE: There will be a possibility issue when importing the proto file to BloomRPC or Postman.
> It is caused by some path and the usage of `protoc-gen-validate` library.
> To solve this issue, there's need a modification for the proto file.

### BloomRPC

blooRPC will have this issue when trying to import the proto file

```
Error while importing protos
no such type: e.Transaction
```

To fix this issue, fix the import path from:

```protobuf
import "proto/entity.proto";
```

To this:

```protobuf
import "../proto/entity.proto";
```

### Postman

There's some issue when importing to Postman. Basically we need to do the same things like BloomRPC and disable the validate import.

```protobuf
import "proto/entity.proto";
import "validate/validate.proto";
```

To this:

```protobuf
import "../proto/entity.proto";
// import "validate/validate.proto";
```

Also don't forget to set the import path e.g. `{YOUR-DIR}/btc/api/proto`

## 8. Load Testing

In order to make sure the service ready to handle a big traffic, it will better if we can do Load Testing to see the performance.

Since the service running in gRPC, we need the tool that support to do HTTP2 request.
In this case we can use https://ghz.sh/ because it is very simple and can generate various output report type.

> NOTE: Like importing the proto file to BloomRPC / Postman,
> when running the `ghz` there's will be issue shown due to the tool can't read the path & validate lib.

Here are some possibility issues when we're trying to run the `ghz` commands:
* `./api/proto/service.proto:5:8: open api/proto/proto/entity.proto: no such file or directory`
* `./api/proto/service.proto:7:8: open api/proto/validate/validate.proto: no such file or directory`
* `./api/proto/service.proto:29:22: field CreateTransactionRequest.user_id: unknown extension validate.rules`

To fix this issue, you need to change some file in proto file:

```protobuf
import "proto/entity.proto";
import "validate/validate.proto";
```

To this:

```protobuf
import "../proto/entity.proto";
// import "validate/validate.proto";
```

And all validation on each field such as:

```protobuf
// CreateTransactionRequest
message CreateTransactionRequest {
  // (Required) The ID of User.
  int64 user_id = 1 [(validate.rules).int64.gte = 1];
  // (Required) The date and time of the created transaction.
  google.protobuf.Timestamp datetime = 2 [(validate.rules).timestamp.required = true];
  // (Required) The amount of the transaction, should not be 0.
  float amount = 3 [(validate.rules).float = {gte: 0.1, lte: -0.1}];
}
```

To this:

```protobuf
// CreateTransactionRequest
message CreateTransactionRequest {
  // (Required) The ID of User.
  int64 user_id = 1;
  // (Required) The date and time of the created transaction.
  google.protobuf.Timestamp datetime = 2;
  // (Required) The amount of the transaction, should not be 0.
  float amount = 3;
}
```

Then, you can run this `ghz` command to do Load Testing for specific RPC, for the example:

### 1. CreateTransaction RPC:

```sh
ghz --insecure --proto ./api/proto/service.proto --call BTCService.CreateTransaction -d '{ "user_id": 1, "datetime": { "seconds": 1676339196, "nanos": 0 }, "amount": 100 }' 0.0.0.0:8080 -O html -o load_testing_create_transaction.html
```

### 2. ListTransaction RPC:

```sh
ghz --insecure --proto ./api/proto/service.proto --call BTCService.ListTransaction -d '{ "user_id": 1, "start_datetime": { "seconds": 1676339196, "nanos": 0 }, "end_datetime": { "seconds": 1676339196, "nanos": 0 } }' 0.0.0.0:8080 -O html -o load_testing_list_transaction.html
```

### 3. GetUserBalance RPC:

```sh
ghz --insecure --proto ./api/proto/service.proto --call BTCService.GetUserBalance -d '{ "user_id": 1 }' 0.0.0.0:8080 -O html -o load_testing_get_user_balance.html
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
