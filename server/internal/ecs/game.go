package ecs

import (
	"errors"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/storage"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/types"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"runtime"
	"strings"
	"time"
)

type Game struct {
	Done              chan bool
	Engine            *Engine
	GameServerMetrics *GameServerMetrics
}

func NewGame() *Game {
	return &Game{
		Done:              make(chan bool),
		Engine:            NewEngine(),
		GameServerMetrics: NewGameServerMetrics(),
	}
}

func (g *Game) Start() {
	ticker := time.NewTicker(1000 / utils.TickRate * time.Millisecond)
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
	//utils.Logger.Debug().Msgf("Time taken to process tick: %d ms", end-start)
}

func (g *Game) processMetrics() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	g.GameServerMetrics.MemoryUsage = m.Sys / 1024 / 1024
	g.GameServerMetrics.PlayerCount = uint8(len(g.Engine.playerIdToEntityId))
	r := g.Engine.storage.GetSharedResource(utils.QuadTreeResource)
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
