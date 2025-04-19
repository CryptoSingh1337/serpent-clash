package ecs

import "github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/storage"

type GameServerMetrics struct {
	PlayerCount uint8             `json:"playerCount"`
	MemoryUsage uint64            `json:"memoryUsageInMiB"`
	QuadTree    *storage.QuadTree `json:"quadTree"`
}

func NewGameServerMetrics() *GameServerMetrics {
	return &GameServerMetrics{}
}
