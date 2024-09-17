package main

import (
	"context"
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
	srv := initHTTPServer(app)

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
	err = srv.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Shutdown failed: %v\n", err)
		return
	}
}
