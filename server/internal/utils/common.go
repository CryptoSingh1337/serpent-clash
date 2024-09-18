package utils

import (
	"encoding/json"
	"fmt"
)

const (
	PlayerMoving        = iota
	TickRate            = 1
	PlayerSpeed         = 10
	DefaultSnakeLength  = 1
	DefaultGrowthFactor = 2
	MaxPlayerAllowed    = 10
	WorldFactor         = 200
	WorldHeight         = 3 * WorldFactor
	WorldWidth          = 4 * WorldFactor
)

const (
	Left  = iota
	Right = iota
	Up    = iota
	Down  = iota
)

const (
	PingMessage       = "ping"
	PongMessage       = "pong"
	GameStateMessage  = "game_state"
	InitializeMessage = "initialize"
)

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Payload struct {
	Type string          `json:"type"`
	Body json.RawMessage `json:"body"`
}

type PlayerState struct {
	Positions []Position `json:"positions"`
	Direction byte       `json:"direction"`
}

type GameState struct {
	PlayerStates map[string]PlayerState `json:"playerStates"`
}

func (p Payload) String() string {
	return fmt.Sprintf("{type=%v, body=%v}", p.Type, string(p.Body))
}
