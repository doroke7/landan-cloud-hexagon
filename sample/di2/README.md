## 安裝 wire
```
go install github.com/google/wire/cmd/wire@latest
```

## 利用 wire 產出 “container” 代碼
```
wire gen -output container.go
```

## wire 到底在幹什麼？
```
1. 繼承不好，組合比較好
2. 開始代碼架構大批湧向 DI， 而不是繼承
3. Di 的問題是：注入關係鏈大的時候，會很複雜。 如 NewAppContoller(NewAppLogic(NewAppModel(NewAesHelper())))
4. wire 提供一種思路，他提供開發著 平行注入的語法，然後幫你生成這種 嵌套注入語法代碼給 main.go 使用。

```

## java, php, nodejs 用的是哪一種思維
```
1. 他用的 “動態” container 思維，呼叫的時候再去 容器裡面找唯一物件
2. 這種思維go也有實作，叫做 dig
```

## 

## wire 與 dig 有什麼不同
```
1. wire 是靜態代思維（compile先把關係建立好）， dig 是動態思維 （runtime 時候需要再去字典查找）
2. 目前 wire 用的人比較多
``` 

