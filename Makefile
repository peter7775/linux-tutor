APP_NAME := linux-tutor
MAIN_PKG := ./cmd/app
BIN_DIR := bin
BIN := $(BIN_DIR)/$(APP_NAME)

GO := go
GOLANGCI_LINT ?= golangci-lint
GOSEC ?= gosec

GO_FILES := $(shell find . -type f -name '*.go' -not -path './bin/*' -not -path './vendor/*')

.PHONY: all build test lint go-sec fmt fmt-check tidy clean

all: fmt lint go-sec test build

$(BIN_DIR):
	@mkdir -p $(BIN_DIR)

tidy:
	$(GO) mod tidy

fmt:
	gofmt -w $(GO_FILES)

fmt-check:
	@test -z "$$(gofmt -l $(GO_FILES))"

build: fmt tidy
	$(GO) build -o $(BIN) $(MAIN_PKG)

test: fmt
	$(GO) test ./...

lint:
	$(GOLANGCI_LINT) run ./...

go-sec:
	$(GOSEC) ./...

clean:
	rm -rf $(BIN_DIR)