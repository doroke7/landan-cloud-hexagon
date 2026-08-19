#. client stream 相當 客戶端 傳資料是傳輸 一點一點連續流資料過去， 全部資料傳遞完成， 才響應


# 好處
1. 好處是 不用每一次資料 都response， 不用一直收 request header
```goregexp

protoc --proto_path=./proto \
    --go_out=./protobuf --go_opt=paths=source_relative \
    --go-grpc_out=./protobuf --go-grpc_opt=paths=source_relative \
    ./proto/monitor.proto
```