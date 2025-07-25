package ecs

import (
	"errors"
	apiutils "github.com/CryptoSingh1337/serpent-clash/server/internal/api/utils"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/storage"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/types"
	gameutils "github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/utils"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/net"
	"math"
	"net/http"
	"runtime"
	"runtime/debug"
	"strings"
	"sync/atomic"
	"time"
)

type Game struct {
	Done              chan bool
	Engine            *Engine
	GameServerMetrics *GameServerMetrics
	gcThreshold       int64
}

func NewGame() *Game {
	gameutils.NewLogger()
	return &Game{
		Done:              make(chan bool),
		Engine:            NewEngine(),
		GameServerMetrics: NewGameServerMetrics(),
		gcThreshold:       20000,
	}
}

func (g *Game) Start() {
	old := debug.SetGCPercent(-1) // Disable GC
	debug.SetMemoryLimit(364904448)
	gameutils.Logger.Info().Msgf("Old GC Percent: %v", old)
	ticker := time.NewTicker(1000 / gameutils.TickRate * time.Millisecond)
	metricsTicker := time.NewTicker(1 * time.Second)
	g.Engine.Start()
	r := g.Engine.storage.GetSharedResource(gameutils.SpawnRegions)
	if r != nil {
		g.GameServerMetrics.SpawnRegions.Radius = gameutils.SpawnRegionRadius
		g.GameServerMetrics.SpawnRegions.Regions = r.([]gameutils.Coordinate)
	}
	go func() {
		for {
			select {
			case <-g.Done:
				ticker.Stop()
				return
			case _ = <-ticker.C:
				g.processTick()
			}
		}
	}()
	go func() {
		for {
			select {
			case <-g.Done:
				metricsTicker.Stop()
				return
			case _ = <-metricsTicker.C:
				g.processMetrics()
			}
		}
	}()
}

func (g *Game) Stop() {
	g.Engine.Stop()
	g.Done <- true
	close(g.Done)
}

func (g *Game) processTick() {
	start := time.Now().UnixMicro()
	g.Engine.UpdateSystems(g.GameServerMetrics.GameMetrics.SystemMetrics)
	end := time.Now().UnixMicro()
	atomic.StoreInt64(&g.GameServerMetrics.GameMetrics.SystemUpdateTimeInLastTick, end-start)
	if g.gcThreshold == 0 {
		runtime.GC()
		g.gcThreshold = 10000 // Reset to default if it was set to 0
	}
}

func (g *Game) processMetrics() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	g.GameServerMetrics.ServerMetrics.HeapInUse = m.HeapAlloc
	g.GameServerMetrics.ServerMetrics.HeapReserved = m.HeapSys
	g.GameServerMetrics.ServerMetrics.TotalHeapAllocated = m.TotalAlloc
	g.GameServerMetrics.ServerMetrics.HeapObjects = m.HeapObjects
	if m.NumGC > 0 {
		g.GameServerMetrics.ServerMetrics.LastGCMs = uint64(time.Since(time.Unix(0, int64(m.LastGC))).Milliseconds())
	} else {
		g.GameServerMetrics.ServerMetrics.LastGCMs = 0
	}
	g.GameServerMetrics.ServerMetrics.GCPauseMicro = uint64(time.Duration(m.PauseNs[(m.NumGC+255)%256]).Microseconds())
	g.GameServerMetrics.ServerMetrics.NumGoroutines = uint64(runtime.NumGoroutine())

	cpuPercent, _ := g.GameServerMetrics.proc.CPUPercent()
	memInfo, _ := g.GameServerMetrics.proc.MemoryInfo()
	uptime, _ := host.Uptime()
	netStats, _ := net.IOCounters(false)
	g.GameServerMetrics.ServerMetrics.CpuUsage = uint64(cpuPercent)
	g.GameServerMetrics.ServerMetrics.MemoryUsage = memInfo.RSS
	g.GameServerMetrics.ServerMetrics.Uptime = uptime
	if len(netStats) > 0 {
		g.GameServerMetrics.ServerMetrics.BytesSent = netStats[0].BytesSent
		g.GameServerMetrics.ServerMetrics.BytesReceived = netStats[0].BytesRecv
		g.GameServerMetrics.ServerMetrics.PacketsSent = netStats[0].PacketsSent
		g.GameServerMetrics.ServerMetrics.PacketsReceived = netStats[0].PacketsRecv
		g.GameServerMetrics.ServerMetrics.ErrorIn = netStats[0].Errin
		g.GameServerMetrics.ServerMetrics.ErrorOut = netStats[0].Errout
		g.GameServerMetrics.ServerMetrics.DropIn = netStats[0].Dropin
		g.GameServerMetrics.ServerMetrics.DropOut = netStats[0].Dropout
	}
	g.GameServerMetrics.ServerMetrics.ActiveConnections = uint8(len(g.Engine.playerIdToEntityId))
	r := g.Engine.storage.GetSharedResource(gameutils.QuadTreeResource)
	if r != nil {
		g.GameServerMetrics.QuadTree = r.(*storage.QuadTree)
	}
	current := atomic.LoadInt64(&g.GameServerMetrics.GameMetrics.SystemUpdateTimeInLastTick)
	g.GameServerMetrics.GameMetrics.PlayerCount = g.GameServerMetrics.ServerMetrics.ActiveConnections
	atomic.StoreInt64(&g.GameServerMetrics.GameMetrics.MaxSystemUpdateTime,
		int64(math.Max(
			float64(atomic.LoadInt64(&g.GameServerMetrics.GameMetrics.MaxSystemUpdateTime)),
			float64(current),
		)))
	if len(g.GameServerMetrics.GameMetrics.SystemUpdateTimeInLastTenTicks) < 10 {
		g.GameServerMetrics.GameMetrics.SystemUpdateTimeInLastTenTicks = append(
			g.GameServerMetrics.GameMetrics.SystemUpdateTimeInLastTenTicks,
			current)
	} else {
		g.GameServerMetrics.GameMetrics.SystemUpdateTimeInLastTenTicks =
			g.GameServerMetrics.GameMetrics.SystemUpdateTimeInLastTenTicks[1:]
		g.GameServerMetrics.GameMetrics.SystemUpdateTimeInLastTenTicks = append(
			g.GameServerMetrics.GameMetrics.SystemUpdateTimeInLastTenTicks,
			current,
		)
	}
	for _, s := range g.GameServerMetrics.GameMetrics.SystemMetrics {
		current = atomic.LoadInt64(&s.UpdateTimeInLastTick)
		atomic.StoreInt64(&s.MaxUpdateTime, int64(math.Max(float64(atomic.LoadInt64(&s.MaxUpdateTime)), float64(current))))
		if len(s.UpdateTimeInLastTenTicks) < 10 {
			s.UpdateTimeInLastTenTicks = append(s.UpdateTimeInLastTenTicks, current)
		} else {
			s.UpdateTimeInLastTenTicks = s.UpdateTimeInLastTenTicks[1:]
			s.UpdateTimeInLastTenTicks = append(s.UpdateTimeInLastTenTicks, current)
		}
	}
}

func (g *Game) AddPlayer(c echo.Context, h *apiutils.WSHandler) error {
	username := c.QueryParam("username")
	username = strings.TrimSpace(username)
	if username == "" {
		return errors.New("invalid username")
	}
	w := c.Response()
	r := c.Request()
	upgrader := websocket.Upgrader{
		CheckOrigin:       func(r *http.Request) bool { return true },
		EnableCompression: true,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	if h.OnOpen != nil {
		h.OnOpen(conn)
	}
	playerId := uuid.NewString()
	defer func() {
		if h.OnClose != nil {
			h.OnClose(playerId, nil)
		}
		_ = conn.Close()
	}()
	g.Engine.JoinQueue <- &types.JoinEvent{
		Connection: conn,
		PlayerId:   playerId,
		Username:   username,
	}
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		if h.OnMessage != nil {
			h.OnMessage(playerId, messageType, message)
		}
	}
	return nil
}
