package main

import (
	"context"

	"google.golang.org/grpc"

	"app/cmd/task/wire"
	lgrpc "app/internal/grpc"
	"app/pkg/config"
	"app/pkg/log"
)

func main() {
	conf := config.NewConfig()
	logger := log.NewLog(conf)
	logger.Info("start task")
	stream := GetStream(logger)
	app, cleanup, err := wire.NewWire(conf, stream, logger)
	defer cleanup()
	if err != nil {
		panic(err)
	}

	if err = app.Run(context.Background()); err != nil {
		panic(err)
	}

}

func GetStream(localLog *log.Logger) *lgrpc.PushMessageService_StreamMessagesClient {
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		localLog.Info("Did not connect: " + err.Error())
		return nil
	}
	defer conn.Close()
	// 创建 gRPC 客户端
	client := lgrpc.NewPushMessageServiceClient(conn)
	// 创建流式请求，启动流连接
	stream, err := client.StreamMessages(context.Background())
	if err != nil {
		localLog.Debug("Failed to start stream: %v", err)
		return nil
	}
	return &stream
}
