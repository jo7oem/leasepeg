DEV_BIN:=$(shell pwd)/dev_tools/bin
GOBIN=$(DEV_BIN)
export GOBIN
MAKEFILE_DIR := $(shell cd $(dir $(lastword $(MAKEFILE_LIST)))&&pwd )
ABS_DEV_BIN := $(MAKEFILE_DIR)/$(DEV_BIN)

.PHONY: build
build:
	go build -o hatsukari ./...

.PHONY: fmt
fmt: $(DEV_BIN)/golangci-lint  $(DEV_BIN)/yamlfmt
	go fix ./...
	$(DEV_BIN)/golangci-lint run --fix --config=.golangci.yml
	$(DEV_BIN)/yamlfmt .

.PHONY: go-lint
go-lint: $(DEV_BIN)/golangci-lint
	go fix ./...
	$(DEV_BIN)/golangci-lint run --config=.golangci.yml

.PHONY: setup
setup: $(DEV_BIN)/air $(DEV_BIN)/dlv $(DEV_BIN)/golangci-lint $(DEV_BIN)/sqlc $(DEV_BIN)/migrate

$(DEV_BIN)/air:
	mkdir -p $(@D)
	go install github.com/air-verse/air@latest


$(DEV_BIN)/golangci-lint:
	mkdir -p $(@D)
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.13.1

$(DEV_BIN)/dlv:
	mkdir -p $(@D)
	go install github.com/go-delve/delve/cmd/dlv@v1.27.1

$(DEV_BIN)/yamlfmt:
	mkdir -p $(@D)
	go install github.com/google/yamlfmt/cmd/yamlfmt@latest
$(DEV_BIN)/sqlc:
	mkdir -p $(@D)
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
$(DEV_BIN)/migrate:
	mkdir -p $(@D)
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest


.PHONY: ci
ci: lint

.PHONY: lint
lint: go-lint redoc-lint yaml-lint

.PHONY: redoc-lint
redoc-lint:
	podman compose --profile=lint run redoc-cli

.PHONY: yaml-lint
yaml-lint:  $(DEV_BIN)/yamlfmt
	$(DEV_BIN)/yamlfmt -lint .

.PHONY: test
test:
	go test --race -parallel 8 ./...

# Migration targets
MIGRATE_BIN := $(DEV_BIN)/migrate
MIGRATE_CMD := $(MIGRATE_BIN) -path ./sql/postgres/migrations -database "postgres://postgres:postgres@localhost:5432/mydatabase?sslmode=disable"

.PHONY: migrate-up
migrate-up: $(MIGRATE_BIN)
	$(MIGRATE_CMD) up

.PHONY: migrate-down
migrate-down: $(MIGRATE_BIN)
	$(MIGRATE_CMD) down 1

.PHONY: migrate-new
migrate-new: $(MIGRATE_BIN)
	./scripts/create-migration.sh $(MIGRATION_NAME)