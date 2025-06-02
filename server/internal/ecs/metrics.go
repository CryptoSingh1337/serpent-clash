package ecs

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/storage"
	gameutils "github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/utils"
	"github.com/shirou/gopsutil/v4/process"
	"os"
)

type GameServerMetrics struct {
	proc          *process.Process
	ServerMetrics struct {
		CpuUsage      uint64 `json:"cpuUsage"`
		MemoryUsage   uint64 `json:"memoryUsageInMB"`
		Uptime        uint64 `json:"uptimeInSec"`
		BytesSent     uint64 `json:"bytesSent"`
		BytesReceived uint64 `json:"bytesReceived"`
		PlayerCount   uint8  `json:"playerCount"`
	} `json:"serverMetrics"`
	QuadTree     *storage.QuadTree `json:"quadTree"`
	SpawnRegions struct {
		Radius  float64                `json:"radius"`
		Regions []gameutils.Coordinate `json:"regions"`
	} `json:"spawnRegions"`
	GameMetrics struct {
		SystemUpdateTimeInLastTick     int64    `json:"systemUpdateTimeInLastTick"`
		MaxSystemUpdateTime            int64    `json:"maxSystemUpdateTime"`
		SystemUpdateTimeInLastTenTicks []int64  `json:"systemUpdateTimeInLastTenTicks"`
		NoOfCollisionsInLastTenTicks   []uint64 `json:"noOfCollisionsInLastTenTicks"`
	} `json:"gameMetrics"`
}

func NewGameServerMetrics() *GameServerMetrics {
	pid := int32(os.Getpid())
	proc, err := process.NewProcess(pid)
	if err != nil {
		gameutils.Logger.Err(err).Msgf("error while fetching process metrics")
		panic(err)
	}
	metrics := &GameServerMetrics{}
	metrics.GameMetrics.SystemUpdateTimeInLastTenTicks = make([]int64, 0, 10)
	metrics.proc = proc
	return metrics
}
