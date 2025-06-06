package main

import (
	"context"
	"errors"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/api"
	apiutils "github.com/CryptoSingh1337/serpent-clash/server/internal/api/utils"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/config"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	Api      *api.Api
	Game     *ecs.Game
	shutdown chan os.Signal
}

func LoadConfig() {
	config.Init()
}

func NewApp() *App {
	LoadConfig()
	app := &App{}
	app.Game = ecs.NewGame()
	app.Api = api.NewApi(app.Game)
	app.shutdown = make(chan os.Signal, 1)
	signal.Notify(app.shutdown, os.Interrupt, syscall.SIGTERM)
	return app
}

func (app *App) Start() {
	go func() {
		app.Game.Start()
		if err := app.Api.Server.Start(":" + config.AppConfig.Port); err != nil && errors.Is(err, http.ErrServerClosed) {
			apiutils.Logger.Fatal().Msgf("shutting down the server, err: %v", err)
		}
	}()
}

func (app *App) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	app.Game.Stop()
	return app.Api.Server.Shutdown(ctx)
}
