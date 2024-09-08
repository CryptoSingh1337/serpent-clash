package main

import (
	"fmt"
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

func initHandler(srv *echo.Echo, app *App) {
	srv.File("/", app.Config.indexFile)
	srv.Static("/assets", app.Config.assetDir)
	srv.Static("/", app.Config.distDir)

	srv.GET("/ws", handleWebsocket)
	srv.GET("/*", handleCatchAll)
}

func handleCatchAll(c echo.Context) error {
	app := c.Get("app").(*App)
	return c.File(app.Config.indexFile)
}

func handleWebsocket(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()
	for {
		err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, World!"))
		if err != nil {
			c.Logger().Error(err)
		}

		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}
		fmt.Println(string(msg))
	}
}
