package api

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/api/handler"
	apiutils "github.com/CryptoSingh1337/serpent-clash/server/internal/api/utils"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/config"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type Api struct {
	Server *echo.Echo
}

func NewApi(game *ecs.Game) *Api {
	api := &Api{}
	apiutils.NewLogger()
	e := initHttpServer()
	e.HideBanner = true
	e.HidePort = true
	api.Server = e

	g := e.Group("/metrics")
	handler.NewMetricsHandler(g, game)

	g = e.Group("/game")
	handler.NewGameHandler(g, game)

	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, nil)
	})
	e.Static("/", config.AppConfig.DistDir)
	e.Static("/assets", config.AppConfig.AssetDir)
	e.File("/", config.AppConfig.IndexFile)
	e.File("/favicon.png", config.AppConfig.Favicon)
	e.GET("/*", func(c echo.Context) error {
		return c.File(config.AppConfig.IndexFile)
	})
	return api
}

func initHttpServer() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(apiutils.LoggingMiddleware)
	return e
}
