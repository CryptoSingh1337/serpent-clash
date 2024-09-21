package main

import (
	"context"
	"github.com/CryptoSingh1337/multiplayer-snake-game/server/internal/services"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	config := LoadConfig()
	app := &App{
		Config: *config,
	}
	log.Println("Loaded config", app.Config)
	game := services.NewGame()
	srv := initHTTPServer(app, game)

	err := srv.Start()
	if err != nil {
		log.Fatalf("nbio.Start failed: %v\n", err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	game.Close()
	err = srv.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Shutdown failed: %v\n", err)
		return
	}
}
