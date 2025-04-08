package main

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs"
	"github.com/labstack/echo/v4"
	"net/http"
)

func initHandler(e *echo.Echo, app *App, game *ecs.Game) {
	e.File("/", app.Config.indexFile)
	e.Static("/assets", app.Config.assetDir)
	e.Static("/", app.Config.distDir)

	e.GET("/ws", func(c echo.Context) error {
		return handleWebsocket(c, game)
	})
	//e.POST("/player/:playerId/teleport", func(c echo.Context) error {
	//	return handlePlayerTeleport(c, game)
	//})
	e.GET("/metrics", func(c echo.Context) error {
		return c.JSON(http.StatusOK, game.GameServerMetrics)
	})
	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, nil)
	})
	e.GET("/*", handleCatchAll)
}

func handleCatchAll(c echo.Context) error {
	app := c.Get("app").(*App)
	return c.File(app.Config.indexFile)
}

func handleWebsocket(c echo.Context, game *ecs.Game) error {
	return game.AddPlayer(c)
}

//func handlePlayerTeleport(c echo.Context, game *ecs.Game) error {
//	app := c.Get("app").(*App)
//	if !app.Config.debugMode {
//		return c.JSON(http.StatusOK, utils.CreateResponse[any](nil,
//			utils.NewError("debug mode is disabled")))
//	}
//	playerId := c.Param("playerId")
//	coordinate := new(utils.Coordinate)
//	if err := c.Bind(coordinate); err != nil {
//		return c.JSON(http.StatusBadRequest, utils.CreateResponse[any](nil,
//			utils.NewError("error in deserialization")))
//	}
//	segments, ok := game.TeleportPlayer(playerId, coordinate)
//	if !ok {
//		return c.JSON(http.StatusBadRequest, utils.CreateResponse[any](nil,
//			utils.NewError(fmt.Sprintf("error in teleporting player: %v, to (%v, %v)", playerId, coordinate.X,
//				coordinate.Y))))
//	}
//	return c.JSON(http.StatusOK, utils.CreateResponse(map[string]any{
//		"segments": segments,
//	}, nil))
//}
