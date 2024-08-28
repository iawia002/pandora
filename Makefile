GOPATH ?= $(shell go env GOPATH)
BIN_DIR := $(GOPATH)/bin
GOLANGCI_LINT := $(BIN_DIR)/golangci-lint

.PHONY: lint test

lint: $(GOLANGCI_LINT)
	@$(GOLANGCI_LINT) run

$(GOLANGCI_LINT):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(BIN_DIR) v1.60.3

test:
	@go test ./... -coverprofile=coverage.out
	@go tool cover -func coverage.out | tail -n 1 | awk '{ print "total: " $$3 }'

IMG ?= iawia002/test:latest

image:
	docker build -f build/Dockerfile . -t ${IMG}

image-push:
	docker buildx build --platform linux/amd64,linux/arm64 -f build/Dockerfile . -t ${IMG} --push
