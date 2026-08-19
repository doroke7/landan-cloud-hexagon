
# A-1 更新索引並安裝 protobuf 主程式
```
apk add --no-cache protobuf protobuf-dev
```

# A-2 撰寫 proto/greeter.proto 文件

# A-3 產生 基本go程序
```
mkdir pb
protoc --go_out=. --go-grpc_out=. proto/greeter.proto


# 當你寫 --go_out=. 時，protoc 會根據 .proto 檔案裡的 option go_package 來決定最終檔案產生的位置。
```

# B-1 撰寫 service 的 greeter_service.go

# C-1 撰寫 gRPC server，並且把 greeter_service.go 註冊進來

# D-1 撰寫 gRPC client。






# 🚀 gRPC 通訊模式快速參考

---

### 1. 簡單模式 (Unary RPC)
* **模式：** 一問一答 (One Request, One Response)
* **流程：** `Client ──▶ Request ──▶ Server ──▶ Response ──▶ Client`
* **場景：** 基礎 API、登入、獲取單個用戶資料。
* **特性：** 最簡單，類似傳統 HTTP 調用。

---

### 2. 伺服器串流 (Server Streaming)
* **模式：** 一發多收 (One Request, Multiple Responses)
* **流程：** `Client ──▶ Request ──▶ Server ──▶ [Data 1, Data 2, Data N...] ──▶ Client`
* **場景：** 大數據導出（如百萬級爬蟲紀錄）、股票行情實時推送。
* **特性：** **「記憶體救星」**，Server 不用在 RAM 湊齊數據，查一筆噴一筆。

---

### 3. 客戶端串流 (Client Streaming)
* **模式：** 多發一收 (Multiple Requests, One Response)
* **流程：** `Client ──▶ [Data 1, Data 2, Data N...] ──▶ Server ──▶ Response ──▶ Client`
* **場景：** 大檔案分塊上傳、批量數據匯入系統。
* **特性：** 客戶端持續餵數據，Server 邊收邊處理，最後才回結論。

---

### 4. 雙向串流 (Bidirectional Streaming)
* **模式：** 全雙工 (Multiple Requests, Multiple Responses)
* **流程：** `Client ◀───(雙向非同步數據流)───▶ Server`
* **場景：** 即時通訊 (IM)、分散式計算節點之間的狀態同步。
* **特性：** 效能最高、延遲最低，雙方完全對等，不需等待對方回應。

---

> **💡 架構提示：** > 當你從 **ThinkPHP** 這種「請求即終止」的思維轉向 **Go + gRPC** 時，
> 最核心的轉變就是學會利用 **Streaming** 來控制記憶體（零拷貝思維）。