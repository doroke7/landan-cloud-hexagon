
================================================================================
          GO 語言【依賴注入 / 參數傳遞】STRUCT 注入 vs POINTER 注入
================================================================================

一、 核心特性的純文字對比表
--------------------------------------------------------------------------------
特性 / 維度      | 結構體注入 (struct)             | 指標注入 (*struct)             
--------------------------------------------------------------------------------
底層傳遞本質     | 傳值 (Value)，記憶體完全複製影本  | 傳引用 (Reference)，傳遞 8 Byte 地址
資料共享性       | 獨享。各服務拿到的零件完全獨立  | 共享。所有服務操控同一個實體
外部副作用       | 無。內部怎麼改，外面都毫無感覺  | 有。內部修改狀態，外面同步改變
記憶體與效能開銷 | 結構體龐大時，高併發會頻繁複製  | 無論結構體多大，永遠只有 8 Byte 地址開銷
防禦 nil 危機    | 100% 安全。不可能是 nil         | 有風險。傳進來可能是 nil，需做 if 檢查
--------------------------------------------------------------------------------
後端常規選型     | 適合輕量、唯讀的 Config/Options | 後端架構 (Service/Repo/DB) 的絕對首選
--------------------------------------------------------------------------------


二、 注入選型的純文字決策樹 (Decision Tree)
--------------------------------------------------------------------------------

  [你要注入的這把武器/零件是什麼？]
                 │
                 ├─── 情況 A：它是「有狀態、要共用連線、會被修改」的基礎設施？
                 │            (例如: *sql.DB, *redis.Client, *UserService)
                 │            └─── ➜ 毫不猶豫，選 【Pointer (*struct) 注入】
                 │
                 ├─── 情況 B：它非常龐大（如超大陣列），傳值會瘋狂消耗 CPU/記憶體？
                 │            └─── ➜ 為了效能，選 【Pointer (*struct) 注入】
                 │
                 └─── 情況 C：它只是純唯讀的系統設定（Config），你希望它進去後絕對安全？
                              └─── ➜ 為了系統防禦力，選 【Struct 注入】


三、 程式碼實戰選型對比
--------------------------------------------------------------------------------

【🟢 推薦：後端元件一律用 Pointer (*struct) 注入】

type UserService struct {
    db *mongo.Database // 拿到的是連線池的「遙控器」
}

func NewUserService(db *mongo.Database) *UserService {
    // 指標注入的標準防禦性檢查
    if db == nil {
        panic("資料庫元件注入失敗，不能為 nil！")
    }
    return &UserService{db: db}
}


【🔴 特殊：不可變的防禦性場景才用 Struct 注入】

type AppConfig struct {
    Timeout int
    Debug   bool
}

type OrderService struct {
    cfg AppConfig // 拿到的是一份「複製品影本」，防範全域設定被不小心改掉
}

func NewOrderService(cfg AppConfig) *OrderService {
    return &OrderService{cfg: cfg}
}

func (s *OrderService) DoSomething() {
    s.cfg.Timeout = 999 // 哪怕不小心改了超時時間，也只是改到影本，不影響全域
}

================================================================================


四. 兩個個 Controller 型態寫了 *AppUserModel 會共用 AppUser嗎？
   需要用同一一個 object 才能共用
```go
// 如果 Controller/App/Authentication/Authenticator.SignIn, Controller/App/Resource/AppUser.ShowOnes

// 要共用一個 AppUserModel 除了必須寫出 * 指標注入外，還必須要 建構的時候指定同一個 oAppUserModel 物件


oAppUser1 := NewAppUserModel()
oAppUser2 := NewAppUserModel()

oC1 := NewControllerAppAuthenticationAuthenticator(oAppUser1)
oC2 := NewControllerAppAuthenticationAuthenticator(oAppUser2)


// 不同個 AppUser 所以沒有共用

```

五. 如果不指定會怎樣？
```go
type AppUserModel {}


type AppUserController {
    *AppUserModel
}


func NewAppUserController () {
    return *AppUserController{
        // 這邊沒有指定 AppUserModel
    }
}

/*

相當

    return *AppUserController{
        AppUserModel: nil
    }

*/

```

六. 為什麼 go 的 wire 不用 PHP 或 Nodejs 那種container 全局 map 實例的架構？
1. 在 Node.js 和 PHP 裡，容器是在**程式跑起來之後（Runtime）**才在記憶體裡動態運作的。
2. Go 語言的哲學追求的是極致的效能與絕對的穩定。Go 的造物主非常不喜歡「執行期反射」這種帶有未知風險的盲盒行為，Go 的 Wire 根本沒有「Container 容器」這個記憶體空間！
3. 


七. Di
1. DI 解決的是：「系統層級依賴（infrastructure / service / repository）」
2. 換句話說 AppUserModel(appId=3) 這種有狀態的不應該注入到 DI 裡面
3. 換言之 Go 裡面的 AppUserModel struct 都是有數據狀態的，無法注入。 
4. laravel 中的 AppUserModel 。只要注入的 appUserModel 只做 數據行為操作，不做數據保存機制， 才可以Di
5. 你以前做過那種 appId=1,2,3,4 的分表id 寫在 appUserModel 中，這個 appId 加上了 必然是有狀態的，不能直接Di
6. Di 講的是 組合依賴關係，container 解決的是全局共用一個物件，是兩件事情