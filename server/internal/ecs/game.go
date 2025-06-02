package ecs

import (
	"errors"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/storage"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/types"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/net"
	"math"
	"strings"
	"time"
)

type Game struct {
	Done              chan bool
	Engine            *Engine
	GameServerMetrics *GameServerMetrics
}

func NewGame() *Game {
	utils.NewLogger()
	return &Game{
		Done:              make(chan bool),
		Engine:            NewEngine(),
		GameServerMetrics: NewGameServerMetrics(),
	}
}

func (g *Game) Start() {
	ticker := time.NewTicker(1000 / utils.TickRate * time.Millisecond)
	metricsTicker := time.NewTicker(1 * time.Second)
	g.Engine.Start()
	r := g.Engine.storage.GetSharedResource(utils.SpawnRegions)
	if r != nil {
		g.GameServerMetrics.SpawnRegions.Radius = utils.SpawnRegionRadius
		g.GameServerMetrics.SpawnRegions.Regions = r.([]utils.Coordinate)
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
	g.Engine.UpdateSystems()
	end := time.Now().UnixMicro()
	g.GameServerMetrics.GameMetrics.SystemUpdateTimeInLastTick = end - start
}

func (g *Game) processMetrics() {
	cpuPercent, _ := g.GameServerMetrics.proc.CPUPercent()
	memInfo, _ := g.GameServerMetrics.proc.MemoryInfo()
	uptime, _ := host.Uptime()
	netStats, _ := net.IOCounters(false)
	g.GameServerMetrics.ServerMetrics.CpuUsage = uint64(cpuPercent)
	g.GameServerMetrics.ServerMetrics.MemoryUsage = memInfo.RSS / (1024 * 1024)
	g.GameServerMetrics.ServerMetrics.Uptime = uptime
	if len(netStats) > 0 {
		g.GameServerMetrics.ServerMetrics.BytesSent = netStats[0].BytesSent
		g.GameServerMetrics.ServerMetrics.BytesReceived = netStats[0].BytesRecv
	}
	g.GameServerMetrics.ServerMetrics.PlayerCount = uint8(len(g.Engine.playerIdToEntityId))
	r := g.Engine.storage.GetSharedResource(utils.QuadTreeResource)
	if r != nil {
		g.GameServerMetrics.QuadTree = r.(*storage.QuadTree)
	}
	g.GameServerMetrics.GameMetrics.MaxSystemUpdateTime = int64(math.Max(
		float64(g.GameServerMetrics.GameMetrics.MaxSystemUpdateTime),
		float64(g.GameServerMetrics.GameMetrics.SystemUpdateTimeInLastTick),
	))
	if len(g.GameServerMetrics.GameMetrics.SystemUpdateTimeInLastTenTicks) < 10 {
		g.GameServerMetrics.GameMetrics.SystemUpdateTimeInLastTenTicks = append(
			g.GameServerMetrics.GameMetrics.SystemUpdateTimeInLastTenTicks,
			g.GameServerMetrics.GameMetrics.SystemUpdateTimeInLastTick)
	} else {
		g.GameServerMetrics.GameMetrics.SystemUpdateTimeInLastTenTicks =
			g.GameServerMetrics.GameMetrics.SystemUpdateTimeInLastTenTicks[1:]
		g.GameServerMetrics.GameMetrics.SystemUpdateTimeInLastTenTicks = append(
			g.GameServerMetrics.GameMetrics.SystemUpdateTimeInLastTenTicks,
			g.GameServerMetrics.GameMetrics.SystemUpdateTimeInLastTick,
		)
	}
}

func (g *Game) AddPlayer(c echo.Context) error {
	username := c.QueryParam("username")
	username = strings.TrimSpace(username)
	if username == "" {
		return errors.New("invalid username")
	}
	w := c.Response()
	r := c.Request()
	upgrader := websocket.NewUpgrader()
	upgrader.EnableCompression(false)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	playerId := uuid.NewString()
	conn.OnMessage(func(c *websocket.Conn, messageType websocket.MessageType, data []byte) {
		g.Engine.ProcessEvent(playerId, messageType, data)
	})
	g.Engine.JoinQueue <- &types.JoinEvent{
		Connection: conn,
		PlayerId:   playerId,
		Username:   username,
	}
	conn.OnClose(func(c *websocket.Conn, err error) {
		g.Engine.LeaveQueue <- playerId
	})
	return nil
}
