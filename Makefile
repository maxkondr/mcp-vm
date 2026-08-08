.PHONY: build-server build-client
.DEFAULT_GOAL := help

BIN_PATH ?= $(CURDIR)/bin
BUILD_FLAGS = -v -trimpath -ldflags="-s -w" -gcflags="all=-N -l"

help:	## Display this help message
	@echo "Please use \`make <target>\`, where <target> is one of the following:"
	@awk 'BEGIN {FS = ":.*## "}; /^[a-zA-Z0-9-]+:.*## / {printf "  %-20s%s\n", $$1, $$2}' $(MAKEFILE_LIST)

$(BIN_PATH):
	mkdir -p $(BIN_PATH)

build-server: $(BIN_PATH)	## Build MCP server release binary
	env CGO_ENABLED=0 go build $(BUILD_FLAGS) -o $(BIN_PATH)/ ./cmd/server

run-server:	## Run MCP server and VictoriaMetrics in dev mode (no TLS, no auth)
	docker compose -f ./docker-compose.yml up --build --detach --renew-anon-volumes --remove-orphans --wait --wait-timeout 100

build-client: $(BIN_PATH)	## Build MCP client release binary
	env CGO_ENABLED=0 go build $(BUILD_FLAGS) -o $(BIN_PATH)/ ./cmd/client

