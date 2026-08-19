## 
1. keepalive, backoff, conn.Connect() 三個缺一不可
2. keepalive：是當連線還是正常狀態的時候，配置間隔多久檢查一次
3. backoff: 是當剛剛進入斷線的時候，設定多久要重新連線
4. conn.Connect：是最基本的主動業務連線
5. 為什麼 在 沒有開啟 conn.Connect 時候，backoff 不會自動連線？
   因為 backoff 在斷線後首次重試後發現連線不上，第二次重試的時候會檢查有沒有業務需要他，沒有就進入 IDLE 不在 backoff