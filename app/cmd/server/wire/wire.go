//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/spf13/viper"

	"app/internal/handler"
	"app/internal/repository"
	"app/internal/server"
	"app/internal/service"
	"app/pkg/app"
	"app/pkg/jwt"
	"app/pkg/log"
	"app/pkg/server/http"
	"app/pkg/sid"
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	//repository.NewRedis,
	repository.NewRepository,
	repository.NewTransaction,
	repository.NewHomeRepository,
	repository.NewExchangeRepository,
	repository.NewAccountRepository,
	repository.NewSettingsRepository,
	repository.NewResourceRepository,
	repository.NewTurntableRepository,
	repository.NewUserAssetRepository,
	repository.NewUserZodiacRepository,
	repository.NewUserAssetRecordRepository,
	repository.NewPunchingRecordRepository,
	repository.NewAdvertiseRepository,
	repository.NewGameRoomRepository,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewHomeService,
	service.NewExchangeService,
	service.NewAccountService,
	service.NewSettingsService,
	service.NewResourceService,
	service.NewTurntableService,
	service.NewPunchingRecordService,
	service.NewAdvertiseService,
	service.NewGameRoomService,
	service.NewUserAssetService,
)

var handlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewExchangeHandler,
	handler.NewHomeHandler,
	handler.NewAccountHandler,
	handler.NewSettingsHandler,
	handler.NewResourceHandler,
	handler.NewTurntableHandler,
	handler.NewPunchingRecordHandler,
	handler.NewAdvertiseHandler,
	handler.NewGameRoomHandler,
)

var serverSet = wire.NewSet(
	server.NewHTTPServer,
	server.NewJob,
)

// build App
func newApp(
	httpServer *http.Server,
	job *server.Job,
) *app.App {
	return app.NewApp(
		app.WithServer(httpServer, job),
		app.WithName("demo-server"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {
	panic(
		wire.Build(
			repositorySet,
			serviceSet,
			handlerSet,
			serverSet,
			sid.NewSid,
			jwt.NewJwt,
			newApp,
		),
	)
}
