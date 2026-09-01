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
http-watch:                             # 用 exec 會直接執行 air 而不是透過 sh
	exec air -c .air.http.toml       

.PHONY: facade-watch
facade-watch:                             # 用 exec 會直接執行 air 而不是透過 sh
	exec air -c .air.facade.toml

.PHONY: websocket-watch
websocket-watch:                             # 用 exec 會直接執行 air 而不是透過 sh
	exec air -c .air.websocket.toml


.PHONY: centrifuge-watch
centrifuge-watch:                             # 用 exec 會直接執行 air 而不是透過 sh
	exec air -c .air.centrifuge.toml

.PHONY: socketio-watch
socketio-watch:                             # 用 exec 會直接執行 air 而不是透過 sh
	exec air -c .air.socketio.toml

.PHONY: resource-watch
resource-watch:                             # 用 exec 會直接執行 air 而不是透過 sh
	exec air -c .air.resource.toml


