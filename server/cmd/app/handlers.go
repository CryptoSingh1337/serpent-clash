package main

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"net/http"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func initHandler(srv *echo.Echo, app *App, game *Game) {
	srv.File("/", app.Config.indexFile)
	srv.Static("/assets", app.Config.assetDir)
	srv.Static("/", app.Config.distDir)

	srv.GET("/ws", func(c echo.Context) error {
		return handleWebsocket(c, game)
	})
	srv.GET("/*", handleCatchAll)
}

func handleCatchAll(c echo.Context) error {
	app := c.Get("app").(*App)
	return c.File(app.Config.indexFile)
}

func handleWebsocket(c echo.Context, game *Game) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	client := &Client{id: uuid.NewString(), session: nil, conn: conn, send: make(chan []string)}
	err = game.addClient(client)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	if client.session == nil {
		err = errors.New("no session assigned to the client")
		c.Logger().Error(err)
		return err
	}

	body := fmt.Sprintf("{\"clientId\":%q, \"sessionId\":%q}", client.id, client.session.id)
	err = conn.WriteJSON(Payload{"initialize", body})
	if err != nil {
		c.Logger().Error(err)
	}
	go client.read()
	go client.write()
	return nil
}
