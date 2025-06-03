package ecs

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/storage"
	gameutils "github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/utils"
	"github.com/shirou/gopsutil/v4/process"
	"os"
)

type ServerMetrics struct {
	CpuUsage           uint64 `json:"cpuUsage"`
	MemoryUsage        uint64 `json:"memoryUsage"`
	HeapInUse          uint64 `json:"heapAllocated"`
	HeapReserved       uint64 `json:"heapReserved"`
	TotalHeapAllocated uint64 `json:"totalHeapAllocated"`
	HeapObjects        uint64 `json:"heapObjects"`
	LastGCMs           uint64 `json:"lastGCMs"`
	GCPauseMicro       uint64 `json:"gcPauseMicro"`
	NumGoroutines      uint64 `json:"numGoroutines"`
	Uptime             uint64 `json:"uptimeInSec"`
	BytesSent          uint64 `json:"bytesSent"`
	BytesReceived      uint64 `json:"bytesReceived"`
	PacketsSent        uint64 `json:"packetsSent"`
	PacketsReceived    uint64 `json:"packetsReceived"`
	ErrorIn            uint64 `json:"errorIn"`
	ErrorOut           uint64 `json:"errorOut"`
	DropIn             uint64 `json:"dropIn"`
	DropOut            uint64 `json:"dropOut"`
	ActiveConnections  uint8  `json:"activeConnections"`
}

type GameMetrics struct {
	PlayerCount                    uint8    `json:"playerCount"`
	SystemUpdateTimeInLastTick     int64    `json:"systemUpdateTimeInLastTick"`
	MaxSystemUpdateTime            int64    `json:"maxSystemUpdateTime"`
	SystemUpdateTimeInLastTenTicks []int64  `json:"systemUpdateTimeInLastTenTicks"`
	NoOfCollisionsInLastTenTicks   []uint64 `json:"noOfCollisionsInLastTenTicks"`
}

type SpawnRegions struct {
	Radius  float64                `json:"radius"`
	Regions []gameutils.Coordinate `json:"regions"`
}

type GameServerMetrics struct {
	proc          *process.Process
	ServerMetrics ServerMetrics     `json:"serverMetrics"`
	GameMetrics   GameMetrics       `json:"gameMetrics"`
	QuadTree      *storage.QuadTree `json:"quadTree"`
	SpawnRegions  SpawnRegions      `json:"spawnRegions"`
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
