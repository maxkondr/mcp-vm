.PHONY: build-server build-client
.DEFAULT_GOAL := help

BIN_PATH ?= $(CURDIR)/../bin

help:	## Display this help message
	@echo "Please use \`make <target>\`, where <target> is one of the following:"
	@awk 'BEGIN {FS = ":.*## "}; /^[a-zA-Z0-9-]+:.*## / {printf "  %-20s%s\n", $$1, $$2}' $(MAKEFILE_LIST)

$(BIN_PATH):
	mkdir -p $(BIN_PATH)

build-server: $(BIN_PATH)	## Build MCP server release binary
	env CGO_ENABLED=0 go build -v -gcflags="all=-N -l" -o $(BIN_PATH)/ ./cmd/server

build-client: $(BIN_PATH)	## Build MCP client release binary
	env CGO_ENABLED=0 go build -v -gcflags="all=-N -l" -o $(BIN_PATH)/ ./cmd/client

