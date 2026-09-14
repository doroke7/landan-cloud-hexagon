package client

import (
	"google.golang.org/grpc"

	pbResourceEvent "example/pb/resource/event"
	pbResourceLogic "example/pb/resource/logic"
	pbResourceModel "example/pb/resource/model"
)

func NewModel(oClientConn *grpc.ClientConn) *Model {
	oModel := &Model{
		AdminUser:            pbResourceModel.NewAdminUserModelClient(oClientConn),
		AppUser:              pbResourceModel.NewAppUserModelClient(oClientConn),
		Game:                 pbResourceModel.NewGameModelClient(oClientConn),
		Table:                pbResourceModel.NewTableModelClient(oClientConn),
		GameType:             pbResourceModel.NewGameTypeModelClient(oClientConn),
		AdminRole:            pbResourceModel.NewAdminRoleModelClient(oClientConn),
		AdminPermission:      pbResourceModel.NewAdminPermissionModelClient(oClientConn),
		AdminPermissionGroup: pbResourceModel.NewAdminPermissionGroupModelClient(oClientConn),
	}

	return oModel
}

func NewLogic(oClientConn *grpc.ClientConn) *Logic {
	oLogic := &Logic{
		AdminUser:            pbResourceLogic.NewAdminUserLogicClient(oClientConn),
		AdminPermissionGroup: pbResourceLogic.NewAdminPermissionGroupLogicClient(oClientConn),
		Game:                 pbResourceLogic.NewGameLogicClient(oClientConn),
		Table:                pbResourceLogic.NewTableLogicClient(oClientConn),
		GameType:             pbResourceLogic.NewGameTypeLogicClient(oClientConn),
	}

	return oLogic
}

func NewEvent(oClientConn *grpc.ClientConn) *Event {
	oEvent := &Event{
		AdminUser: pbResourceEvent.NewAdminUserEventClient(oClientConn),
	}

	return oEvent
}

type Model struct {
	AdminUser            pbResourceModel.AdminUserModelClient
	AdminRole            pbResourceModel.AdminRoleModelClient
	AdminPermission      pbResourceModel.AdminPermissionModelClient
	AdminPermissionGroup pbResourceModel.AdminPermissionGroupModelClient
	AppUser              pbResourceModel.AppUserModelClient
	Game                 pbResourceModel.GameModelClient
	Table                pbResourceModel.TableModelClient
	GameType             pbResourceModel.GameTypeModelClient
}

type Logic struct {
	AdminUser            pbResourceLogic.AdminUserLogicClient
	AdminPermissionGroup pbResourceLogic.AdminPermissionGroupLogicClient
	Game                 pbResourceLogic.GameLogicClient
	Table                pbResourceLogic.TableLogicClient
	GameType             pbResourceLogic.GameTypeLogicClient
}

type Event struct {
	AdminUser pbResourceEvent.AdminUserEventClient
}

/*
    為何這裡寫 【AdminUser pbResourceModel.AdminUserModelClient】 而不是 【AdminUser *pbResourceModel.AdminUserModelClient】
	是可行的？
	1. 首先寫 * 的用意是為了 全系統 連線變量唯一
	2. 何時可能不唯一 -> 存在兩個地方使用， 譬如 Controller 跟 Middleware 都注入了 Helper類
	3. 但是這邊 AdminUserModelClient 只有 Resource.Model 使用


	4. AdminUserModelClient 是 Interface， 我們不需要對 interface 寫 *
*/

////////////////////////////////////////////////////////////////////////////

func NewResourceClient(oClientConn *grpc.ClientConn, oModel *Model, oLogic *Logic, oEvent *Event) *ResourceClient {

	return &ResourceClient{
		conn:  oClientConn,
		Model: oModel,
		Logic: oLogic,
		Event: oEvent,
	}
}

type ResourceClient struct {
	conn  *grpc.ClientConn
	Model *Model // 這樣不厭其煩的命名 【嵌套結構】，是為了與 server 【命名空間一致性】，增加可讀性。
	Logic *Logic
	Event *Event
}

func (oClient *ResourceClient) Close() error {
	oErr := oClient.conn.Close()

	return oErr
}
