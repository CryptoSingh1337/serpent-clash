package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/CryptoSingh1337/multiplayer-snake-game/server/internal/services"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"log"
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
	// TODO: add readDeadline and writeDeadline in ws connection
	upgrader.OnMessage(func(c *websocket.Conn, messageType websocket.MessageType, data []byte) {
		// TODO: add dispatcher
		log.Println("Message received:", string(data))
	})
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	client := &services.Client{Id: uuid.NewString(), Conn: conn, Send: make(chan []string)}
	err = game.AddClient(client)
	if err != nil {
		c.Logger().Error(err)
		conn.CloseWithError(err)
	}
	if client.Session == nil {
		err = errors.New("no session assigned to the client")
		c.Logger().Error(err)
		conn.CloseWithError(err)
	}
	body := fmt.Sprintf("{\"clientId\":%q, \"sessionId\":%q}", client.Id, client.Session.Id)
	payload, err := json.Marshal(services.Payload{Type: "initialize", Body: []byte(body)})
	err = conn.WriteMessage(websocket.TextMessage, payload)
	if err != nil {
		conn.CloseWithError(err)
	}
	conn.OnClose(func(c *websocket.Conn, err error) {
		log.Println("OnClose:", c.RemoteAddr().String(), err)
		err = game.RemoveClient(client)
		if err != nil {
			log.Println("Error removing client:", err)
		}
	})
	log.Println("OnOpen:", conn.RemoteAddr().String())
	return nil
}
