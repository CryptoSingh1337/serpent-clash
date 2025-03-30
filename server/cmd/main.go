package main

import (
	"context"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/utils"
	"os"
	"os/signal"
	"time"
)

func main() {
	config := LoadConfig()
	app := &App{
		Config: *config,
	}
	game := ecs.NewGame()
	srv := initHTTPServer(app, game)
	utils.Logger.LogInfo().Msgf("Loaded config: %v", app.Config)

	err := srv.Start()
	game.Start()
	if err != nil {
		utils.Logger.LogFatal().Msgf("nbio.Start failed: %v", err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	game.Stop()
	err = srv.Shutdown(ctx)
	if err != nil {
		utils.Logger.LogFatal().Msgf("Shutdown failed: %v", err)
		return
	}
}
