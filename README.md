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
- [Running Unit Tests Locally](#running-unit-tests-locally)
    - [1. Run Unit Tests](#2-run-unit-tests)

<!-- /code_chunk_output -->

## Project Summary

| Item                       | Description                                                                                                           |
|----------------------------|-----------------------------------------------------------------------------------------------------------------------|
| Golang Version             | [1.19](https://golang.org/doc/go1.19)                                                                                 |
| moq                        | [mockgen](https://github.com/golang/mock)                                                                             |
| Linter                     | [GolangCI-Lint](https://github.com/golangci/golangci-lint)                                                            |
| Testing                    | [testing](https://golang.org/pkg/testing/) and [testify/assert](https://godoc.org/github.com/stretchr/testify/assert) |
| API                        | [gRPC](https://grpc.io/docs/tutorials/basic/go/)                                                                      |
| Application Architecture   | [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)                    |
| Directory Structure        | [Standard Go Project Layout](https://github.com/golang-standards/project-layout)                                      |
| CI (Lint & Test)           | [GitHubActions](https://github.com/features/actions)                                                                  |

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

### 1. Run Unit Tests

You can simply execute the following command to run all test cases in this service:

```sh
make test
```

### 1. Run Linter

Check the Go and Proto code style using lint can be done with this command:

```sh
make lint
```
