package main

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"app/cmd/server/wire"
	"app/pkg/config"
	"app/pkg/log"

	"github.com/joho/godotenv"
)

// @title           鸵小妥API
// @version         1.0.2
// @contact.email  wq0426@163.com
// @host      127.0.0.1:8289
// @securityDefinitions.apiKey Bearer
// @in header
// @name Authorization
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	conf := config.NewConfig()
	logger := log.NewLog(conf)
	app, cleanup, err := wire.NewWire(conf, logger)
	if err != nil {
		panic(err)
	}
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}
	defer cleanup()
	logger.Info(
		"server start",
		zap.String("host", fmt.Sprintf("http://%s:%d", conf.GetString("http.host"), conf.GetInt("http.port"))),
	)
	if err = app.Run(context.Background()); err != nil {
		panic(err)
	}
}
