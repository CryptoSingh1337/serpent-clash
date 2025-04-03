package ecs

import (
	"errors"
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
	Done             chan bool
	Engine           *Engine
	GameServerMetric *GameServerMetric
}

func NewGame() *Game {
	return &Game{
		Done:             make(chan bool),
		Engine:           NewEngine(),
		GameServerMetric: NewGameServerMetric(),
	}
}

func (g *Game) Start() {
	ticker := time.NewTicker(1000 / utils.TickRate * time.Millisecond)
	engineTicker := make(chan time.Time)
	metricTicker := make(chan time.Time)
	g.Engine.Start()
	go func() {
		for {
			select {
			case <-g.Done:
				ticker.Stop()
				return
			case t := <-ticker.C:
				engineTicker <- t
				metricTicker <- t
			}
		}
	}()
	go func() {
		for {
			select {
			case _ = <-engineTicker:
				g.processTick()
			}
		}
	}()
	go func() {
		for {
			select {
			case _ = <-metricTicker:
				g.processMetric()
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

func (g *Game) processMetric() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	g.GameServerMetric.MemoryUsage = m.Sys / 1024 / 1024
	g.GameServerMetric.PlayerCount = uint8(len(g.Engine.playerIdToEntityId))
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
	playerId := uuid.New().String()
	upgrader.OnMessage(func(c *websocket.Conn, messageType websocket.MessageType, data []byte) {
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
