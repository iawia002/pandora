GOPATH ?= $(shell go env GOPATH)
BIN_DIR := $(GOPATH)/bin
GOLANGCI_LINT := $(BIN_DIR)/golangci-lint

.PHONY: lint

lint: $(GOLANGCI_LINT)
	@$(GOLANGCI_LINT) run

$(GOLANGCI_LINT):
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(BIN_DIR) v1.46.2
