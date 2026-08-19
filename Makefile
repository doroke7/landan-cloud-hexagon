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



.PHONY: watch-http
watch-http:
	air -c .air.http.toml

.PHONY: watch-facade
watch-facade:
	air -c .air.facade.toml

.PHONY: watch-resource
watch-resource:
	air -c .air.resource.toml
