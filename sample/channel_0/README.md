| # | Pattern |
|---|---------|
| 1 | Producer / Consumer |
| 2 | Fan-Out (Split Channel)|
| 3 | Fan-In (Merge Channel)|
| 4 | Pipeline |
| 5 | Worker Pool |
| 6 | Semaphore |
| 7 | Timeout |
| 8 | Context Cancel |               + 取消
| 9 | Broadcast (`close`) |          + close 做廣播
| 10 | Request-Response |            + 把 response 綁定到 request 上， client 就能異步解析
| 11 | Future-Promise |              + 基本上就是 模仿await 語法
| 12 | Rate Limiter |                + 限流
| 13 | Boker（Pub-Sub）(Event-bus) |  + 支持 一個訊息同時給多個訂閱·（channel 只能隨機送個單一）
| 14 | Actor |                       + 指 Acotor 會自己動作，外面只要給他訊息就可以
| 15 | Tee Channel |
| 16 | Bridge Channel |              + bridge channel 基本上就是攤平 channel
| 17 | Or-Done |                     + 優化 for-range 跟 select 的語法糖
| 18 | Heartbeat |
| 19 | Graceful Shutdown |           + 利用close 把channel 給所有 worker ，然後worker 監控 close
| 20 | Bounded Queue |                 Go 的 channel 本身就是一種 Bounded Queue（有界佇列）。
| 21 | Unbounded Queue |               無限大小的 queue
| 22 | Batch |
| 23 | Scatter/Gather |              + 本質就是 fade-out + fade-in 聯合使用
| 24 | Retry Queue |                 + 把錯誤的任務丟到另外一個隊列，適當時機再放會原隊列
| 25 | Workers (Load-Balance) |      + 多 Worker 消費同一個 Channel
| 26 | Token Ring |                  + 利用 beforeToken 啟動 本goruntine， 用 afterToken 觸發後一個 goruntine 
| 27 | Barrier |                     + 門
| 28 | Lock-Free Queue |               簡單的說 go-channel 已經實現 lock-free queue, 不然用array 做會需要上鎖
| 29 | Backpressure |                  go-channel 利用 特定個數的channel 已經實現了 減低壓 queue
| 30 | Cancellation Tree |
| 31 | Dynamic Worker |              + 可變動的 worker 來 拉取 channel
| 32 | Channel Ownership |           + 就是只有生產 channel 資料的程才能關閉
| 33 | Queue Drain |                   關閉之前先把 channel 資料 特殊處理
| 34 | Nil Channel Disable |         + nil channel 可以動態阻塞其他channel
| 35 | Select Multiplex |            + channels 的多路復用
| 36 | Sleeping Barber |             * 經典併發問題
| 47 | Dining Philosophers |         * 經典併發問題
