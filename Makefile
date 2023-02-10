GOPATH=$(shell go env GOPATH)

GOLANGCI_LINT_VERSION=v1.51.1

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
