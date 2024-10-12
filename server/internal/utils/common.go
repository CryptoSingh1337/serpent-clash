package utils

import (
	"encoding/json"
	"fmt"
	"math"
)

// Common constants
const (
	TickRate             = 60
	PlayerSpeed          = 4
	MaxTurnRate          = 0.03
	DefaultSnakeLength   = 10
	DefaultGrowthFactor  = 2
	SnakeSegmentDistance = 15
	SnakeSegmentDiameter = 50
	MaxPlayerAllowed     = 10
	WorldBoundaryMinX    = -3000
	WorldBoundaryMaxX    = 3000
	WorldBoundaryMinY    = -3000
	WorldBoundaryMaxY    = 3000
)

// Message types
const (
	HelloMessage     = "hello"
	PingMessage      = "ping"
	PongMessage      = "pong"
	GameStateMessage = "game_state"
	Movement         = "movement"
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

type PingEvent struct {
	Timestamp uint32 `json:"timestamp"`
}

type PlayerState struct {
	Color    string       `json:"color"`
	Segments []Coordinate `json:"positions"`
	Seq      uint64       `json:"seq"`
}

type GameState struct {
	PlayerStates map[string]PlayerState `json:"playerStates"`
}

func LerpAngle(a, b, t float64) float64 {
	diff := b - a
	// Handle wrapping from -π to π
	for diff < -math.Pi {
		diff += 2 * math.Pi
	}
	for diff > math.Pi {
		diff -= 2 * math.Pi
	}
	return a + diff*math.Min(t, 1.0)
}
