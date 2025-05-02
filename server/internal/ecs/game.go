package ecs

import (
	"errors"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/storage"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/types"
	gameutils "github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/net"
	"strings"
	"time"
)

type Game struct {
	Done              chan bool
	Engine            *Engine
	GameServerMetrics *GameServerMetrics
}

func NewGame() *Game {
	gameutils.NewLogger()
	return &Game{
		Done:              make(chan bool),
		Engine:            NewEngine(),
		GameServerMetrics: NewGameServerMetrics(),
	}
}

func (g *Game) Start() {
	ticker := time.NewTicker(1000 / gameutils.TickRate * time.Millisecond)
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
	//start := time.Now().UnixMilli()
	g.Engine.UpdateSystems()
	//end := time.Now().UnixMilli()
	//gameutils.Logger.Debug().Msgf("Time taken to process tick: %d ms", end-start)
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
	r := g.Engine.storage.GetSharedResource(gameutils.QuadTreeResource)
	if r != nil {
		g.GameServerMetrics.QuadTree = r.(*storage.QuadTree)
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
