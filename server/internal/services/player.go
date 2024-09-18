package services

import (
	"fmt"
	"github.com/CryptoSingh1337/multiplayer-snake-game/server/internal/utils"
	"github.com/google/uuid"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"math/rand"
)

type Player struct {
	Id        string
	Conn      *websocket.Conn
	Positions []utils.Position
	Direction byte
	Color     string
}

func NewPlayer(conn *websocket.Conn) *Player {
	player := &Player{
		Id:        uuid.NewString(),
		Conn:      conn,
		Color:     fmt.Sprintf("hsl(%v, 100%%, 50%%)", rand.Intn(360)),
		Direction: 0,
	}
	player.generateRandomPosition()
	return player
}

func (player *Player) generateRandomPosition() {
	x := 100 + rand.Intn(utils.WorldWidth-100)
	y := 100 + rand.Intn(utils.WorldHeight-100)

	positions := make([]utils.Position, utils.DefaultSnakeLength)
	positions[0] = utils.Position{X: x, Y: y}
	for i := 1; i < len(positions); i++ {
		positions[i].X = x + utils.DefaultGrowthFactor
		positions[i].Y = y + utils.DefaultGrowthFactor
	}
	player.Positions = positions
}
