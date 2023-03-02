# BTC Service

---

[![CI Workflow](https://github.com/moemoe89/btc/actions/workflows/ci.yml/badge.svg)](https://github.com/moemoe89/btc/actions/workflows/ci.yml) <!-- start-coverage --><img src="https://img.shields.io/badge/coverage-84.1%25-brightgreen"><!-- end-coverage -->

BTC Service handles BTC transaction and User balance related data.

## Table of Contents

<!-- @import "[TOC]" {cmd="toc" depthFrom=1 depthTo=6 orderedList=false} -->

<!-- code_chunk_output -->

- [Table of Contents](#table-of-contents)
- [Project Summary](#project-summary)
- [Architecture Diagram](#architecture-diagram)
- [Installation](#installation)
    - [1. Set Up Golang Development Environment](#1-set-up-golang-development-environment)
    - [2. Install Development Utility Tools](#2-install-development-utility-tools)
- [Development workflow and guidelines](#development-workflow-and-guidelines)
    - [1. API](#1-api)
    - [2. TimescaleDB + GUI](#2-timescaledb--gui)
    - [3. Migration](#3-migration)
    - [4. Database Schema](#4-database-schema)
    - [5. Cache](#5-cache)
    - [6. Instrumentation](#6-instrumentation)
    - [7. Unit Test](#7-unit-test)
    - [8. Linter](#8-linter)
    - [9. Mock](#9-mock)
    - [10. Run the service](#10-run-the-service)
    - [11. Test the service](#11-test-the-service)
    - [12. Load Testing](#12-load-testing)
    - [13. Messaging](#13-messaging)
- [Project Structure](#project-structure)
- [GitHub Actions CI](#github-actions-ci)
- [Documentation](#documentation)
  - [Visualize Code Diagram](#visualize-code-diagram)
  - [RPC Sequence Diagram](#rpc-sequence-diagram)
- [TODO](#todo)

<!-- /code_chunk_output -->

## Project Summary

| Item                      | Description                                                                                                           |
|---------------------------|-----------------------------------------------------------------------------------------------------------------------|
| Golang Version            | [1.19](https://golang.org/doc/go1.19)                                                                                 |
| Database                  | [TimescaleDB](https://www.timescale.com) and [pgx](https://github.com/jackc/pgx)                                      |
| Database Documentation    | [SchemaSpy](https://schemaspy.org)                                                                                    |
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

## Architecture Diagram

---

[Excalidraw link](https://excalidraw.com/#json=l2lSpGvdr3pzL7VN9__qD,gkfkHNTucXGdPXNNC7DBHw)

![Architecture-Diagram](https://user-images.githubusercontent.com/7221739/222329857-93e39281-7440-476c-ab97-40f561728539.png)


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

#### NOTE:

If you have any difficulties installing all dependencies needed for generating the proto files,
you can easily build the docker image and use that instead.

Here are 2 commands for building and generating:

```shell
make build-protoc
make docker-protoc
```

### 2. TimescaleDB + GUI

![pgAdmin](https://user-images.githubusercontent.com/7221739/222328956-6db95f91-cffc-43fe-8e31-6c051c11f8e5.png)

#### NOTE
In this project, the Database Replication could be implemented, then we need 2 databases Master and Slave.
But if only 1 database exists, we can easily the replication on the App side by setting the env variables.

From `IS_REPLICA` `true` to `false`.

Run TimescaleDB locally with the GUI (pgAdmin) can be executed with the following docker-compose command:

```sh
$ docker-compose -f ./development/docker-compose.yml up timescaledb-master
$ docker-compose -f ./development/docker-compose.yml up timescaledb-slave
$ docker-compose -f ./development/docker-compose.yml up pgadmin
```

> NOTE:
> TimescaleDB will use port 5432 (Master) and 5433 (Slave) and pgAdmin will use 5050, please make sure those port are unused in your system.
> If the port conflicted, you can change the port on the [development/docker-compose.yml](development/docker-compose.yml) file.

The default email & password for pgAdmin are:
* email: `admin@admin.com` 
* password: `admin123`

With this following TimescaleDB Master info:
* host: `timescaledb-master` -> change this to `localhost` if you try to connect from outside Docker
* port: `5432`
* username: `test`
* password: `test`
* db: `test`

And this is the info for TimescaleDB Slave:
* host: `timescaledb-slave` -> change this to `localhost` if you try to connect from outside Docker
* port: `5433`
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

### 4. Database Schema

![SchemaSpy](https://user-images.githubusercontent.com/7221739/222328524-7b8178dd-1acc-4093-9e00-12d35d4c5a6c.png)

If you want to check the overall Database Schema, you can use a UI tool based on browser
using SchemaSpy.

The docker-compose for SchemaSpy already exist, but make sure to run the `tiemscaledb-master` and do the `migration`, thne you can just run this command:

```shell
$ docker-compose -f ./development/docker-compose.yml up schemaspy
```

The HTML and assets file will be generated under `development/schemaspy/output` directory.

### 5. Cache

When getting transactions list and user balance, there's a cache implemented using Redis
in order to have middle layer and avoid call the main DB frequently.

To start running Redis, there's a docker-compose command available:

```sh
$ docker-compose -f ./development/docker-compose.yml up redis
```

### 6. Instrumentation

![Jaeger](https://user-images.githubusercontent.com/7221739/222329540-55f8c982-becd-43d5-a4a7-fca8661f1c25.png)

This service implements [https://opentelemetry.io/](https://opentelemetry.io/) to enable instrumentation in order to measure the performance.
The data exported to Jaeger and can be seen in the Jaeger UI [http://localhost:16686](http://localhost:16686)

For running the Jaeger exporter, easily run with docker-compose command:

```sh
$ docker-compose -f ./development/docker-compose.yml up jaeger
```

### 7. Unit Test

Make sure the database already running, then you can simply execute the following command to run all test cases in this service:

```sh
$ make test
```

### 8. Linter

For running the linter make sure these libraries already installed in your system:

* https://github.com/golangci/golangci-lint
* https://github.com/yoheimuta/protolint

Then checks the Go and Proto code style using lint can be done with this command:

```sh
$ make lint
```

### 9. Mock

This service using Mock in some places like in the repository, usecase, pkg, etc.
To automatically updating the mock if the interface changed, easily run with `go generate` command:

```sh
$ make mock
```

### 10. Run the service

For running the service, you need the database running and set up some env variables:

```
# app config
export APP_ENV=dev
export SERVER_PORT=8080

# master db config
export POSTGRES_USER_MASTER=test
export POSTGRES_PASSWORD_MASTER=test
export POSTGRES_HOST_MASTER=localhost
export POSTGRES_PORT_MASTER=5432
export POSTGRES_DB_MASTER=test

# slave db config
export POSTGRES_USER_SLAVE=test
export POSTGRES_PASSWORD_SLAVE=test
export POSTGRES_HOST_SLAVE=localhost
export POSTGRES_PORT_SLAVE=5433
export POSTGRES_DB_SLAVE=test

# use replica config
export IS_REPLICA=true

# tracing config
export OTEL_AGENT=http://localhost:14268/api/traces

# cache config
export REDIS_HOST=localhost:6379
```

Or you can just execute the sh file:

```sh
$ ./scripts/run.sh
```

### 11. Test the service

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

![BloomRPC](https://user-images.githubusercontent.com/7221739/222319521-2bb079a9-ff78-43b8-902f-705d7f816a20.png)

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

![Postman](https://user-images.githubusercontent.com/7221739/222329685-5a1c7499-11c3-4985-9f7b-0fd504da341a.png)

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

![Swagger](https://user-images.githubusercontent.com/7221739/222329119-dad08930-4878-4f4e-a976-31ffa91f863e.png)

This service has HTTP server built on gRPC-Gateway, if you prefer to test using HTTP instead HTTP2 protocol,
you can copy the Swagger file here [api/openapiv2/proto/service.swagger.json](api/openapiv2/proto/service.swagger.json) and then copy paste to this URL https://editor.swagger.io/

By default, HTTP server running on gRPC port + 1, if the gRPC port is 8080, then HTTP server will run on 8081.

### 12. Load Testing

![ghz](https://user-images.githubusercontent.com/7221739/222329410-c29564da-e4ca-4870-b0d0-ecccfdcf4593.png)

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

### 13. Messaging

In order to avoid failing when creates the transaction and support for easily retry,
there's a simple Event based system using RabbitMQ.

To test the event based you need to run the rabbitmq, the server and the consumer server.

```shell
$ docker-compose -f ./development/docker-compose.yml up timescaledb-master timescaledb-slave pgadmin jaeger rabbitmq
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
> Then you will have all services running like `timescaledb-master`, `timescaledb-slave`, `pgadmin`, `jaeger`, `rabbitmq`, `redis`
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

![GitHubActionsCI](https://user-images.githubusercontent.com/7221739/222317570-2f8098e1-4a66-4f77-bac4-f7628a8dec0b.png)

This project has GitHub Actions CI to do some automation such as:

* [lint](.github/workflows/lint.yml): check the code style.
* [test](.github/workflows/test.yml): run unit testing and uploaded code coverage artifact.
* [generate-proto](.github/workflows/generate-proto.yml): generates protobuf files.
* [generate-rpc-diagram](.github/workflows/generate-rpc-diagram.yml): generates RPC sequence diagram.
* [generate-diagram](.github/workflows/generate-diagram.yml): generates graph code visualization.
* [push-file](.github/workflows/push-file.yml): commit and push generated proto, diagram as github-actions[bot] user.

## Documentation

### Visualize Code Diagram

![GraphDiagram](https://user-images.githubusercontent.com/7221739/222317676-b4b33482-2ddb-459e-8cbe-97fdf93cd789.png)

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

![SequenceDiagram](https://user-images.githubusercontent.com/7221739/222317786-aaf2b078-b26b-43cb-98f8-1599f12b46c0.png)

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

## TODO

---

In the future, here is the Architecture Diagram want to achieve to improve the performance.

* Basically need to separate the service into Create & Search service, following CQRS pattern and using event based for the communication.

* Create service will have responsibility only for inserting the data and trigger message to Search service. To create the transaction need to publish an event from client.

* Search service will have responsibility regarding searching data such as search transactions and get user balance. Client can call this service directly using gRPC or REST (gRPC-Gateway).
Search service also do indexing the transaction data into search engine service like Elasticsearch, triggered from Create service by event.
For the user balance still getting from replica database, and both transaction and user balance will have a middle layer cache using Redis. Just in case connection to Elasticsearch fail, the search service still be able to get the data from replica database directly.

[Excalidraw](https://excalidraw.com/#json=VZfAhQiwJ7FpvFutyrU4t,Uoutk-t-s-NbcGOzxid-WQ)

![Future-Architecture-Diagram](https://user-images.githubusercontent.com/7221739/222329982-b963cab4-9735-4054-9c7b-0b1dfa4984cc.png)
