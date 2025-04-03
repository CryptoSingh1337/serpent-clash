package ecs

type GameServerMetric struct {
	PlayerCount uint8  `json:"playerCount"`
	MemoryUsage uint64 `json:"memoryUsageInMiB"`
}

func NewGameServerMetric() *GameServerMetric {
	return &GameServerMetric{}
}
