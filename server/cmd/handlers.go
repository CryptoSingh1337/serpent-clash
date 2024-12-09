package main

import (
	"errors"
	"fmt"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/services"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/lesismal/llib/std/net/http"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"strings"
)

func initHandler(e *echo.Echo, app *App, game *services.GameDriver) {
	e.File("/", app.Config.indexFile)
	e.Static("/assets", app.Config.assetDir)
	e.Static("/", app.Config.distDir)

	e.GET("/ws", func(c echo.Context) error {
		return handleWebsocket(c, game)
	})
	e.POST("/player/:playerId/teleport", func(c echo.Context) error {
		return handlePlayerTeleport(c, game)
	})
	e.GET("/*", handleCatchAll)
}

func handleCatchAll(c echo.Context) error {
	app := c.Get("app").(*App)
	return c.File(app.Config.indexFile)
}

func handleWebsocket(c echo.Context, game *services.GameDriver) error {
	username := c.QueryParam("username")
	username = strings.TrimSpace(username)
	if username == "" {
		return errors.New("invalid username")
	}
	w := c.Response()
	r := c.Request()
	upgrader := websocket.NewUpgrader()
	upgrader.EnableCompression(false)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	player := services.NewPlayer(&username, conn, w)
	upgrader.OnMessage(func(c *websocket.Conn, messageType websocket.MessageType, data []byte) {
		game.ProcessEvent(player, messageType, data)
	})
	game.JoinQueue <- player
	conn.OnClose(func(c *websocket.Conn, err error) {
		game.LeaveQueue <- player
	})
	return nil
}

func handlePlayerTeleport(c echo.Context, game *services.GameDriver) error {
	app := c.Get("app").(*App)
	if !app.Config.debugMode {
		return c.JSON(http.StatusOK, utils.CreateResponse[any](nil,
			utils.NewError("debug mode is disabled")))
	}
	playerId := c.Param("playerId")
	coordinate := new(utils.Coordinate)
	if err := c.Bind(coordinate); err != nil {
		return c.JSON(http.StatusBadRequest, utils.CreateResponse[any](nil,
			utils.NewError("error in deserialization")))
	}
	segments, ok := game.TeleportPlayer(playerId, coordinate)
	if !ok {
		return c.JSON(http.StatusBadRequest, utils.CreateResponse[any](nil,
			utils.NewError(fmt.Sprintf("error in teleporting player: %v, to (%v, %v)", playerId, coordinate.X,
				coordinate.Y))))
	}
	return c.JSON(http.StatusOK, utils.CreateResponse(map[string]any{
		"segments": segments,
	}, nil))
}
