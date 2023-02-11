GOPATH=$(shell go env GOPATH)

GOLANGCI_LINT_VERSION=v1.51.1

install: install-gomock install-linter

install-gomock:
	@echo "\n>>> Install gomock\n"
	go install github.com/golang/mock/mockgen

install-linter:
	@echo "\n>>> Install GolangCI-Lint"
	@echo ">>> https://github.com/golangci/golangci-lint/releases \n"
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/${GOLANGCI_LINT_VERSION}/install.sh | \
	sh -s -- -b ${GOPATH}/bin ${GOLANGCI_LINT_VERSION}

	@echo "\n>>> Install protolint"
	@echo ">>> https://github.com/yoheimuta/protolint/releases \n"
	@go install github.com/yoheimuta/protolint/cmd/protolint

lint:
	@echo "\n>>> Run GolangCI-Lint\n"
	/bin/bash ./scripts/lint.sh

	@echo "\n>>> Run Proto-Lint\n"
	protolint api/proto

test:
	mkdir -p .coverage/html
	go test -v -race -cover -coverprofile=.coverage/internal.coverage.tmp ./internal/... && \
	cat .coverage/internal.coverage.tmp | grep -v "_mock.go\|_mockgen.go" > .coverage/internal.coverage && \
	go tool cover -html=.coverage/internal.coverage -o .coverage/html/internal.coverage.html;

diagram:
	go run ./scripts/generate-diagram
