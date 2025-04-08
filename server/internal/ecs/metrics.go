package ecs

type GameServerMetrics struct {
	PlayerCount uint8  `json:"playerCount"`
	MemoryUsage uint64 `json:"memoryUsageInMiB"`
}

func NewGameServerMetrics() *GameServerMetrics {
	return &GameServerMetrics{}
}
