//go:build wireinject
// +build wireinject

package wire

import (
	"app/internal/handler"
	"app/internal/repository"
	"app/internal/server"
	"app/internal/service"
	"app/pkg/app"
	"app/pkg/log"

	"github.com/google/wire"
	"github.com/spf13/viper"

	"app/pkg/jwt"
	"app/pkg/sid"
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	repository.NewRepository,
	repository.NewTransaction,
	repository.NewActivityRepository,
	repository.NewSettingsRepository,
	repository.NewAccountRepository,
	repository.NewEscortRecordRepository,
	repository.NewUnionMemberRepository,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewActivityService,
)

var handlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewActivityHandler,
)

var serverSet = wire.NewSet(
	server.NewTask,
)

// build App
func newApp(
	task *server.Task,
) *app.App {
	return app.NewApp(
		app.WithServer(task),
		app.WithName("demo-task"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {
	panic(
		wire.Build(
			repositorySet,
			serviceSet,
			serverSet,
			handlerSet,
			sid.NewSid,
			jwt.NewJwt,
			newApp,
		),
	)
}
