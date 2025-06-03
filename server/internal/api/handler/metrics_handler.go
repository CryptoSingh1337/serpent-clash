package handler

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/storage"
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
	g.GET("/info", h.GetInfo)
	g.GET("/quad-tree", h.GetQuadTree)
}

func (h *MetricsHandler) GetInfo(c echo.Context) error {
	metrics := struct {
		ServerMetrics ecs.ServerMetrics `json:"serverMetrics"`
		GameMetrics   ecs.GameMetrics   `json:"gameMetrics"`
	}{
		ServerMetrics: h.game.GameServerMetrics.ServerMetrics,
		GameMetrics:   h.game.GameServerMetrics.GameMetrics,
	}
	return c.JSON(http.StatusOK, metrics)
}

func (h *MetricsHandler) GetQuadTree(c echo.Context) error {
	quadTree := struct {
		QuadTree     *storage.QuadTree `json:"quadTree"`
		SpawnRegions ecs.SpawnRegions  `json:"spawnRegions"`
	}{
		QuadTree:     h.game.GameServerMetrics.QuadTree,
		SpawnRegions: h.game.GameServerMetrics.SpawnRegions,
	}
	return c.JSON(http.StatusOK, quadTree)
}
