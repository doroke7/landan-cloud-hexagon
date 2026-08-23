#!/bin/bash


# 預設目標
.PHONY: help
help:



.PHONY: protoc
protoc:
	@protoc \
	-I ./proto \
	--go_out=paths=source_relative:./pb \
	--go-grpc_out=paths=source_relative:./pb \
	$(shell find ./proto -name "*.proto")


.PHONY: wire
wire:
	cd container/ && wire



.PHONY: http-watch
http-watch:
	air -c .air.http.toml

.PHONY: facade-watch
facade-watch:
	air -c .air.facade.toml

.PHONY: websocket-watch
websocket-watch:
	air -c .air.websocket.toml

.PHONY: resource-watch
resource-watch:
	air -c .air.resource.toml
