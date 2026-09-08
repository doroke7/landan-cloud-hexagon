## 六角架構圖

文字版（由下往上）：

```
+----------------------------------------------------------------+
|                            input                               |
|        HTTP / gRPC / CLI / Cron / WebSocket / GraphQL          |
+----------------------------------------------------------------+
                               |
                               v
                    +----------------------+
                    |     input_port       |
                    +----------------------+
                               ^
                               |
                    +----------------------+      +----------------------+
                    |                      |----->|                      |
                    |                      |      |                      |
                    |                      |      |                      |
                    |                      |      |        domain        |
                    |       usecase        |      |                      |
                    |                      |      |                      |
                    |                      |      |                      |
                    |                      |      |                      |
                    +----------+-----------+      +----------------------+
                               |
                               v
                    +----------------------+
                    |     output_port      |
                    +----------------------+
                               ^
                               |
+----------------------------------------------------------------+
|                            output                              |
|      MySQL / Redis / Kafka / S3 / MQ / Third-party API         |
+----------------------------------------------------------------+

```
## 六角架構主要文獻
1. Clean Architecture 實作篇：在整潔的架構上弄髒你的手 (第二版) (Get Your Hands Dirty on Clean Architecture, 2/e)
2. 無瑕的程式碼－整潔的軟體設計與架構篇 (Clean Architecture: A Craftsman's Guide to Software Structure and Design)

## 六角框架使用重點

1. input 不同輸入 adaptor 相同的路徑應該是相同的設備方法, 譬如 
- (定時任務實現後台登入): cron/admin/authentication/authenticator_handler.go,   
- (CLI任務實現後台登入): command/admin/authentication/authenticator_handler.go,  
- (HTTP實現後台登入): http/admin/authentication/authenticator_handler.go,  
- (gRPC實現後台登入): facade/admin/authentication/authenticator_handler.go,  
- (WebSocket實現後台登入): websocket/admin/authentication/authenticator_handler.go,  

    以不同協議但是相同業務邏輯實現 登入取得 token 這個業務。
2. output 不同輸入 adaptor 相同的路徑應該是相同的設備方法, 譬如. 
- (Mysql 實現後台用戶簡單數據)：internal/output/application/mysql/model/admin_user_model.go,  
- (rabbitmq 實現後台用戶事件輸出)：internal/output/application/rabbitmq/event/admin_user_event.go,  
- (gRPC resource client 實現後台用戶簡單數據)：internal/output/application/resource/model/admin_user_model.go,  
      

3. usecase *一套代碼*實現所有domain業務邏輯。
4. 六角框架的也間接看出，後端的本質其實是在做*【消息傳遞】*，不論輸入輸出怎麼變

## 如何利用這套框架 從 0到100 建立一個獨立 grpc服務端-服務

1. proto基本文檔撰寫:
  建立獨立目錄的 proto/source/announcement/lottery.proto, 並生成
2. config設定當基本參數:
  設定 config/services.yaml  與 bootstrap/config.go
3. cmd 撰寫啟動入口程序：
  建立 cmd/source.go 與 container tree 與 register
4. input & handler 協議輸入代碼建立：
  建立 input/source/announcement/lottery.go 的服務類，並且綁定 pb， 並注入 container (記得注入 abstract 類)，並且註冊 register/

## 如何 建立一個獨立 grpc客戶端-連線 conn。

1. proto: 建立獨立目錄的 proto/source/announcement/lottery.proto, 並生成
2. config: 設定 config/clients.yaml  與 bootstrap/config.go
3. bootstrap： bootstrap/source.go 基本客戶端連線（conn） + internal/client/source_client.go
4. cmd ：建立 cmd/daemon.go 與 container tree 與 register/daemon.go

## 如何利用這套框架 在已經存在的 gRPC服務加功能
1. 建立 domain
2. 建立 output port
3. 請 AI 做其他實作，但是除了 resource
4. 請 AI 做 resource 實作. 


## 如何 output 從 model 改用 logic
1. 在 logic 建立 model 的相似方法，包含 resource
2. 在原本用 到 model 的 {admin,app,table}usecase 都換成 logic
3. 刪除 resource + 其他output 的 model 方法，包含 proto


## 如何拔掉 GameModel.ShowOneById
1. 先把 output 的 resource 的 GameModel.ShowOneById 改為 return nil, nil
2. 去掉 gRPC resource proto input
3. 去掉 gRPC resource 的 usecase
4. 去掉 output 的 port 去掉 GameModel ShowOneById

## 框架特點
1. 六角框架，多彈性輸入，多彈性輸出
2. 全局單例模組使用 DI 模式統一注入 一個唯一物件使用。
  （例如 MySQL 連線, Amqp 連線 這些模組的變量應該是跟 cmd/main.go 的生命週期一致，不需要一直 New物件）
3. 利用Container 模式 ，全局唯一,減少 NewW(NewX(NewY(NewZ()))) 嵌套代碼 ， 大幅增加可維護性。
4. container 模式最大好處，提前組裝 框架元件，在業務代碼上大量減少 import 行為，
   大幅度減少業務代碼 同時操作框架代碼的問題，代碼量減少，好維護。
   適用於 【框架組件 ，也就是生命週期 = main.go 的模組】。
<br>
<br>
<br>
<br>
<br>              
   再說一次，這樣的設計最大好處是，你的業務代碼裡面 再也沒有框架代碼，沒有組裝代碼，
   業務核心常常改動的地方 變得非常精簡，單純解耦。  
<br>
<br>
<br>
<br>
<br> 

5. 抽象類（AbstractXXXX） 注入，統一將需要的元件，注入到統一的抽象類裡面，以此減少重複注入代碼。  
6. 全局注入一個 context ，可以優雅關機。  
7. 使用 Aop 組件，減少緩存代碼  
8. 高度對稱性，減少人腦記憶，直覺開發，最少的定義組合出出複雜的組合 go
    input 各個實作對稱
    output 各個實作對稱
    cmd register container 啟動單元程式對稱
    結果就是 9個載體服務， 8 個邏輯微服務，高度彈性，但是代碼量卻很少。

## 我們的服務分成兩層

- **第一層是【協議層服務】**：  
  本質上只是一個服務載體，對應不同的協議實作方式（http / grpc / command / cron ...），本身不包含任何業務邏輯。  
- **第二層才是【邏輯服務】**：真正的業務邏輯放在這一層，
  例如 `admin` 服務屬於邏輯層，同一份邏輯可以同時掛載到 http、grpc、command 等不同的協議載體上，不綁定特定協議。

1. 協議服務（實例服務）：
  - Source：開獎資料來源服務載體
  - Daemon：常駐任務服務載體
  - Facade：對外主要 gRPC 服務載體
  - Resource：資料 gRPC 服務載體
  - Http：對外主要 HTTP 服務載體
  - Command：CLI 服務載體
  - Cron：排程服務載體
  - Consumer：消息隊列服務載體
  - Websocket：WebSocket 服務載體
2. 邏輯服務（虛擬服務）：
  - admin：後台介面邏輯服務，負責所有後台相關業務邏輯
  - app：前台介面邏輯服務，負責所有前台相關業務邏輯
  - third：第三方介接邏輯服務，負責所有第三方串接相關業務邏輯
  - game：前台遊戲介面邏輯服務，負責所有前台遊戲相關業務邏輯
  - table：前台資料介面邏輯服務，負責所有前台地端上報相關業務邏輯
  - register：前台驗證邏輯服務，負責所有前台身份驗證相關業務邏輯
  - model：簡單數據——單一資源的增刪改查（CRUD）
  - logic：複雜數據——跨多個資源、需要額外組合運算的資料邏輯（事務放這邊）
  - event：事件輸出——fire-and-forget，往 MQ／gRPC 丟事件，不回資料
  - announcement：開獎邏輯服務
  - watcher：採集開獎資料的邏輯服務
3. 實際運作服務堆疊

```
協議服務（載體）        實際掛載的服務（邏輯）
+-------------+     +--------------------------------+
|   Source    | --> | announcement                   |
+-------------+     +--------------------------------+
|   Daemon    | --> | watcher                        |
+-------------+     +--------------------------------+
|   Facade    | --> | game, table, register, (admin) |
+-------------+     +--------------------------------+
|   Resource  | --> | logic, model                   |
+-------------+     +--------------------------------+
|   Http      | --> | admin, app, third              |
+-------------+     +--------------------------------+
|   Command   | --> | (admin)                        |
+-------------+     +--------------------------------+
|   Cron      | --> | (admin)                        |
+-------------+     +--------------------------------+
|   Websocket | --> | (admin)                        |
+-------------+     +--------------------------------+
```

## 傳統 ThinkPHP MVC 遇到的主要問題

1. ThinkPHP 裡面的 Controller 雖然是業務主邏輯，
   但是裡面包含了 http 協議的代碼，這個會造成多 不同輸入實作有困難

## 這套框架的各個職責

1. input : 只是寫協議的對接 （如 grpc http command），決定這個服務要用在哪一個服務載體
2. usecase : 業務邏輯，基本上就是 Tp 的 C 去掉了協議的部分。
3. output/\*\*/model：簡單數據（單一資源的 CRUD）。基本上就是 Tp 的 M-model  
   output/\*\*/logic：複雜數據（跨多資源、需額外組合運算；事務放這邊！）。基本上就是 Tp 的 M-logic  
   output/\*\*/event：事件輸出（fire-and-forget，往 MQ／gRPC 丟事件，不回資料）
4. 基本上，就是這四個元件交互

## 目錄結構

```
.
├── main.go                        #   全程式進入點，只呼叫 cmd.Execute()
├── cmd/                           #   cobra 指令，每個檔案對應一個可獨立啟動的服務／進程
│   ├── root.go                    #   root command
│   ├── facade.go                  #   啟動 facade gRPC 服務（對外入口）
│   ├── resource.go                #   啟動 resource gRPC 服務（資料服務，僅供 facade / http 呼叫）
│   ├── http.go                    #   啟動 HTTP（Gin）服務（對外入口）
│   ├── consumer.go                #   啟動 AMQP consumer 服務
│   ├── source.go                  #   啟動開獎資料來源服務
│   ├── daemon.go                  #   啟動常駐任務
│   ├── cron.go                    #   啟動排程服務
│   ├── websocket.go               #   啟動 websocket 服務（對外入口）
│   └── command.go                 #   啟動 一次性 CLI 指令
│ 
│ 
│ 
├── bootstrap/                    # 讀 CONFIG、建立各種基礎設施連線（mysql / redis / amqp / mongo / grpc client） 
│ 
│ 
├── internal/
│   ├── domain/                    # 領域物件（entity）：AdminUser、AppUser、User 也就是即將操作的業務數據模型
│   ├── helper/                    # 通用工具（AES、RSA、JWT、Cache 讀寫……），跟業務邏輯無關可到處注入
│   ├── client/                    # 對外部 gRPC server 服務的 client 封裝成類別使用
│   │
│   ├── input/                     #   協議輸入端（driving adapter），只有實作，沒有介面
│   │   └── application/           #   不同協議輸入端，不同的端且相同的相對目錄 代表同一個業務輸入。
│   │       ├── facade/            #   對外 gRPC 入口：register/、table/、admin/authentication/
│   │       ├── resource/          #   resource 內部 gRPC 服務（僅供 facade / http 呼叫）：model/
│   │       ├── http/              #   對外 Http 入口：admin/authentication/、admin/resource/
│   │       ├── client/            #   gRPC client（訂閱外部 stream）：admin/resource/
│   │       ├── consumer/          #   AMQP consumer：admin/resource/
│   │       ├── source/            #   開獎資料來源服務載體：announcement/
│   │       ├── daemon/            #   常駐任務服務載體：watcher/source/
│   │       ├── cron/              #   排程任務：admin/authentication/、admin/resource/
│   │       ├── websocket/         #   websocket 入口：admin/authentication/、admin/resource/
│   │       └── command/           #   CLI 指令：admin/authentication/、admin/resource/
│   │       （每個 adapter 底下都有自己獨立的 abstract_handler.go，彼此不共用；
│   │        adapter 內部依 leaf 功能再分 admin/resource、admin/authentication 這種子資料夾——
│   │        這是跨層的對應 key，不是字面語意：input/<adapter>/admin/resource、
│   │        usecase/application/any/admin/resource、usecase/port/any/admin/resource
│   │        三者相對路徑相同，代表同一條 usecase 邏輯，不代表跟 HTTP 後台路由有關；
│   │        目前五個 adapter：http / command / cron / facade / websocket 都已經掛上同一條
│   │        admin/authentication 登入邏輯，是這套框架「一個 usecase、多種載體」的主要範例）
│   │
│   ├── middleware/                       # HTTP 專用 middleware 鏈
│   │         
│   │
│   ├── interceptor/                      # gRPC 專用攔截器鏈
│   │   
│   │   
│   │   
│   │
│   ├── usecase/                           # 商務案例：實作 + 端口介面
│   │   ├── application/                   # 實作
│   │   │   └── any/                       #  「any」表示這份 usecase 不綁定特定 input，可以被多個 input-driving 使用
│   │   │       ├── admin/                 #  後台服務
│   │   │       │   ├── authentication/    #  登入邏輯，被 http/command/cron/facade/websocket 五個 adapter 共用
│   │   │       │   └── resource/          #  CRUD 控制 列表，裡面就是增刪改查
│   │   │       │   └── option/            #  提供後台需要的 UI 數據，與 CRUD 接口 分開，避免業務耦合
│   │   │       │
│   │   │       ├── model/                 #  單一數據服務
│   │   │       ├── logic/                 #  複雜數據服務
│   │   │       ├── game/                  #  遊戲邏輯服務
│   │   │       ├── register/              #  前台驗證邏輯服務
│   │   │       ├── table/                 #  地端上報服務
│   │   │       ├── annoucement/           #  開獎邏輯服務（資料夾拼字沿用舊名，跟 port/any/announcement 拼法不同步，屬已知瑕疵）
│   │   │       └── watcher/
│   │   │           └── source/            #  採集開獎資料的邏輯服務
│   │   └── port/                          #  案例介面
│   │       └── any/                       #  對應上面每一組的 interface，同樣用「功能」命名
│   │           ├── admin/
│   │           │   ├── authentication/
│   │           │   └── resource/
│   │           ├── model/                 # model/logic 底下 struct 命名相同是正確的。 同時有 AppUserModel 處理 單一數據問題；也有 AppUserLogic 處理事務問題
│   │           ├── logic/                 # model/logic 底下 struct 命名相同是正確的。 同時有 AppUserModel 處理 單一數據問題；也有 AppUserLogic 處理事務問題
│   │           ├── game/
│   │           ├── register/
│   │           ├── table/
│   │           ├── announcement/
│   │           └── watcher/
│   │               └── source/
│   │
│   ├── output/                    # 輸出端（driven adapter）：實作 + 端口介面
│   │   ├── application/           # 輸出端-實作。每個 adapter 底下依用途分：
│   │   │   │                      #   model：簡單數據（單一資源 CRUD）
│   │   │   │                      #   logic：複雜數據（跨多資源、組合運算、事務）
│   │   │   │                      #   event：事件輸出（fire-and-forget，往 MQ／gRPC 丟事件）
│   │   │   ├── mysql/             # mysql 輸出          model/logic 底下 struct 命名相同是正確的：AppUserModel 處理單一數據、AppUserLogic 處理事務
│   │   │   │   ├── model/
│   │   │   │   └── logic/
│   │   │   ├── resource/          # resource gRPC client 輸出
│   │   │   │   ├── model/
│   │   │   │   ├── logic/
│   │   │   │   └── event/
│   │   │   ├── cache/             # redis 輸出（目前只有 model）
│   │   │   │   └── model/
│   │   │   ├── memory/            # 記憶體輸出（目前只有 model）
│   │   │   │   └── model/
│   │   │   ├── elasticsearch/     # ES 輸出（model + logic）
│   │   │   └── rabbitmq/          # rabbitmq 輸出（event）
│   │   │       └── event/
│   │   └── port/                  # 輸出端-介面
│   │       └── any/
│   │           ├── model/         # 簡單數據 port
│   │           ├── logic/         # 複雜數據 port
│   │           └── event/         # 事件輸出 port
│   │
│   └── register/                  # 組裝層：把 container 生好的 handler 註冊到對應的 server/router
│                                    #   （grpc.RegisterXxxServer / gin.Group / cron.AddFunc ...），
│                                    #   cmd/ 只管呼叫 XxxInit 拿到 server 物件再 Serve，不碰組裝細節
│
├── container/                      # 利用 wire 組裝 模組為 Di tree， 各個套件為全局單例套件
│                                   # 減少 多餘 的 不斷建構 嵌套代碼！！！！！
│       
│      
│
├── pkg/                            # 跟 domain 無關、可重用的通用元件（logger / router / cache / response / aop 等泛用工具）
│
├── config/                         # viper 讀取的 yaml 設定檔，
│
├── proto/                          # protobuf 原始定義（facade/ 對外、resource/ 資料服務、client/ 外部訂閱、source/ 資料來源）
└── pb/                             # protoc 產生的程式碼，對應 proto/ 底下的定義
```

## HTTP 服務的流程是什麼？ 
1. 在 handler 裡面，如果訊息是正常的 -> 改寫 context, context 傳遞給 ResponseMiddleware -> ResponseMiddleware 寫接口訊息
2. 在 handler 裡面，如果訊息是異常的 -> context 寫入 error -> ResponseMiddleware 有 error 不做事 -> ErrorMiddleware 寫接口訊息
3. http 服務，如果是 DefaultError 代表 可以被查看訊息的異常，resoure 服務，如果是 abort 代表業務異常，也可以被查看訊息
4. ErrorMiddleware 針對 3 的異常 選擇展示
*. 我們不使用 panic 統一捕獲錯誤後返回的策略，panic 性能差

## DI 依賴注入樹狀圖（ResourceContainer）

文字版（由下往上）：

```
┌ bootstrap ──┐
│             ├─────┐
└──────┬──────┘     │
       ▼            │
┌ pkg ────────┐     │
│             │     │
└──────┬──────┘     │      
       ▼            ▼
┌ internal/resource ───────────────────────────────────────────────────────────────────────────────────────────┐
│                                               ┌─────────┐                                                    │
│   ┌───────────────────────────────────────────┤ Helpers ├─────────────────────────────────────────────────┐  │
│   │                                           └────┬────┘                                                 │  │
│   │  ┌┄┄┄┄┄┄┄┄┄┄┄┐                                 │                                                      │  │
│   │  ┆           ┆                                 ▼                                                      │  │
│   │  ┆       ┌ internal/output ──────────────────────────────────────────────────────────┐                │  │
│   │  ┆       │ ┌─────────────────────┐  ┌─────────────────────┐  ┌─────────────────────┐ │                │  │
│   │  ┆       │ │  Mysql/Reposities   │  │  Memory/Reposities  │  │ Producer/Reposities │ │                │  │
│   │  ┆       │ └──────────┬──────────┘  └──────────┬──────────┘  └──────────┬──────────┘ │                │  │
│   │  ┆       └───────────────────────────────────────────────────────────────────────────┘                │  │
│   │  ┆           ▲                                 │                                                      │  │
│   │  └┄┄┄┄┄┄┄┄┄┄┄┘                                 ▼                                                      │  │
│   │                                          ┌───────────┐                                                │  │
│   │                                          │ Usecases  │◀───────────────────────────────────────────────┘  │
│   │                                          └─────┬─────┘                                                   │
│   │                                                ▼                                                         │
│   └───────────────────────────────────────┬────────┴──────────┐                                              │
│                                           ▼                   ▼                                              │
│                                 ┌─────────────────┐  ┌──────────────────┐                                    │
│                                 │   GrpcHandlers  │  │   Interceptors   │                                    │
│                                 └────────┬────────┘  └─────┬────────────┘                                    │
│                                          └────────┬────────┘                                                 │
│                                                   ▼                                                          │
│                                        ┌─────────────────────┐                                               │
│                                        │  ResourceContainer  │                                               │
│                                        └─────────────────────┘                                               │
└──────────────────────────────────────────────────────────────────────────────────────────────────────────────┘

```



## DI 依賴注入樹狀圖（HttpContainer）

文字版（由下往上）：

```
┌ bootstrap ──┐
│             ├─────┐
└──────┬──────┘     │
       ▼            │
┌ pkg ────────┐     │
│             │     │
└──────┬──────┘     │      
       ▼            ▼
┌ internal/http ───────────────────────────────────────────────────────────────────────────────────────────────┐
│                                               ┌─────────┐                                                    │
│   ┌───────────────────────────────────────────┤ Helpers ├─────────────────────────────────────────────────┐  │
│   │                                           └────┬────┘                                                 │  │
│   │  ┌┄┄┄┄┄┄┄┄┄┄┄┐                                 │                                                      │  │
│   │  ┆           ┆                                 ▼                                                      │  │
│   │  ┆       ┌ internal/output ──────────────────────────────────────────────────────────┐                │  │
│   │  ┆       │             ┌─────────────────────┐  ┌─────────────────────┐              │                │  │
│   │  ┆       │             │ Resource/Reposities │  │  Memory/Reposities  │              │                │  │
│   │  ┆       │             └──────────┬──────────┘  └──────────┬──────────┘              │                │  │
│   │  ┆       └───────────────────────────────────────────────────────────────────────────┘                │  │
│   │  ┆           ▲                                 │                                                      │  │
│   │  └┄┄┄┄┄┄┄┄┄┄┄┘                                 ▼                                                      │  │
│   │                                          ┌───────────┐                                                │  │
│   │                                          │  Usecases │◀───────────────────────────────────────────────┘  │
│   │                                          └─────┬─────┘                                                   │
│   │                                                ▼                                                         │
│   └───────────────────────────────────────┬────────┴──────────┐                                              │
│                                           ▼                   ▼                                              │
│                                 ┌─────────────────┐  ┌──────────────────┐                                    │
│                                 │   HttpHandlers  │  │    Middleware    │                                    │
│                                 └────────┬────────┘  └─────┬────────────┘                                    │
│                                          └────────┬────────┘                                                 │
│                                                   ▼                                                          │
│                                        ┌─────────────────────┐                                               │
│                                        │    HttpContainer    │                                               │
│                                        └─────────────────────┘                                               │
└──────────────────────────────────────────────────────────────────────────────────────────────────────────────┘


```

## 服務拓樸

```

                                         ┌───────────┐
                                         │  Source   │                                 ┌──────────────┐
                                         └───────────┘                                 │              │
                                               ▲                                       ▲              │
                                               │ gRPC 呼叫                              │              │
                                               │ (訂閱 Watch stream)                    │              │
                                               │                                       │              │
┌──────────┐  ┌─────────┐  ┌───────────┐  ┌──────────┐  ┌─────────┐  ┌────────┐  ┌───────────┐        │
│  Facade  │  │   Http  │  │ Websocket │  │  Daemon  │  │ Command │  │  Cron  │  │ Consumer  │        │
└──────────┘  └─────────┘  └───────────┘  └──────────┘  └─────────┘  └────────┘  └───────────┘        │
     │             │             │             │             │           │             │              │
     └─────────────┼─────────────┘             │             │           │             │              │
                   │                           │             │           │             │              │
                   │ gRPC 呼叫                  │             │           │             │              │
                   ▼                           │             │           │             │              │
        ┌────────────────────┐                 │             │           │             │              │
        │      Resource      │                 │             │           │             │              │
        └────────────────────┘                 │             │           │             │              │
                   │                           │             │           │             │              │
             ┌─────┴─────┐                     │             │           │             │              │
             │           │                     │             │           │             │              │
             ▼           ▼                     ▼             ▼           ▼             ▼              │
 ┌─────────────────────────────────────────────────────────────────────────────────────────┐          │
 │                ┌───────┐              ┌───────┐                   ┌──────────┐          │          │
 │                │ Redis │              │ MySQL │                   │ RabbitMQ │          │          │
 │                └───────┘              └───────┘                   └────▲─────┘          │          │
 └────────────────────────────────────────────────────────────────────────+────────────────┘          │
                                                                          │                           │
                                                                          └───────────────────────────┘
                                                                                                      
```

## 如何 watch 开发 (熱更新開發)

1. go mod 安装下载 air 套件

```zsh
go install github.com/air-verse/air@latest
```


2. 以 air 啟動 go 服務 （以開發 http 服務為例）

```zsh
air -c .air.resource.toml # 啟動 gRPC resource 資料庫服務
air -c .air.http.toml     # 啟動 http 服務
```

## 微服務地端如果有強大的數據需求，又對網路延遲鳴感，怎麼辦？

1. 在地端 使用 SQL lite 思維，本地 SQL思維
2. 再統一上報數據

## HTTP 使用標準 RESTful 四劍客（傳值輸入, 過濾查詢, 排序規則, 分頁機制）
1. value (傳值輸入）
   post body 輸入的資料，用於新增修改
2. filters (過濾塞選)
   get param 查詢過濾器列表
3. sorters (排序規則）
   get param 排序器列表
4. piganition (分頁機制）
   get param 分頁數據


## 排查步驟

1. **確認 API 用的是哪一個 usecase**：從 `input/application/<載體>/**/xxx_handler.go` 找到它注入的 usecase 介面（`usecase/port/any/...`）。
2. **看這個 usecase 呼叫了哪些數據 port**：進 `usecase/application/any/...` 的實作，確認它用的是 `outputPortAnyModel.XxxModel` 還是 `outputPortAnyLogic.XxxLogic`（或兩者）。
3. **數據輸出是 gRPC 時，檢查 client 呼叫方式**：`output/application/resource/{model,logic}/xxx.go` 裡打的是哪一支 rpc、`ResourceLogicClient` 還是 `ResourceModelClient`、request 欄位有沒有對齊 proto。
4. **確認 resource 端用的是哪一個 usecase**：resource 服務走 `usecase/application/any/{model,logic}/xxx_usecase.go`，對照它掛在 `input/application/resource/{model,logic}/xxx_handler.go`。
5. **層次規則**：resource 的 `logic` usecase 只能呼叫 output 的 `logic`；`model` usecase 只能呼叫 output 的 `model`。跨過去就是接錯層。
6. **最後看 output model / logic 的實作**：確認真正打 DB 的那段（gorm / bson / es query）、soft-delete 條件、事務範圍、Aop 快取 key 是否一致。


```
使用案例 usecase              數據 port（呼叫 model / logic）      output 實作
─────────────────────────    ──────────────────────────────    ────────────────────────────────
ResourceAdminUserUsecase     outputAdminUserLogic  ─────────▶   mysql
  .RemoveOne                 （介面）                            mongodb
                                                                elasticsearch
                                                                resource ──┬─▶ logicClient   ✔ 用這個
                                                                           └─▶ modelClient   ✗ 不是這個
```

**resource 載體**（資料 gRPC 服務，被 http / facade 呼叫）：

```
使用案例 usecase                    數據 port（resource logic 只能呼叫 logic）   output 實作
─────────────────────────────      ─────────────────────────────────────────   ─────────────
usecase/application/any/logic/     outputAdminUserLogic  ──────────────────▶   mysql
  admin_user_usecase.go            （介面）                                     mongodb

                                   outputAdminUserModel   ✗ resource 的 logic 不會走 model
```

## 刪掉數據跟 關係架構的注意事項
A. 無關係 。（configs 表），可以直接刪除 config 一筆
B. 屬於關係(Person 屬於 Nation), nation_id 在 person 裡面。可以直接刪除 person 一筆
C. 擁有關係(AppUser 擁有 AppUserOrder), app_user_id 在 AppUserOrder 裡面。刪除AppUser 一筆之前，需要把旗下的 AppUserOrder 都刪除
   如：Menu 跟 Menu，GameType 跟 GameType，
D. 多對多關係(AdminUser <-> AdminRole)。刪除 AdminUser 一筆之前，先把 admin_users_to_admin_roles 旗 admin_user 的關係刪掉，資料不動。


## 關於變量縮寫的問題
1. 傾向於不縮寫，因為不同的人，不同的時間縮寫習慣不同。會造成代碼翻譯的問題。
2. 1990年代以前 IDE 不成熟，沒有自動 lint, 電腦屏幕窄，沒有git，只有 vim，語言多為低階語言，記憶體有限，使用全拼寫撰寫的成本較高
3. 現在IDE 的ui 成熟，自動補全機制，電腦可寬，語言高階，有 git， 記憶體極巨大冗余，有AI。使用短命名的成本較高
4. 這麼喜歡短命名怎麼不用 abcd 命名
5. 有共識的資料可以縮寫，譬如: API JSON DOM