package main

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/lesismal/nbio/nbhttp/websocket"
)

func initHandler(e *echo.Echo, app *App, game *services.Game) {
	e.File("/", app.Config.indexFile)
	e.Static("/assets", app.Config.assetDir)
	e.Static("/", app.Config.distDir)

	e.GET("/ws", func(c echo.Context) error {
		return handleWebsocket(c, game)
	})
	e.GET("/*", handleCatchAll)
}

func handleCatchAll(c echo.Context) error {
	app := c.Get("app").(*App)
	return c.File(app.Config.indexFile)
}

func handleWebsocket(c echo.Context, game *services.Game) error {
	w := c.Response()
	r := c.Request()
	upgrader := websocket.NewUpgrader()
	upgrader.EnableCompression(false)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	player := services.NewPlayer(conn)
	upgrader.OnMessage(func(c *websocket.Conn, messageType websocket.MessageType, data []byte) {
		game.ProcessEvent(player, messageType, data)
	})
	game.JoinQueue <- player
	conn.OnClose(func(c *websocket.Conn, err error) {
		game.LeaveQueue <- player
	})
	return nil
}
