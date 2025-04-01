package ecs

import (
	"errors"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"strings"
	"time"
)

type Game struct {
	World *World
	Done  chan bool
}

func NewGame() *Game {
	return &Game{
		World: NewWorld(),
		Done:  make(chan bool),
	}
}

func (g *Game) Start() {
	ticker := time.NewTicker(1000 / utils.TickRate * time.Millisecond)
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
}

func (g *Game) Stop() {
	g.World.Stop()
	g.Done <- true
	close(g.Done)
}

func (g *Game) processTick() {
	start := time.Now().UnixMilli()
	g.World.Update(1000 / utils.TickRate * time.Millisecond)
	end := time.Now().UnixMilli()
	utils.Logger.Debug().Msgf("Time taken to process tick: %d ms", end-start)
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
		g.World.ProcessEvent(playerId, messageType, data)
	})
	g.World.JoinQueue <- &utils.JoinEvent{
		PlayerId:   playerId,
		Connection: conn,
	}
	conn.OnClose(func(c *websocket.Conn, err error) {
		g.World.LeaveQueue <- playerId
	})
	return nil
}
