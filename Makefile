DEV_BIN:=$(shell pwd)/dev_tools/bin
GOBIN=$(DEV_BIN)
export GOBIN
MAKEFILE_DIR := $(shell cd $(dir $(lastword $(MAKEFILE_LIST)))&&pwd )
ABS_DEV_BIN := $(MAKEFILE_DIR)/$(DEV_BIN)

.PHONY: build
build:
	go build -o hatsukari ./...

.PHONY: fmt
fmt: $(DEV_BIN)/golangci-lint
	go fix ./...
	$(DEV_BIN)/golangci-lint run --fix --config=.golangci.yml

.PHONY: lint
lint: $(DEV_BIN)/golangci-lint
	$(DEV_BIN)/golangci-lint run --config=.golangci.yml

.PHONY: setup
setup: $(DEV_BIN)/air $(DEV_BIN)/dlv $(DEV_BIN)/golangci-lint

$(DEV_BIN)/air:
	mkdir -p $(@D)
	go install github.com/air-verse/air@latest


$(DEV_BIN)/golangci-lint:
	mkdir -p $(@D)
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.10.1

$(DEV_BIN)/dlv:
	mkdir -p $(@D)
	go install github.com/go-delve/delve/cmd/dlv@v1.26.0

.PHONY: ci
ci: lint

.PHONY: test
test:
	go test --race -parallel 8 ./...