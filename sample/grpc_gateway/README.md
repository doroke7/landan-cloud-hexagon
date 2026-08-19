# 沒錯，它就是標準的 Adapter Pattern（轉接器模式）。在軟體架構的層次上，它更像是一個 Anti-Corruption Layer（抗腐蝕層），
# 確保你的核心服務不需要為了遷就前端或外部客戶端，而把原本純淨的 gRPC 邏輯弄髒。









# 編譯「轉接合約」 - 執行 protoc 指令。這一步會同時產出 核心通訊碼 與 適配器代碼。
```
# 在專案根目錄執行
# 1. -I ./proto: 這讓 protoc 把 proto/ 當成搜尋根目錄。當你最後寫 service.proto 時，它會從這裡找到檔案，且因為在它眼裡 service.proto 就在根部，所以產出時不會建立 proto/ 資料夾。

# 2. -I .: 這讓 protoc 也能從專案根目錄找東西。當它解析 service.proto 看到 import "google/api/annotations.proto" 時，它會透過這個路徑找到你放在根目錄的 google/ 資料夾。

# 3. 輸入參數 service.proto: 注意！這裡不要寫 proto/service.proto，否則路徑又會被加上去。

# 4. --go_out=pb --go_opt=paths=source_relative
   輸出到 pb 目錄, 並且按照 proto 文件的文件相對目錄結構建立文件

# 5. --go-grpc_out=pb --go-grpc_opt=paths=source_relative 同上

# 6. --grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative 同上


protoc \
  -I ./proto \
  -I . \
  --go_out=pb --go_opt=paths=source_relative \
  --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
  --grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
  service.proto
```

# 客戶端請求 http ，並且代理到 gRPC server
```goregexp
curl -X POST http://localhost:8080/v1/translate \
     -H "Content-Type: application/json" \
     -d '{"text": "你好，適配器"}'

```