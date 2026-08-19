## 安裝 tar to go 套件
```sh
go install github.com/TarsCloud/TarsGo/tars/tools/tars2go@latest
```

## 生成 tar go 文件
```sh
tars2go --outdir=./protocol $(find ./tars -name "*.tars")
```


## 運行 tars rpc server
```sh
tars2go --outdir=./protocol $(find ./tars -name "*.tars")
```

## tar 比較老 
```
看起來是有一點笨，缺少了 約定勝於配置思維，其實 tars 的 結構，路由，還有 config.xml 的 servant 屬於同一個事情，卻需要分別配置
```
