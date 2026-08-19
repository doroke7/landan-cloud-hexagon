// 定義生成的 Go 包名與路徑
namespace go app

service Hello {
    // 定義一個方法：傳入 sReq (字串)，回傳字串
    // 注意：Thrift 強制要求參數必須有編號（例如 1:）
    string testHello(1: string sReq)
}