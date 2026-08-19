## 為何需要學習拆包黏包
1. 因為對於大型遊戲，帶寬成本極高，需要自行定義 傳輸格式
1. 比 http2-grpc 省掉 HTTP/2 Frame Header （定義 http2 協議的 表頭）
2. 省 http1.1 省掉 HTTP Header （http 應用資料的表頭）