// @program:     countrybattle
// @file:        golbal.go.go
// @author:      ac
// @create:      2024-11-08 09:50
// @description:
package global

import (
	"github.com/spf13/viper"

	pb "app/internal/grpc"
	"app/internal/handler"
	"app/internal/repository"
	"app/internal/service"
	"app/pkg/log"
)

var HomeHandler *handler.HomeHandler

func InitGlobals(viperViper *viper.Viper, logger *log.Logger) *handler.HomeHandler {
	if HomeHandler == nil {
		v := repository.NewDB(viperViper, logger)
		repositoryRepository := repository.NewRepository(logger, v)
		homeRepository := repository.NewHomeRepository(repositoryRepository)
		handlerHandler := handler.NewHandler(logger)
		serviceService := service.NewService(nil, logger, nil, nil)
		homeService := service.NewHomeService(serviceService, homeRepository)
		HomeHandler = handler.NewHomeHandler(handlerHandler, homeService)
		HomeHandler.EventChannel = make(map[string]chan *pb.PushMessageResponse)
		HomeHandler.EventHistoryChannel = make(map[string]pb.PushMessageResponse)
		HomeHandler.Streams = make(map[string]pb.PushMessageService_StreamMessagesServer)
	}
	return HomeHandler
}

func GetHomeHandler() *handler.HomeHandler {
	return HomeHandler
}
