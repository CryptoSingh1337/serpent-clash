package handler

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs"
	"github.com/labstack/echo/v4"
	"net/http"
)

type MetricsHandler struct {
	game *ecs.Game
}

func NewMetricsHandler(g *echo.Group, game *ecs.Game) {
	h := &MetricsHandler{
		game,
	}
	g.GET("", h.SubScribeGameMetrics)
}

func (h *MetricsHandler) SubScribeGameMetrics(c echo.Context) error {
	return c.JSON(http.StatusOK, h.game.GameServerMetrics)
}
