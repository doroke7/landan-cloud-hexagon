package main

import "sync/atomic"

func main() {

}

func old_function() {
	var count int64

	atomic.AddInt64(&count, 1)
}

func new_function() {
	var count atomic.Int64

	count.Add(1)
}

/*



+----+--------------------------+---------------------------+----------------------------------------+
| #  | 功能                     | 用途                      | 範例                                   |
+----+--------------------------+---------------------------+----------------------------------------+
| 1  | Counter                  | 原子遞增/遞減             | 線上人數、QPS、TPS、API 次數           |
| 2  | Flag                     | 儲存布林或狀態            | Shutdown、Pause、Maintenance           |
| 3  | State Machine            | 管理物件生命週期          | Init → Running → Stopping → Stopped    |
| 4  | Compare And Swap (CAS)   | 條件更新                  | Singleton、Lock-Free 演算法            |
| 5  | Swap                     | 原子交換值                | Config 切換、物件切換                  |
| 6  | Pointer                  | 原子替換指標              | Config Hot Reload、Router Hot Reload   |
| 7  | Value                    | 原子更新任意物件          | Config、Cache                          |
| 8  | Reference Count          | 引用計數                  | Shared Object、資源管理                |
| 9  | Lock-Free Queue          | 無鎖佇列                  | 高效能 Message Queue                   |
| 10 | Lock-Free Stack          | 無鎖堆疊                  | Treiber Stack                          |
| 11 | Ring Buffer              | 無鎖環形緩衝區            | Log、Network、MQ                       |
| 12 | Concurrent Map           | 無鎖雜湊表                | 高效能 Concurrent Map                  |
| 13 | Spin Lock                | 簡易自旋鎖                | 極短臨界區                             |
| 14 | Mutex（底層）            | 鎖狀態管理                | sync.Mutex                             |
| 15 | RWMutex（底層）          | 讀寫鎖                    | sync.RWMutex                           |
| 16 | Once（底層）             | 單次初始化                | sync.Once                              |
| 17 | Double Checked Locking   | Lazy Initialization       | 降低 Lock 開銷                         |
| 18 | Rate Limiter             | 限制併發                  | 最大連線數、API 限流                   |
| 19 | Statistics               | 收集統計                  | Success、Error、Latency                |
| 20 | Feature Flag             | 功能開關                  | A/B Test、灰度發布                     |
| 21 | Hot Reload               | 即時切換資料              | Config、Routing Table                  |
+----+--------------------------+---------------------------+----------------------------------------+
*/
