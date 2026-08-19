
## base64 這個 client server 對接的重點
1. client 端記得要把 origin-base64 改成 url-base64 （+ -> -, / -> _ , = -> ）
2. server 端 可以把 收到的 base64 都強轉成 origin-base64，再處理
3. 不論哪一種轉換策略，client 端收到string 之後 都要先把它轉成正常的 base64 後再處理


## 有2種主流策略 避開 base64 url 字元是把
A. 使用 URL-base64， 簡單的說就是把 url 特殊符號的字 用別的字代替
B. 使用 urlencode 思維 ，簡單的說他是為了設計成可以把任何字串 安全放入的格式 （但是需要注意 不同語言的 安全url 不一樣）





## encodeURIComponent 的具體編碼
空格  → %20
+     → %2B
/     → %2F
?     → %3F
=     → %3D
&     → %26
#     → %23