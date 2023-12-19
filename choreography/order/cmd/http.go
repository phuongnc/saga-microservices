package cmd

import (
	"context"
	"fmt"

	"infra/log"
	"infra/middleware"
	"order-service/src"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Api struct {
	echo   *echo.Echo
	config *AppConfig
	logger *log.Logger
}

func NewApi(config *AppConfig, db *gorm.DB, logger *log.Logger, orderHandler src.OrderHandler) *Api {
	a := &Api{
		config: config,
		logger: logger,
	}
	a.echo = createHttpServer(db, orderHandler, a.logger)
	return a
}

func (a *Api) Run() error {
	a.logger.Info("starting order service")
	return a.echo.Start(fmt.Sprintf(":%s", a.config.HttpServerConfig.Port))
}

func (a *Api) Stop() {
	_ = a.echo.Shutdown(context.Background())
}

func createHttpServer(db *gorm.DB, orderHandler src.OrderHandler, log *log.Logger) *echo.Echo {
	e := echo.New()
	v1 := e.Group("/v1")
	e.Use(middleware.HttpDb(db))
	e.Use(middleware.Cors())
	orderHandler.RegisterEndpoints(v1)
	return e
}