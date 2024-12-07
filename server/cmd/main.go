package main

import (
	"context"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/services"
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
	utils.Logger.LogInfo().Msgf("Loaded config: %v", app.Config)
	game := services.NewGame()
	srv := initHTTPServer(app, game)

	err := srv.Start()
	if err != nil {
		utils.Logger.LogFatal().Msgf("nbio.Start failed: %v", err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	game.Close()
	err = srv.Shutdown(ctx)
	if err != nil {
		utils.Logger.LogFatal().Msgf("Shutdown failed: %v", err)
		return
	}
}
