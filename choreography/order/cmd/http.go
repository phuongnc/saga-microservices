package cmd

import (
	"context"
	"fmt"

	"infra/common/log"
	"infra/common/middleware"
	"order-service/src"

	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Api struct {
	echo   *echo.Echo
	config *AppConfig
	logger *log.Logger
}

func NewApi(config *AppConfig, gormDB *gorm.DB, logger *log.Logger, orderHandler src.OrderHandler) *Api {
	a := &Api{
		config: config,
		logger: logger,
	}
	a.echo = createHttpServer(gormDB, orderHandler, a.logger)
	return a
}

func (a *Api) Run() error {
	a.logger.Info("starting order service")
	return a.echo.Start(fmt.Sprintf(":%s", a.config.HttpServerConfig.Port))
}

func (a *Api) Stop() {
	a.logger.Info("Shutdown order service")
	_ = a.echo.Shutdown(context.Background())
}

func createHttpServer(gormDB *gorm.DB, orderHandler src.OrderHandler, log *log.Logger) *echo.Echo {
	e := echo.New()
	e.GET("/health", healthCheck)
	path := e.Group("/orders")

	e.Use(middleware.HttpDb(gormDB))
	e.Use(middleware.Cors())
	orderHandler.RegisterEndpoints(path)
	return e
}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, "Running")
}
