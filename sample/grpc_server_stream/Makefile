# 變數
APP = myapp

# 預設執行
all: build

build:
	protoc --proto_path=./proto \
          --go_out=. --go_opt=paths=source_relative \
          --go-grpc_out=. --go-grpc_opt=paths=source_relative \
          proto/monitor.proto

