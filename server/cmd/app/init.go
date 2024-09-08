package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Config struct {
	port      string
	distDir   string
	assetDir  string
	indexFile string
}

type App struct {
	Config Config
}

func initLogger(srv *echo.Echo) {
	srv.Logger.SetLevel(2)
}

func initHTTPServer(app *App) *echo.Echo {
	srv := echo.New()

	srv.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			if c.Request().RequestURI == "/ws" {
				return false
			}
			return true
		},
	}))
	srv.Use(middleware.Recover())

	// Register app (*App) to be injected into all HTTP handlers.
	srv.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("app", app)
			return next(c)
		}
	})

	initHandler(srv, app)
	return srv
}
