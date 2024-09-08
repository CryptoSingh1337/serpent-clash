package main

import (
	"errors"
	"net/http"
)

func main() {
	app := &App{
		Config: Config{
			port:      ":8080",
			distDir:   "../client/dist",
			assetDir:  "../client/dist/assets",
			indexFile: "../client/dist/index.html",
		},
	}
	srv := initHTTPServer(app)

	initLogger(srv)

	if err := srv.Start(app.Config.port); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			srv.Logger.Printf("HTTP server shut down")
		} else {
			srv.Logger.Printf("HTTP server error: %v", err)
		}
	}
}
