package handler

import (
	apiutils "github.com/CryptoSingh1337/serpent-clash/server/internal/api/utils"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/config"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs"
	gameutils "github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/utils"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"net/http"
)

type GameHandler struct {
	game    *ecs.Game
	handler *apiutils.WSHandler
}

func NewGameHandler(g *echo.Group, game *ecs.Game) {
	h := &GameHandler{
		game,
		&apiutils.WSHandler{
			OnOpen: func(conn *websocket.Conn) {
				gameutils.Logger.Info().Msgf("Connection open")
			},
			OnMessage: func(playerId string, messageType int, data []byte) {
				game.Engine.ProcessEvent(playerId, messageType, data)
			},
			OnClose: func(playerId string, err error) {
				if err != nil {
					gameutils.Logger.Err(err).Msgf("Error while closing ws connection, player id: %s", playerId)
				}
				game.Engine.LeaveQueue <- playerId
			},
		},
	}
	g.GET("/ws", h.RegisterPlayer)
	g.POST("/player/:playerId/teleport", h.Teleport)
}

func (h *GameHandler) RegisterPlayer(c echo.Context) error {
	return h.game.AddPlayer(c, h.handler)
}

func (h *GameHandler) Teleport(c echo.Context) error {
	if !config.AppConfig.DebugMode {
		return c.JSON(http.StatusOK, apiutils.CreateResponse[any](nil,
			apiutils.NewError("debug mode is disabled")))
	}
	_ = c.Param("playerId")
	coordinate := new(gameutils.Coordinate)
	if err := c.Bind(coordinate); err != nil {
		return c.JSON(http.StatusBadRequest, apiutils.CreateResponse[any](nil,
			apiutils.NewError("error in deserialization")))
	}
	return c.JSON(http.StatusOK, nil)
	//segments, ok := game.TeleportPlayer(playerId, coordinate)
	//if !ok {
	//	return c.JSON(http.StatusBadRequest, utils.CreateResponse[any](nil,
	//		utils.NewError(fmt.Sprintf("error in teleporting player: %v, to (%v, %v)", playerId, coordinate.X,
	//			coordinate.Y))))
	//}
	//return c.JSON(http.StatusOK, utils.CreateResponse(map[string]any{
	//	"segments": segments,
	//}, nil))
}
