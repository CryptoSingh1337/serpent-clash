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
	app := &App{
		Config: Config{
			addr:      "localhost",
			port:      "8080",
			distDir:   "../client/dist",
			assetDir:  "../client/dist/assets",
			indexFile: "../client/dist/index.html",
		},
	}
	game := services.NewGame()
	srv := initHTTPServer(app, game)

	err := srv.Start()
	if err != nil {
		log.Fatalf("nbio.Start failed: %v\n", err)
	}

	log.Println("serving [labstack/echo] on [nbio]")

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
