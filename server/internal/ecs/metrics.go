package ecs

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/storage"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/utils"
)

type GameServerMetrics struct {
	PlayerCount  uint8             `json:"playerCount"`
	MemoryUsage  uint64            `json:"memoryUsageInMiB"`
	QuadTree     *storage.QuadTree `json:"quadTree"`
	SpawnRegions struct {
		Radius  float64            `json:"radius"`
		Regions []utils.Coordinate `json:"regions"`
	} `json:"spawnRegions"`
}

func NewGameServerMetrics() *GameServerMetrics {
	return &GameServerMetrics{}
}
