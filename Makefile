APP_NAME := linux-tutor
MAIN_PKG := ./cmd/app
BIN_DIR := bin
BIN := $(BIN_DIR)/$(APP_NAME)

GOLANGCI_LINT ?= golangci-lint
GOSEC ?= gosec

.PHONY: all build test lint go-sec clean tidy

all: build

$(BIN_DIR):
	@mkdir -p $(BIN_DIR)

tidy:
	go mod tidy

build: $(BIN_DIR) tidy
	go build -o $(BIN) $(MAIN_PKG)

test:
	go test ./...

lint:
	$(GOLANGCI_LINT) run ./...

go-sec:
	$(GOSEC) ./...

clean:
	rm -rf $(BIN_DIR)