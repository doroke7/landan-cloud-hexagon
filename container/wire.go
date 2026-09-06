//go:build wireinject
// +build wireinject

package container

/*

  我好像應該把一個全局 context 注入 到 各個元件中，如果程序終止的話， 各個程序也停止

*/

import (
	"context"

	"github.com/google/wire"
	"github.com/nats-io/nats.go"

	bootstrap "example/bootstrap"
	pkgCache "example/pkg/cache"
	pkgGin "example/pkg/gin"
	pkgNats "example/pkg/nats"
	pkgUtility "example/pkg/utility"

	client "example/internal/client"

	helper "example/internal/helper"

	outputApplicationCache "example/internal/output/application/cache"
	outputApplicationCacheModel "example/internal/output/application/cache/model"

	outputApplicationMemory "example/internal/output/application/memory"
	outputApplicationMemoryModel "example/internal/output/application/memory/model"

	outputApplicationMysql "example/internal/output/application/mysql"
	outputApplicationMysqlLogic "example/internal/output/application/mysql/logic"
	outputApplicationMysqlModel "example/internal/output/application/mysql/model"

	outputApplicationRabbitmq "example/internal/output/application/rabbitmq"
	outputApplicationRabbitmqEvent "example/internal/output/application/rabbitmq/event"
	outputApplicationResource "example/internal/output/application/resource"
	outputApplicationResourceLogic "example/internal/output/application/resource/logic"
	outputApplicationResourceModel "example/internal/output/application/resource/model"
	outputPortAnyEvent "example/internal/output/port/any/event"

	usecasePortAnyModel "example/internal/usecase/port/any/model"

	usecaseApplicationAnyAdmin "example/internal/usecase/application/any/admin"
	usecaseApplicationAnyAdminAuthentication "example/internal/usecase/application/any/admin/authentication"
	usecaseApplicationAnyAdminOption "example/internal/usecase/application/any/admin/option"
	usecaseApplicationAnyAdminResource "example/internal/usecase/application/any/admin/resource"
	usecaseApplicationAnyAnnouncement "example/internal/usecase/application/any/annoucement"
	usecaseApplicationAnyEvent "example/internal/usecase/application/any/event"
	usecaseApplicationAnyGame "example/internal/usecase/application/any/game"
	usecaseApplicationAnyGameAuthentication "example/internal/usecase/application/any/game/authentication"
	usecaseApplicationAnyLogic "example/internal/usecase/application/any/logic"
	usecaseApplicationAnyModel "example/internal/usecase/application/any/model"
	usecaseApplicationAnyWatcher "example/internal/usecase/application/any/watcher"
	usecaseApplicationAnyWatcherSource "example/internal/usecase/application/any/watcher/source"

	middlewareHttpAdmin "example/internal/middleware/http/admin"
	middlewareHttpGame "example/internal/middleware/http/game"
	middlewareHttpTable "example/internal/middleware/http/table"

	interceptorFacadeAdmin "example/internal/interceptor/facade/admin"
	interceptorFacadeGame "example/internal/interceptor/facade/game"
	interceptorResourceLogic "example/internal/interceptor/resource/logic"
	interceptorResourceModel "example/internal/interceptor/resource/model"

	inputApplicationCommand "example/internal/input/application/command"
	inputApplicationCommandAdminAuthentication "example/internal/input/application/command/admin/authentication"
	inputApplicationCommandAdminResource "example/internal/input/application/command/admin/resource"

	inputApplicationRabbitmq "example/internal/input/application/rabbitmq"
	inputApplicationRabbitmqAdminResource "example/internal/input/application/rabbitmq/admin/resource"

	inputApplicationCron "example/internal/input/application/cron"
	inputApplicationCronAdminAuthentication "example/internal/input/application/cron/admin/authentication"
	inputApplicationCronAdminResource "example/internal/input/application/cron/admin/resource"

	inputApplicationSource "example/internal/input/application/source"
	inputApplicationSourceAnnouncement "example/internal/input/application/source/announcement"

	inputApplicationDaemon "example/internal/input/application/daemon"
	inputApplicationDaemonWatcherSource "example/internal/input/application/daemon/watcher/source"

	inputApplicationResource "example/internal/input/application/resource"
	inputApplicationResourceEvent "example/internal/input/application/resource/event"
	inputApplicationResourceLogic "example/internal/input/application/resource/logic"
	inputApplicationResourceModel "example/internal/input/application/resource/model"

	inputApplicationFacade "example/internal/input/application/facade"
	inputApplicationFacadeAdminAuthentication "example/internal/input/application/facade/admin/authentication"
	inputApplicationFacadeAdminResource "example/internal/input/application/facade/admin/resource"
	inputApplicationFacadeRegister "example/internal/input/application/facade/register"
	inputApplicationFacadeTable "example/internal/input/application/facade/table"

	inputApplicationHttp "example/internal/input/application/http"
	inputApplicationHttpAdminAuthentication "example/internal/input/application/http/admin/authentication"
	inputApplicationHttpAdminOption "example/internal/input/application/http/admin/option"
	inputApplicationHttpAdminResource "example/internal/input/application/http/admin/resource"
	inputApplicationHttpGameAuthentication "example/internal/input/application/http/game/authentication"

	inputApplicationTcp "example/internal/input/application/tcp"
	inputApplicationTcpAdminAuthentication "example/internal/input/application/tcp/admin/authentication"

	inputApplicationWebsocket "example/internal/input/application/websocket"
	inputApplicationWebsocketAdminAuthentication "example/internal/input/application/websocket/admin/authentication"

	middlewareWebsocketAdmin "example/internal/middleware/websocket/admin"
)

// HttpContainer 只給 `http` Gin 服務使用。
type HttpContainer struct {

	// pkg
	*pkgGin.Response
	Clock *pkgUtility.Clock

	// Helper
	*helper.AbstractHelper
	*helper.AesHelper
	*helper.RsaHelper
	*helper.JwtHelper

	// Clients
	ResourceClient *client.ResourceClient

	// HTTP server -Controller
	HttpAdminAuthenticationAuthenticator *inputApplicationHttpAdminAuthentication.AuthenticatorHandler
	HttpAdminResourceGame                *inputApplicationHttpAdminResource.GameHandler
	HttpAdminResourceTable               *inputApplicationHttpAdminResource.TableHandler
	HttpAdminResourceAdminUser           *inputApplicationHttpAdminResource.AdminUserHandler
	HttpAdminResourceAdminRole           *inputApplicationHttpAdminResource.AdminRoleHandler
	HttpAdminResourceAdminPermission     *inputApplicationHttpAdminResource.AdminPermissionHandler
	HttpAdminResourceGameType            *inputApplicationHttpAdminResource.GameTypeHandler
	HttpAdminOptionGameType              *inputApplicationHttpAdminOption.GameTypeHandler
	HttpAdminOptionAdminRole             *inputApplicationHttpAdminOption.AdminRoleHandler
	HttpAdminOptionAdminPermission       *inputApplicationHttpAdminOption.AdminPermissionHandler
	HttpGameAuthenticationAuthenticator  *inputApplicationHttpGameAuthentication.AuthenticatorHandler

	// HTTP server -Middleware
	// Middleware 部分
	HttpAdminAbstractMiddleware       *middlewareHttpAdmin.AbstractMiddleware
	HttpAdminAdminMiddleware          *middlewareHttpAdmin.AdminMiddleware
	HttpAdminAuthenticationMiddleware *middlewareHttpAdmin.AuthenticationMiddleware
	HttpAdminDecryptionMiddleware     *middlewareHttpAdmin.DecryptionMiddleware
	HttpAdminEncryptionMiddleware     *middlewareHttpAdmin.EncryptionMiddleware
	HttpAdminErrorMiddleware          *middlewareHttpAdmin.ErrorMiddleware
	HttpAdminLoggerMiddleware         *middlewareHttpAdmin.LoggerMiddleware
	HttpAdminNonexistentMiddleware    *middlewareHttpAdmin.NonexistentMiddleware
	HttpAdminRequestMiddleware        *middlewareHttpAdmin.RequestMiddleware
	HttpAdminResponseMiddleware       *middlewareHttpAdmin.ResponseMiddleware
	HttpAdminSignatureMiddleware      *middlewareHttpAdmin.SignatureMiddleware

	HttpTableAbstractMiddleware       *middlewareHttpTable.AbstractMiddleware
	HttpTableTableMiddleware          *middlewareHttpTable.TableMiddleware
	HttpTableAuthenticationMiddleware *middlewareHttpTable.AuthenticationMiddleware
	HttpTableDecryptionMiddleware     *middlewareHttpTable.DecryptionMiddleware
	HttpTableEncryptionMiddleware     *middlewareHttpTable.EncryptionMiddleware
	HttpTableErrorMiddleware          *middlewareHttpTable.ErrorMiddleware
	HttpTableLoggerMiddleware         *middlewareHttpTable.LoggerMiddleware
	HttpTableNonexistentMiddleware    *middlewareHttpTable.NonexistentMiddleware
	HttpTableRequestMiddleware        *middlewareHttpTable.RequestMiddleware
	HttpTableResponseMiddleware       *middlewareHttpTable.ResponseMiddleware
	HttpTableSignatureMiddleware      *middlewareHttpTable.SignatureMiddleware

	HttpGameAbstractMiddleware       *middlewareHttpGame.AbstractMiddleware
	HttpGameGameMiddleware           *middlewareHttpGame.GameMiddleware
	HttpGameAuthenticationMiddleware *middlewareHttpGame.AuthenticationMiddleware
	HttpGameDecryptionMiddleware     *middlewareHttpGame.DecryptionMiddleware
	HttpGameEncryptionMiddleware     *middlewareHttpGame.EncryptionMiddleware
	HttpGameErrorMiddleware          *middlewareHttpGame.ErrorMiddleware
	HttpGameLoggerMiddleware         *middlewareHttpGame.LoggerMiddleware
	HttpGameNonexistentMiddleware    *middlewareHttpGame.NonexistentMiddleware
	HttpGameRequestMiddleware        *middlewareHttpGame.RequestMiddleware
	HttpGameResponseMiddleware       *middlewareHttpGame.ResponseMiddleware
	HttpGameSignatureMiddleware      *middlewareHttpGame.SignatureMiddleware
}

func InitHttpContainer(ctx context.Context) (*HttpContainer, error) {
	wire.Build(

		// pkg
		pkgGin.NewResponse,
		pkgUtility.NewClock,

		// bootstrap
		bootstrap.NewResource,

		// helper
		helper.NewAbstractHelper,
		helper.NewAesHelper,
		helper.NewRsaHelper,
		helper.NewJwtHelper,
		helper.NewValiatorHelper,

		// output
		outputApplicationResource.NewAbstractResource,
		outputApplicationResourceLogic.NewAbstractLogic,
		outputApplicationResourceModel.NewAbstractModel,
		outputApplicationResourceModel.NewAdminUserModel,
		outputApplicationResourceModel.NewAdminRoleModel,
		outputApplicationResourceModel.NewAdminPermissionModel,
		outputApplicationResourceModel.NewGameModel,
		outputApplicationResourceModel.NewTableModel,
		outputApplicationResourceModel.NewGameTypeModel,
		outputApplicationResourceLogic.NewAdminUserLogic,
		outputApplicationResourceLogic.NewGameLogic,
		outputApplicationResourceLogic.NewTableLogic,
		outputApplicationResourceLogic.NewGameTypeLogic,
		outputApplicationResourceModel.NewAppUserModel,

		// usecase
		usecaseApplicationAnyAdmin.NewAbstractUsecase,
		usecaseApplicationAnyAdminAuthentication.NewAuthenticatorUsecase,
		usecaseApplicationAnyGame.NewAbstractUsecase,
		usecaseApplicationAnyGameAuthentication.NewAuthenticatorUsecase,
		usecaseApplicationAnyAdminResource.NewGameUsecase,
		usecaseApplicationAnyAdminResource.NewTableUsecase,
		usecaseApplicationAnyAdminResource.NewAdminUserUsecase,
		usecaseApplicationAnyAdminResource.NewAdminRoleUsecase,
		usecaseApplicationAnyAdminResource.NewAdminPermissionUsecase,
		usecaseApplicationAnyAdminResource.NewGameTypeUsecase,
		usecaseApplicationAnyAdminOption.NewGameTypeUsecase,
		usecaseApplicationAnyAdminOption.NewAdminRoleUsecase,
		usecaseApplicationAnyAdminOption.NewAdminPermissionUsecase,

		// client
		client.NewModel,
		client.NewEvent,
		client.NewLogic,
		client.NewResourceClient,

		// input-http
		inputApplicationHttp.NewAbstractHandler,

		inputApplicationHttpAdminAuthentication.NewAuthenticatorHandler,
		inputApplicationHttpGameAuthentication.NewAuthenticatorHandler,
		inputApplicationHttpAdminResource.NewGameHandler,
		inputApplicationHttpAdminResource.NewTableHandler,
		inputApplicationHttpAdminResource.NewAdminUserHandler,
		inputApplicationHttpAdminResource.NewAdminRoleHandler,
		inputApplicationHttpAdminResource.NewAdminPermissionHandler,
		inputApplicationHttpAdminResource.NewGameTypeHandler,
		inputApplicationHttpAdminOption.NewGameTypeHandler,
		inputApplicationHttpAdminOption.NewAdminRoleHandler,
		inputApplicationHttpAdminOption.NewAdminPermissionHandler,

		// Middleware 部分
		middlewareHttpAdmin.NewAbstractMiddleware,
		middlewareHttpAdmin.NewAdminMiddleware,
		middlewareHttpAdmin.NewAuthenticationMiddleware,
		middlewareHttpAdmin.NewDecryptionMiddleware,
		middlewareHttpAdmin.NewEncryptionMiddleware,
		middlewareHttpAdmin.NewErrorMiddleware,
		middlewareHttpAdmin.NewLoggerMiddleware,
		middlewareHttpAdmin.NewNonexistentMiddleware,
		middlewareHttpAdmin.NewRequestMiddleware,
		middlewareHttpAdmin.NewResponseMiddleware,
		middlewareHttpAdmin.NewSignatureMiddleware,

		middlewareHttpTable.NewAbstractMiddleware,
		middlewareHttpTable.NewTableMiddleware,
		middlewareHttpTable.NewAuthenticationMiddleware,
		middlewareHttpTable.NewDecryptionMiddleware,
		middlewareHttpTable.NewEncryptionMiddleware,
		middlewareHttpTable.NewErrorMiddleware,
		middlewareHttpTable.NewLoggerMiddleware,
		middlewareHttpTable.NewNonexistentMiddleware,
		middlewareHttpTable.NewRequestMiddleware,
		middlewareHttpTable.NewResponseMiddleware,
		middlewareHttpTable.NewSignatureMiddleware,

		middlewareHttpGame.NewAbstractMiddleware,
		middlewareHttpGame.NewGameMiddleware,
		middlewareHttpGame.NewAuthenticationMiddleware,
		middlewareHttpGame.NewDecryptionMiddleware,
		middlewareHttpGame.NewEncryptionMiddleware,
		middlewareHttpGame.NewErrorMiddleware,
		middlewareHttpGame.NewLoggerMiddleware,
		middlewareHttpGame.NewNonexistentMiddleware,
		middlewareHttpGame.NewRequestMiddleware,
		middlewareHttpGame.NewResponseMiddleware,
		middlewareHttpGame.NewSignatureMiddleware,

		wire.Struct(new(HttpContainer), "*"),
	)
	return nil, nil
}

type FacadeContainer struct {

	// pkg
	Clock *pkgUtility.Clock

	// Helper
	*helper.AbstractHelper
	*helper.AesHelper
	*helper.RsaHelper
	*helper.JwtHelper

	// Clients
	ResourceClient *client.ResourceClient

	// gRPC Facade service
	FacadeAbstract                         *inputApplicationFacade.AbstractHandler
	FacadeTableScanner                     *inputApplicationFacadeTable.ScannerHandler
	FacadeTableAuthenticator               *inputApplicationFacadeRegister.AuthenticatorHandler
	FacadeAdminAuthenticationAuthenticator *inputApplicationFacadeAdminAuthentication.AuthenticatorHandler
	FacadeAdminResourceGame                *inputApplicationFacadeAdminResource.GameHandler

	// gRPC Facade Interceptor
	FacadeGameErrorInterceptor          *interceptorFacadeGame.ErrorInterceptor
	FacadeGameStatusInterceptor         *interceptorFacadeGame.StatusInterceptor
	FacadeGameLoggerInterceptor         *interceptorFacadeGame.LoggerInterceptor
	FacadeGameAuthenticationInterceptor *interceptorFacadeGame.AuthenticationInterceptor
	FacadeAdminErrorInterceptor         *interceptorFacadeAdmin.ErrorInterceptor
	FacadeAdminStatusInterceptor        *interceptorFacadeAdmin.StatusInterceptor
	FacadeAdminLoggerInterceptor        *interceptorFacadeAdmin.LoggerInterceptor
	FacadeAdminSignatureInterceptor     *interceptorFacadeAdmin.SignatureInterceptor
	FacadeAdminDecryptionInterceptor    *interceptorFacadeAdmin.DecryptionInterceptor
	FacadeAdminEncryptionInterceptor    *interceptorFacadeAdmin.EncryptionInterceptor
}

func InitFacadeContainer(ctx context.Context) (*FacadeContainer, error) {
	wire.Build(

		// pkg
		pkgUtility.NewClock,

		// bootstrap
		bootstrap.NewResource,

		// helper
		helper.NewAbstractHelper,
		helper.NewAesHelper,
		helper.NewRsaHelper,
		helper.NewJwtHelper,
		helper.NewValiatorHelper,

		// output
		outputApplicationResource.NewAbstractResource,
		outputApplicationResourceLogic.NewAbstractLogic,
		outputApplicationResourceModel.NewAbstractModel,
		outputApplicationResourceModel.NewAdminUserModel,
		outputApplicationResourceModel.NewGameModel,
		outputApplicationResourceLogic.NewGameLogic,

		// usecase
		usecaseApplicationAnyAdmin.NewAbstractUsecase,
		usecaseApplicationAnyAdminAuthentication.NewAuthenticatorUsecase,
		usecaseApplicationAnyAdminResource.NewGameUsecase,

		// client
		client.NewModel,
		client.NewEvent,
		client.NewLogic,
		client.NewResourceClient,

		// input-facade
		inputApplicationFacade.NewAbstractHandler,
		inputApplicationFacadeTable.NewScannerHandler,
		inputApplicationFacadeRegister.NewAuthenticatorHandler,
		inputApplicationFacadeAdminAuthentication.NewAuthenticatorHandler,
		inputApplicationFacadeAdminResource.NewGameHandler,

		// interceptor-facade
		interceptorFacadeGame.NewAbstractInterceptor,
		interceptorFacadeGame.NewErrorInterceptor,
		interceptorFacadeGame.NewStatusInterceptor,
		interceptorFacadeGame.NewLoggerInterceptor,
		interceptorFacadeGame.NewAuthenticationInterceptor,
		interceptorFacadeAdmin.NewAbstractInterceptor,
		interceptorFacadeAdmin.NewErrorInterceptor,
		interceptorFacadeAdmin.NewStatusInterceptor,
		interceptorFacadeAdmin.NewLoggerInterceptor,
		interceptorFacadeAdmin.NewSignatureInterceptor,
		interceptorFacadeAdmin.NewDecryptionInterceptor,
		interceptorFacadeAdmin.NewEncryptionInterceptor,

		wire.Struct(new(FacadeContainer), "*"),
	)
	return nil, nil
}

//////////////////////////////////////////////////////////////////////////////

type ResourceContainer struct {

	// pkg
	Clock *pkgUtility.Clock

	// Helper
	*helper.AbstractHelper
	*helper.AesHelper
	*helper.RsaHelper

	*usecaseApplicationAnyModel.AbstractUsecase
	usecasePortAnyModel.AdminUserUsecase
	usecasePortAnyModel.AppUserUsecase
	usecasePortAnyModel.GameUsecase
	usecasePortAnyModel.TableUsecase
	usecasePortAnyModel.GameTypeUsecase

	// gRPC Resource server
	ResourceAbstract                  *inputApplicationResource.AbstractHandler
	ResourceModelAdminUser            *inputApplicationResourceModel.AdminUserHandler
	ResourceModelAdminRole            *inputApplicationResourceModel.AdminRoleHandler
	ResourceModelAdminPermission      *inputApplicationResourceModel.AdminPermissionHandler
	ResourceModelAppUser              *inputApplicationResourceModel.AppUserHandler
	ResourceModelGame                 *inputApplicationResourceModel.GameHandler
	ResourceModelTable                *inputApplicationResourceModel.TableHandler
	ResourceModelGameType             *inputApplicationResourceModel.GameTypeHandler
	ResourceLogicAdminUser            *inputApplicationResourceLogic.AdminUserHandler
	ResourceLogicGame                 *inputApplicationResourceLogic.GameHandler
	ResourceLogicTable                *inputApplicationResourceLogic.TableHandler
	ResourceLogicGameType             *inputApplicationResourceLogic.GameTypeHandler
	ResourceLogicAdminPermissionGroup *inputApplicationResourceLogic.AdminPermissionGroupHandler
	ResourceEventAdminUser            *inputApplicationResourceEvent.AdminUserHandler

	// gRPC Resource Interceptor
	ResourceLogicAuthenticationInterceptor *interceptorResourceLogic.AuthenticationInterceptor
	ResourceLogicErrorInterceptor          *interceptorResourceLogic.ErrorInterceptor
	ResourceLogicLoggerInterceptor         *interceptorResourceLogic.LoggerInterceptor
	ResourceModelAuthenticationInterceptor *interceptorResourceModel.AuthenticationInterceptor
	ResourceModelErrorInterceptor          *interceptorResourceModel.ErrorInterceptor
	ResourceModelLoggerInterceptor         *interceptorResourceModel.LoggerInterceptor

	// MQ 生產者
	ResourceRabbitmqAdminUserEvent outputPortAnyEvent.AdminUserEvent
}

func InitResourceContainer(ctx context.Context) (*ResourceContainer, error) {
	wire.Build(

		// bootstrap
		bootstrap.NewMysql,
		bootstrap.NewAmqp,
		bootstrap.NewRedis,
		pkgCache.NewAop,
		pkgUtility.NewClock,

		// helper
		helper.NewAbstractHelper,
		helper.NewAesHelper,
		helper.NewRsaHelper,

		// output
		outputApplicationMysql.NewAbstractMysql,
		outputApplicationMysqlModel.NewAbstractModel,
		outputApplicationMysqlModel.NewAdminUserModel,
		outputApplicationMysqlModel.NewAdminRoleModel,
		outputApplicationMysqlModel.NewAdminPermissionModel,
		outputApplicationMysqlModel.NewAppUserModel,
		outputApplicationMysqlModel.NewGameModel,
		outputApplicationMysqlModel.NewTableModel,
		outputApplicationMysqlModel.NewGameTypeModel,
		outputApplicationMysqlLogic.NewAbstractLogic,
		outputApplicationMysqlLogic.NewAdminUserLogic,
		outputApplicationMysqlLogic.NewGameLogic,
		outputApplicationMysqlLogic.NewTableLogic,
		outputApplicationMysqlLogic.NewGameTypeLogic,
		outputApplicationMysqlLogic.NewAdminPermissionGroupLogic,
		outputApplicationRabbitmq.NewAbstractRabbitmq,
		outputApplicationRabbitmqEvent.NewAbstractEvent,
		outputApplicationRabbitmqEvent.NewAdminUserEvent,

		// usecase
		usecaseApplicationAnyModel.NewAbstractUsecase,
		usecaseApplicationAnyModel.NewAdminUserUsecase,
		usecaseApplicationAnyModel.NewAdminRoleUsecase,
		usecaseApplicationAnyModel.NewAdminPermissionUsecase,
		usecaseApplicationAnyModel.NewAppUserUsecase,
		usecaseApplicationAnyModel.NewGameUsecase,
		usecaseApplicationAnyModel.NewTableUsecase,
		usecaseApplicationAnyModel.NewGameTypeUsecase,
		usecaseApplicationAnyLogic.NewAbstractUsecase,
		usecaseApplicationAnyLogic.NewAdminUserUsecase,
		usecaseApplicationAnyLogic.NewGameUsecase,
		usecaseApplicationAnyLogic.NewTableUsecase,
		usecaseApplicationAnyLogic.NewGameTypeUsecase,
		usecaseApplicationAnyLogic.NewAdminPermissionGroupUsecase,
		usecaseApplicationAnyEvent.NewAdminUserUsecase,

		// input-resource
		inputApplicationResource.NewAbstractHandler,
		inputApplicationResourceModel.NewAdminUserHandler,
		inputApplicationResourceModel.NewAdminRoleHandler,
		inputApplicationResourceModel.NewAdminPermissionHandler,
		inputApplicationResourceModel.NewAppUserHandler,
		inputApplicationResourceModel.NewGameHandler,
		inputApplicationResourceModel.NewTableHandler,
		inputApplicationResourceModel.NewGameTypeHandler,
		inputApplicationResourceLogic.NewAdminUserHandler,
		inputApplicationResourceLogic.NewGameHandler,
		inputApplicationResourceLogic.NewTableHandler,
		inputApplicationResourceLogic.NewGameTypeHandler,
		inputApplicationResourceLogic.NewAdminPermissionGroupHandler,
		inputApplicationResourceEvent.NewAdminUserHandler,

		// interceptor-resource
		interceptorResourceLogic.NewAbstractInterceptor,
		interceptorResourceLogic.NewAuthenticationInterceptor,
		interceptorResourceLogic.NewErrorInterceptor,
		interceptorResourceLogic.NewLoggerInterceptor,
		interceptorResourceModel.NewAbstractInterceptor,
		interceptorResourceModel.NewAuthenticationInterceptor,
		interceptorResourceModel.NewErrorInterceptor,
		interceptorResourceModel.NewLoggerInterceptor,

		wire.Struct(new(ResourceContainer), "*"),
	)
	return nil, nil
}

//////////////////////////////////////////////////////////////////////////////

// RabbitmqContainer 只給 `rabbitmq` MQ 消費者服務使用。
type RabbitmqContainer struct {

	// pkg
	Clock *pkgUtility.Clock

	// Helper
	*helper.AbstractHelper
	*helper.AesHelper

	// MQ 消費者
	*inputApplicationRabbitmq.AbstractHandler
	ConsumerAdminResourceAppUser *inputApplicationRabbitmqAdminResource.AppUserHandler
}

func InitRabbitmqContainer(ctx context.Context) (*RabbitmqContainer, error) {
	wire.Build(

		// bootstrap
		bootstrap.NewMysql,
		bootstrap.NewAmqp,
		bootstrap.NewRedis,
		pkgCache.NewAop,
		pkgUtility.NewClock,

		// helper
		helper.NewAbstractHelper,
		helper.NewAesHelper,
		helper.NewJwtHelper,
		helper.NewValiatorHelper,

		// output
		outputApplicationMysql.NewAbstractMysql,
		outputApplicationMysqlModel.NewAbstractModel,
		outputApplicationMysqlModel.NewAppUserModel,

		// usecase
		usecaseApplicationAnyAdmin.NewAbstractUsecase,
		usecaseApplicationAnyAdminResource.NewAppUserUsecase,

		// input-consumer
		inputApplicationRabbitmq.NewAbstractHandler,
		inputApplicationRabbitmqAdminResource.NewAppUserHandler,

		wire.Struct(new(RabbitmqContainer), "*"),
	)
	return nil, nil
}

//////////////////////////////////////////////////////////////////////////////

// CronContainer 只給 `cron` 排程服務使用。
type CronContainer struct {

	// pkg
	Clock *pkgUtility.Clock

	// Helper
	*helper.AbstractHelper
	*helper.AesHelper
	*helper.JwtHelper

	// 排程 server
	CronAdminResourceAppUser             *inputApplicationCronAdminResource.AppUserHandler
	CronAdminAuthenticationAuthenticator *inputApplicationCronAdminAuthentication.AuthenticatorHandler
}

func InitCronContainer(ctx context.Context) (*CronContainer, error) {
	wire.Build(

		// bootstrap
		bootstrap.NewMysql,
		bootstrap.NewRedis,
		pkgCache.NewAop,
		pkgUtility.NewClock,

		// helper
		helper.NewAbstractHelper,
		helper.NewAesHelper,
		helper.NewJwtHelper,
		helper.NewValiatorHelper,

		// output
		outputApplicationMysql.NewAbstractMysql,
		outputApplicationMysqlModel.NewAbstractModel,
		outputApplicationMysqlModel.NewAppUserModel,
		outputApplicationMysqlModel.NewAdminUserModel,

		// usecase
		usecaseApplicationAnyAdmin.NewAbstractUsecase,
		usecaseApplicationAnyAdminResource.NewAppUserUsecase,
		usecaseApplicationAnyAdminAuthentication.NewAuthenticatorUsecase,

		// input-cron
		inputApplicationCron.NewAbstractHandler,
		inputApplicationCronAdminResource.NewAppUserHandler,
		inputApplicationCronAdminAuthentication.NewAuthenticatorHandler,

		wire.Struct(new(CronContainer), "*"),
	)
	return nil, nil
}

//////////////////////////////////////////////////////////////////////////////

// WebsocketContainer 只給 `websocket` 服務使用。
type WebsocketContainer struct {

	// pkg
	Clock *pkgUtility.Clock

	// Helper
	*helper.AbstractHelper
	*helper.AesHelper
	*helper.RsaHelper
	*helper.JwtHelper

	// Clients
	ResourceClient *client.ResourceClient

	// websocket
	*inputApplicationWebsocket.AbstractHandler
	WebsocketAdminAuthenticationAuthenticator *inputApplicationWebsocketAdminAuthentication.AuthenticatorHandler

	// Websocket server -Middleware
	WebsocketAdminAbstractMiddleware    *middlewareWebsocketAdmin.AbstractMiddleware
	WebsocketAdminAdminMiddleware       *middlewareWebsocketAdmin.AdminMiddleware
	WebsocketAdminErrorMiddleware       *middlewareWebsocketAdmin.ErrorMiddleware
	WebsocketAdminLoggerMiddleware      *middlewareWebsocketAdmin.LoggerMiddleware
	WebsocketAdminNonexistentMiddleware *middlewareWebsocketAdmin.NonexistentMiddleware
	WebsocketAdminRequestMiddleware     *middlewareWebsocketAdmin.RequestMiddleware
}

func InitWebsocketContainer(ctx context.Context) (*WebsocketContainer, error) {
	wire.Build(

		// pkg
		pkgUtility.NewClock,

		// bootstrap
		bootstrap.NewResource,

		// helper
		helper.NewAbstractHelper,
		helper.NewAesHelper,
		helper.NewRsaHelper,
		helper.NewJwtHelper,
		helper.NewValiatorHelper,

		// output
		outputApplicationResource.NewAbstractResource,
		outputApplicationResourceModel.NewAbstractModel,
		outputApplicationResourceModel.NewAdminUserModel,

		// usecase
		usecaseApplicationAnyAdmin.NewAbstractUsecase,
		usecaseApplicationAnyAdminAuthentication.NewAuthenticatorUsecase,

		// client
		client.NewModel,
		client.NewEvent,
		client.NewLogic,
		client.NewResourceClient,

		// websocket
		inputApplicationWebsocket.NewAbstractHandler,
		inputApplicationWebsocketAdminAuthentication.NewAuthenticatorHandler,

		// Middleware 部分
		middlewareWebsocketAdmin.NewAbstractMiddleware,
		middlewareWebsocketAdmin.NewAdminMiddleware,
		middlewareWebsocketAdmin.NewErrorMiddleware,
		middlewareWebsocketAdmin.NewLoggerMiddleware,
		middlewareWebsocketAdmin.NewNonexistentMiddleware,
		middlewareWebsocketAdmin.NewRequestMiddleware,

		wire.Struct(new(WebsocketContainer), "*"),
	)
	return nil, nil
}

type CentrifugeContainer struct {

	// pkg
	Clock *pkgUtility.Clock

	// Helper
	*helper.AbstractHelper
	*helper.AesHelper
	*helper.RsaHelper
	*helper.JwtHelper

	// Clients
	ResourceClient *client.ResourceClient

	// NATS
	Nats       *nats.Conn
	NatsBroker *pkgNats.NATSBroker
}

func InitCentrifugeContainer(ctx context.Context) (*CentrifugeContainer, error) {
	wire.Build(

		// pkg
		pkgUtility.NewClock,

		// bootstrap
		bootstrap.NewResource,
		bootstrap.NewNats,

		// helper
		helper.NewAbstractHelper,
		helper.NewAesHelper,
		helper.NewRsaHelper,
		helper.NewJwtHelper,

		// client
		client.NewModel,
		client.NewEvent,
		client.NewLogic,
		client.NewResourceClient,

		// centrifuge
		pkgNats.NewNATSBroker,

		wire.Struct(new(CentrifugeContainer), "*"),
	)
	return nil, nil
}

//////////////////////////////////////////////////////////////////////////////

// ClientContainer 只給 `client` （訂閱外部 gRPC stream）服務使用。
type ClientContainer struct {

	// pkg
	Clock *pkgUtility.Clock

	// Helper
	*helper.AbstractHelper
	*helper.AesHelper
}

func InitClientContainer(ctx context.Context) (*ClientContainer, error) {
	wire.Build(

		// pkg
		pkgUtility.NewClock,

		// helper 部份
		helper.NewAbstractHelper,
		helper.NewAesHelper,

		wire.Struct(new(ClientContainer), "*"),
	)
	return nil, nil
}

// 、
type CommandContainer struct {

	// pkg
	Clock *pkgUtility.Clock

	// Helper
	*helper.AbstractHelper
	*helper.AesHelper
	*helper.JwtHelper

	// command
	*inputApplicationCommand.AbstractHandler
	CommandAdminReourceAppUser       *inputApplicationCommandAdminResource.AppUserHandler
	CommandAdminAuthenticationSignIn *inputApplicationCommandAdminAuthentication.AuthenticatorHandler
}

func InitCommandContainer(ctx context.Context) (*CommandContainer, error) {
	wire.Build(

		// bootstrap
		bootstrap.NewMysql,
		bootstrap.NewRedis,
		pkgCache.NewAop,
		pkgUtility.NewClock,

		// helper
		helper.NewAbstractHelper,
		helper.NewAesHelper,
		helper.NewJwtHelper,
		helper.NewValiatorHelper,

		// output
		outputApplicationMysql.NewAbstractMysql,
		outputApplicationMysqlModel.NewAbstractModel,
		outputApplicationMysqlModel.NewAppUserModel,
		outputApplicationMysqlModel.NewAdminUserModel,

		// usecase
		usecaseApplicationAnyAdmin.NewAbstractUsecase,
		usecaseApplicationAnyAdminResource.NewAppUserUsecase,
		usecaseApplicationAnyAdminAuthentication.NewAuthenticatorUsecase,

		// command
		inputApplicationCommand.NewAbstractHandler,
		inputApplicationCommandAdminResource.NewAppUserHandler,
		inputApplicationCommandAdminAuthentication.NewAuthenticatorHandler,

		wire.Struct(new(CommandContainer), "*"),
	)
	return nil, nil
}

type TcpContainer struct {

	// pkg
	Clock *pkgUtility.Clock

	// Helper
	*helper.AbstractHelper
	*helper.AesHelper
	*helper.JwtHelper

	// Clients
	ResourceClient *client.ResourceClient

	// tcp
	*inputApplicationTcp.AbstractHandler
	TcpAdminAuthenticationSignIn *inputApplicationTcpAdminAuthentication.AuthenticatorHandler
}

func InitTcpContainer(ctx context.Context) (*TcpContainer, error) {
	wire.Build(

		// pkg
		pkgUtility.NewClock,

		// bootstrap
		bootstrap.NewResource,

		// helper
		helper.NewAbstractHelper,
		helper.NewAesHelper,
		helper.NewJwtHelper,
		helper.NewValiatorHelper,

		// output
		outputApplicationResource.NewAbstractResource,
		outputApplicationResourceModel.NewAbstractModel,
		outputApplicationResourceModel.NewAdminUserModel,

		// usecase
		usecaseApplicationAnyAdmin.NewAbstractUsecase,
		usecaseApplicationAnyAdminAuthentication.NewAuthenticatorUsecase,

		// client
		client.NewModel,
		client.NewEvent,
		client.NewLogic,
		client.NewResourceClient,

		// tcp
		inputApplicationTcp.NewAbstractHandler,
		inputApplicationTcpAdminAuthentication.NewAuthenticatorHandler,

		wire.Struct(new(TcpContainer), "*"),
	)
	return nil, nil
}

type SourceContainer struct {

	// pkg
	Clock *pkgUtility.Clock

	// Helper
	*helper.AbstractHelper
	*helper.AesHelper

	SourceAnnouncementLottery *inputApplicationSourceAnnouncement.LotteryHandler
}

func InitSourceContainer(ctx context.Context) (*SourceContainer, error) {
	wire.Build(

		// pkg
		pkgUtility.NewClock,

		// bootstrap

		// helper
		helper.NewAbstractHelper,
		helper.NewAesHelper,

		outputApplicationMemory.NewAbstractMemory,
		outputApplicationMemoryModel.NewAbstractModel,
		outputApplicationMemoryModel.NewLotteryModel,
		usecaseApplicationAnyAnnouncement.NewAbstractUsecase,
		usecaseApplicationAnyAnnouncement.NewLotteryUsecase,

		inputApplicationSource.NewAbstractHandler,
		inputApplicationSourceAnnouncement.NewLotteryHandler,

		wire.Struct(new(SourceContainer), "*"),
	)
	return nil, nil
}

type DaemonContainer struct {

	// pkg
	Clock *pkgUtility.Clock

	// Helper
	*helper.AbstractHelper
	*helper.AesHelper

	// Clients
	SourceClient *client.SourceClient

	DaemonWatcherSourceAnnouncementLottery *inputApplicationDaemonWatcherSource.AnnouncementLotteryHandler
}

func InitDaemonContainer(ctx context.Context) (*DaemonContainer, error) {
	wire.Build(

		// pkg
		pkgUtility.NewClock,

		// bootstrap
		bootstrap.NewSource,
		bootstrap.NewRedis,

		// client
		client.NewAnnouncement,
		client.NewSourceClient,

		// helper
		helper.NewAbstractHelper,
		helper.NewAesHelper,
		helper.NewCacheHelper,

		// output
		outputApplicationCache.NewAbstractCache,
		outputApplicationCacheModel.NewAbstractModel,
		outputApplicationCacheModel.NewLotteryModel,

		// usecase
		usecaseApplicationAnyWatcher.NewAbstractUsecase,
		usecaseApplicationAnyWatcherSource.NewAnnouncementLotteryUsecase,

		// input-daemon
		inputApplicationDaemon.NewAbstractHandler,
		inputApplicationDaemonWatcherSource.NewAnnouncementLotteryHandler,

		wire.Struct(new(DaemonContainer), "*"),
	)
	return nil, nil
}
