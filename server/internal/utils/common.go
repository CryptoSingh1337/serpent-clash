package utils

import (
	"encoding/json"
	"fmt"
)

// Common constants
const (
	TickRate             = 60
	PlayerSpeed          = 5
	PlayerBoostSpeed     = 3
	MaxTurnRate          = 0.03
	DefaultSnakeLength   = 10
	DefaultGrowthFactor  = 2
	SnakeSegmentDistance = 15
	SnakeSegmentDiameter = 50
	MaxPlayerAllowed     = 10
	WorldBoundaryRadius  = 2750
)

// ChatMessage types
const (
	HelloMessage     = "hello"
	PingMessage      = "ping"
	PongMessage      = "pong"
	GameStateMessage = "game_state"
	Movement         = "movement"
	SpeedBoost       = "boost"
	Kill             = "kill"
)

// Entity types
const (
	PlayerEntity = "player"
	FoodEntity   = "food"
)

// Component names
const (
	InputComponent      = "input"
	NetworkComponent    = "network"
	PlayerInfoComponent = "playerInfo"
	PositionComponent   = "position"
	SnakeComponent      = "snake"
)

type Coordinate struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Payload struct {
	Type string          `json:"type"`
	Body json.RawMessage `json:"body"`
}

func (p Payload) String() string {
	return fmt.Sprintf("{type=%v, body=%v}", p.Type, string(p.Body))
}

type MouseEvent struct {
	Seq        uint64     `json:"seq"`
	Coordinate Coordinate `json:"coordinate"`
}

type SpeedBoostEvent struct {
	Seq     uint64 `json:"seq"`
	Enabled bool   `json:"enabled"`
}

type PingEvent struct {
	Timestamp uint32 `json:"timestamp"`
}

type DeathEvent struct {
	PlayerId string `json:"playerId"`
}

type PlayerState struct {
	Color    string       `json:"color"`
	Segments []Coordinate `json:"positions"`
	Seq      uint64       `json:"seq"`
}

type GameState struct {
	PlayerStates map[string]PlayerState `json:"playerStates"`
}

type ChatMessage struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}
