package main

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/services"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lesismal/nbio/nbhttp"
	"github.com/rs/zerolog"
	"os"
	"path/filepath"
	"strings"
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
	utils.NewLogger()
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	e.Use(utils.LoggingMiddleware)
}

func LoadConfig() *Config {
	env := os.Getenv("GO_ENV")

	// Set default values
	config := Config{
		addr: "0.0.0.0",
		port: "8080",
	}
	if env == "PROD" {
		// Read from environment variables
		if env := os.Getenv("SERVER_ADDR"); env != "" {
			config.addr = env
		}
		if env := os.Getenv("SERVER_PORT"); env != "" {
			config.port = env
		}
		if env := os.Getenv("DIST_DIR"); env != "" {
			config.distDir = env
			config.assetDir = filepath.Join(config.distDir, "assets")
			config.indexFile = filepath.Join(config.distDir, "index.html")
		}
	} else {
		data, err := os.ReadFile(".env")
		if err != nil {
			utils.Logger.LogError().Err(err)
		}
		content := string(data)
		lines := strings.Split(content, "\n")
		for _, line := range lines {
			env := strings.Split(line, "=")
			if env[0] == "DIST_DIR" {
				config.distDir = strings.TrimSuffix(env[1], "\r")
				config.assetDir = filepath.Join(config.distDir, "assets")
				config.indexFile = filepath.Join(config.distDir, "index.html")
			}
		}
	}
	return &config
}

func initHTTPServer(app *App, game *services.Game) *nbhttp.Engine {
	e := echo.New()
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
