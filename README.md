# BTC Service

---

[![CI Workflow](https://github.com/moemoe89/btc/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/moemoe89/btc/actions/workflows/ci.yml?query=workflow%3Atest) <!-- start-coverage --><img src="https://img.shields.io/badge/coverage-83.6%25-yellowgreen"><!-- end-coverage -->

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
    - [4. Cache](#4-cache)
    - [5. Instrumentation](#5-instrumentation)
    - [6. Unit Test](#6-unit-test)
    - [7. Linter](#7-linter)
    - [8. Mock](#8-mock)
    - [9. Run the service](#9-run-the-service)
    - [10. Test the service](#10-test-the-service)
    - [11. Load Testing](#11-load-testing)
    - [12. Messaging](#12-messaging)
- [Project Structure](#project-structure)
- [GitHub Actions CI](#github-actions-ci)
- [Documentation](#documentation)
  - [Visualize Code Diagram](#visualize-code-diagram)
  - [RPC Sequence Diagram](#rpc-sequence-diagram)

<!-- /code_chunk_output -->

## Project Summary

| Item                      | Description                                                                                                           |
|---------------------------|-----------------------------------------------------------------------------------------------------------------------|
| Golang Version            | [1.19](https://golang.org/doc/go1.19)                                                                                 |
| Database                  | [TimescaleDB](https://www.timescale.com) and [pgx](https://github.com/jackc/pgx)                                      |
| Cache                     | [Redis](https://redis.com) and [go-redis](https://github.com/redis/go-redis)                                          |
| Migration                 | [migrate](https://github.com/golang-migrate/migrate)                                                                  |
| moq                       | [mockgen](https://github.com/golang/mock)                                                                             |
| Linter                    | [GolangCI-Lint](https://github.com/golangci/golangci-lint)                                                            |
| Testing                   | [testing](https://golang.org/pkg/testing/) and [testify/assert](https://godoc.org/github.com/stretchr/testify/assert) |
| Load Testing              | [ghz](https://ghz.sh)                                                                                                 |
| API                       | [gRPC](https://grpc.io/docs/tutorials/basic/go/) and [gRPC-Gateway](https://github.com/grpc-ecosystem/grpc-gateway)   |
| Application Architecture  | [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)                    |
| Directory Structure       | [Standard Go Project Layout](https://github.com/golang-standards/project-layout)                                      |
| CI (Lint, Test, Generate) | [GitHubActions](https://github.com/features/actions)                                                                  |
| Visualize Code Diagram    | [go-callviz](https://github.com/ofabry/go-callvis)                                                                    |
| Sequence Diagram          | [Mermaid](https://mermaid.js.org)                                                                                     |
| Protobuf Operations       | [buf](https://buf.build)                                                                                              |
| Instrumentation           | [OpenTelemetry](https://opentelemetry.io) and [Jaeger](https://www.jaegertracing.io)                                  |
| Logger                    | [zap](https://github.com/uber-go/zap)                                                                                 |
| Messaging                 | [RabbitMQ](https://www.rabbitmq.com) and [amqp091-go](https://github.com/rabbitmq/amqp091-go)                         |

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

This service also implementing gRPC-Gateway with this library:

* https://github.com/grpc-ecosystem/grpc-gateway

For generating the gRPC-Gateway and OpenAPI files, there's required additional package such as:

* github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
* github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2

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
> TimescaleDB will use port 5432 and pgAdmin will use 5050, please make sure those port are unused in your system.
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

This migration also seeds some test data, because when creating a transaction, will require existing User ID.
By this seeds, we will have 5 users test data, from ID 1 to 5.

### 4. Cache

When getting transactions list and user balance, there's a cache implemented using Redis
in order to have middle layer and avoid call the main DB frequently.

To start runing Redis, there's a docker-compose command available:

```sh
$ docker-compose -f ./development/docker-compose.yml up redis
```

### 5. Instrumentation

This service implements [https://opentelemetry.io/](https://opentelemetry.io/) to enable instrumentation in order to measure the performance.
The data exported to Jaeger and can be seen in the Jaeger UI [http://localhost:16686](http://localhost:16686)

For running the Jaeger exporter, easily run with docker-compose command:

```sh
$ docker-compose -f ./development/docker-compose.yml up jaeger
```

### 6. Unit Test

Make sure the database already running, then you can simply execute the following command to run all test cases in this service:

```sh
$ make test
```

### 7. Linter

For running the linter make sure these libraries already installed in your system:

* https://github.com/golangci/golangci-lint
* https://github.com/yoheimuta/protolint

Then checks the Go and Proto code style using lint can be done with this command:

```sh
$ make lint
```

### 8. Mock

This service using Mock in some places like in the repository, usecase, pkg, etc.
To automatically updating the mock if the interface changed, easily run with `go generate` command:

```sh
$ make mock
```

### 9. Run the service

For running the service, you need the database running and set up some env variables:

```
export APP_ENV=dev
export SERVER_PORT=8080
export POSTGRES_USER=test
export POSTGRES_PASSWORD=test
export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432
export POSTGRES_DB=test
export OTEL_AGENT=http://localhost:14268/api/traces
export REDIS_HOST=localhost:6379
```

Or you can just execute the sh file:

```sh
$ ./scripts/run.sh
```

### 10. Test the service

The example how to call the gRPC service written in Golang can be seen on this [example-client](scripts/example-client) file.

> NOTE: To test this service need the migration to be done. After that you can choose the User ID's from 1 to 5.
 
If you want to test by GUI client, you can use either BloomRPC (although already no longer active) or Postman.
For the detail please visit these links:
* https://github.com/bloomrpc/bloomrpc
* https://www.postman.com

Basically you just need to import the [api/proto/service.proto](api/proto/service.proto) file if you want to test via BloomRPC / Postman.

> NOTE: There will be a possibility issue when importing the proto file to BloomRPC or Postman.
> It is caused by some path issue, the usage of `gRPC Gateway` and `protoc-gen-validate` library.
> To solve this issue, there's need a modification for the proto file.

#### BloomRPC

BloomRPC will have these issues when trying to import the proto file:

```
Error while importing protos
illegal name ';' (/path/btc/api/proto/service.proto, line 20)

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

and remove gRPC Gateway related annotations:

```protobuf
option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
  ...
};
```

#### Postman

There's some issue when importing to Postman. Basically we need to do the same things like BloomRPC and disable the validate import.

```protobuf
import "proto/entity.proto";
import "validate/validate.proto";
import "protoc-gen-openapiv2/options/annotations.proto";
```

To this:

```protobuf
import "../proto/entity.proto";
// import "validate/validate.proto";
// import "protoc-gen-openapiv2/options/annotations.proto";
```

Also don't forget to set the import path e.g. `{YOUR-DIR}/btc/api/proto`

#### gRPC-Gateway

This service has HTTP server built on gRPC-Gateway, if you prefer to test using HTTP instead HTTP2 protocol,
you can copy the Swagger file here [api/openapiv2/proto/service.swagger.json](api/openapiv2/proto/service.swagger.json) and then copy paste to this URL https://editor.swagger.io/

By default HTTP server running on gRPC port + 1, if the gRPC port is 8080, then HTTP server will run on 8081.

### 11. Load Testing

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
import "google/api/annotations.proto";
import "protoc-gen-openapiv2/options/annotations.proto";
```

To this:

```protobuf
import "../proto/entity.proto";
// import "validate/validate.proto";
// import "google/api/annotations.proto";
// import "protoc-gen-openapiv2/options/annotations.proto";
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

and remove gRPC Gateway related annotations:

```protobuf
option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
  ...
};
```

Then, you can run this `ghz` command to do Load Testing for specific RPC, for the example:

#### 1. CreateTransaction RPC:

```sh
ghz --insecure --proto ./api/proto/service.proto --call BTCService.CreateTransaction -d '{ "user_id": 1, "datetime": { "seconds": 1676339196, "nanos": 0 }, "amount": 100 }' 0.0.0.0:8080 -O html -o load_testing_create_transaction.html
```

#### 2. ListTransaction RPC:

```sh
ghz --insecure --proto ./api/proto/service.proto --call BTCService.ListTransaction -d '{ "user_id": 1, "start_datetime": { "seconds": 1676339196, "nanos": 0 }, "end_datetime": { "seconds": 1676339196, "nanos": 0 } }' 0.0.0.0:8080 -O html -o load_testing_list_transaction.html
```

#### 3. GetUserBalance RPC:

```sh
ghz --insecure --proto ./api/proto/service.proto --call BTCService.GetUserBalance -d '{ "user_id": 1 }' 0.0.0.0:8080 -O html -o load_testing_get_user_balance.html
```

### 12. Messaging

In order to avoid failing when creates the transaction and support for easily retry,
there's a simple Event based system using RabbitMQ.

To test the event based you need to run the rabbitmq, the server and the consumer server.

```shell
$ docker-compose -f ./development/docker-compose.yml up timescaledb pgadmin jaeger rabbitmq
$ ./scripts/run.sh
$ ./scripts/run-consumer.sh
```

After that you can try to send a message by publishing a message.

```shell
go run ./scripts/example-publish
```

# NOTE

> If you have any difficulties to run the service, easily just run all dependencies by docker-compose for the example:
> 
> `docker-compose -f ./development/docker-compose.yml up`
>
> Then you will have all services running like `timescaledb`, `pgadmin`, `jaeger`, `rabbitmq`, `redis`
> also running the `migration` and run `btc-server` + `btc-consumer`.

## Project Structure

This project follow https://github.com/golang-standards/project-layout

However, for have a clear direction when working in this project, here are some small guide about each directory:

* [api](api): contains Protobuf files, generated protobuf, swagger, etc.
* [build](build): Docker file for the service, migration, etc.
* [cmd](cmd): main Go file for running the service, producer, consumer, etc.
* [development](development): file to support development like docker-compose.
* [docs](docs): file about project documentations such as diagram, sequence diagram, etc.
* [internal](internal): internal code that can't be shared.
  * [internal/adapters/grpchandler](internal/adapters/grpchandler): adapter layer that serve into gRPC service.
  * [internal/di](internal/di): dependencies injection for connecting each layer.
  * [internal/entities/repository](internal/entities/repository): data entities to connect with repository layer.
  * [internal/infrastructure](internal/infrastructure): codes that doing database operations.
  * [internal/usecases](internal/usecases): business logic that connect to repository layer, RPC & HTTP client, etc.
* [migrations](migrations): database migration files.
* [pkg](pkg): package code that can be shared.
* [scripts](scripts): shell script, go script to help build or testing something.
* [tools](tools): package that need to store on go.mod in order to easily do installation.

## GitHub Actions CI

This project has GitHub Actions CI to do some automation such as:

* [lint](.github/workflows/lint.yml): check the code style.
* [test](.github/workflows/test.yml): run unit testing and uploaded code coverage artifact.
* [generate-proto](.github/workflows/generate-proto.yml): generates protobuf files.
* [generate-rpc-diagram](.github/workflows/generate-rpc-diagram.yml): generates RPC sequence diagram.
* [generate-diagram](.github/workflows/generate-diagram.yml): generates graph code visualization.
* [push-file](.github/workflows/push-file.yml): commit and push generated proto, diagram as github-actions[bot] user.

## Documentation

### Visualize Code Diagram

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

### RPC Sequence Diagram

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
