package handler

import (
	"net/http"
	"time"

	"github.com/CryptoSingh1337/serpent-clash/server/internal/api/utils"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/storage"
	"github.com/labstack/echo/v4"
)

type MetricsHandler struct {
	game *ecs.Game
}

func NewMetricsHandler(g *echo.Group, game *ecs.Game) {
	h := &MetricsHandler{
		game,
	}
	g.GET("/subscribe/info", h.SubscribeMetricsInfo)
	g.GET("/subscribe/quad-tree", h.SubscribeQuadTreeInfo)
}

func (h *MetricsHandler) SubscribeMetricsInfo(c echo.Context) error {
	utils.Logger.Debug().Msgf("SSE client connected, ip: %v", c.RealIP())
	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-c.Request().Context().Done():
			utils.Logger.Debug().Msgf("SSE client disconnected, ip: %v", c.RealIP())
			return nil
		case <-ticker.C:
			p, err := utils.ToJsonB(struct {
				ServerMetrics ecs.ServerMetrics `json:"serverMetrics"`
				GameMetrics   ecs.GameMetrics   `json:"gameMetrics"`
			}{
				ServerMetrics: h.game.GameServerMetrics.ServerMetrics,
				GameMetrics:   h.game.GameServerMetrics.GameMetrics,
			})
			if err != nil {
				continue
			}
			event := utils.Event{
				Data: p,
			}
			if err := event.MarshalTo(w); err != nil {
				return err
			}
			w.Flush()
		}
	}
}

func (h *MetricsHandler) SubscribeQuadTreeInfo(c echo.Context) error {
	utils.Logger.Debug().Msgf("SSE client connected, ip: %v", c.RealIP())
	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-c.Request().Context().Done():
			utils.Logger.Debug().Msgf("SSE client disconnected, ip: %v", c.RealIP())
			return nil
		case <-ticker.C:
			p, err := utils.ToJsonB(struct {
				QuadTree     *storage.QuadTree `json:"quadTree"`
				SpawnRegions ecs.SpawnRegions  `json:"spawnRegions"`
			}{
				QuadTree:     h.game.GameServerMetrics.QuadTree,
				SpawnRegions: h.game.GameServerMetrics.SpawnRegions,
			})
			if err != nil {
				continue
			}
			event := utils.Event{
				Data: p,
			}
			if err := event.MarshalTo(w); err != nil {
				return err
			}
			w.Flush()
		}
	}
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
