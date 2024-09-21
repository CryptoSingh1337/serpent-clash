package main

import (
	"github.com/CryptoSingh1337/multiplayer-snake-game/server/internal/services"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lesismal/nbio/nbhttp"
	"os"
	"path/filepath"
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

func LoadConfig() *Config {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Set default values
	config := Config{
		addr: "0.0.0.0",
		port: "8080",
	}

	// Override with environment variables if they exist
	if env := os.Getenv("SERVER_ADDR"); env != "" {
		config.addr = env
	}
	if env := os.Getenv("SERVER_PORT"); env != "" {
		config.port = env
	}
	if env := os.Getenv("DIST_DIR"); env != "" {
		config.distDir = env
		config.assetDir = filepath.Join(env, "assets")
		config.indexFile = filepath.Join(env, "index.html")
	}
	return &config
}

func initHTTPServer(app *App, game *services.Game) *nbhttp.Engine {
	e := echo.New()

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
