package main

import (
	"github.com/CryptoSingh1337/multiplayer-snake-game/server/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lesismal/nbio/nbhttp"
)

type Config struct {
	addr      string
	port      string
	distDir   string
	assetDir  string
	indexFile string
}

type App struct {
	Config Config
}

func initLogger(e *echo.Echo) {
	e.Logger.SetLevel(2)
}

func initHTTPServer(app *App) *nbhttp.Engine {
	e := echo.New()
	game := services.NewGame()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			if c.Request().RequestURI == "/ws" {
				return false
			}
			return true
		},
	}))
	e.Use(middleware.Recover())

	// Register app (*App) to be injected into all HTTP handlers.
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("app", app)
			return next(c)
		}
	})

	initLogger(e)

	initHandler(e, app, game)

	return nbhttp.NewEngine(nbhttp.Config{
		Network: "tcp",
		Addrs:   []string{app.Config.addr + ":" + app.Config.port},
		Handler: e,
	})
}
