package services

import (
	"encoding/json"
	"github.com/CryptoSingh1337/multiplayer-snake-game/server/internal/utils"
	"github.com/google/uuid"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"math/rand"
)

type Player struct {
	Id        string
	Conn      *websocket.Conn
	Positions []utils.Position
}

type Payload struct {
	Type string          `json:"type"`
	Body json.RawMessage `json:"body"`
}

func NewPlayer(conn *websocket.Conn) *Player {
	player := &Player{
		Id:   uuid.NewString(),
		Conn: conn,
	}
	player.generateRandomPosition()
	return player
}

func (player *Player) generateRandomPosition() {
	x := 100 + rand.Intn(worldWidth-100)
	y := 100 + rand.Intn(worldHeight-100)

	positions := make([]utils.Position, defaultSnakeLength)
	positions[0] = utils.Position{X: x, Y: y}
	for i := 1; i < len(positions); i++ {
		positions[i].X = x + defaultGrowthFactor
		positions[i].Y = y + defaultGrowthFactor
	}
	player.Positions = positions
}
