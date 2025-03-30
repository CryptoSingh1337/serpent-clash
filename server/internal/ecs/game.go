package ecs

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/utils"
	"github.com/labstack/echo/v4"
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
	//username := c.QueryParam("username")
	//username = strings.TrimSpace(username)
	//if username == "" {
	//	return errors.New("invalid username")
	//}
	//w := c.Response()
	//r := c.Request()
	//upgrader := websocket.NewUpgrader()
	//upgrader.EnableCompression(false)
	//conn, err := upgrader.Upgrade(w, r, nil)
	//if err != nil {
	//	return err
	//}
	//player := services.NewPlayer(&username, conn, w)
	//upgrader.OnMessage(func(c *websocket.Conn, messageType websocket.MessageType, data []byte) {
	//	g.ProcessEvent(player, messageType, data)
	//})
	//g.JoinQueue <- player
	//conn.OnClose(func(c *websocket.Conn, err error) {
	//	g.LeaveQueue <- player
	//})
	return nil
}
