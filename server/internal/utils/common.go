package utils

import (
	"encoding/json"
	"fmt"
	"math"
)

const (
	PlayerMoving         = iota
	TickRate             = 60
	PlayerSpeed          = 3
	MaxTurnRate          = 0.05
	DefaultSnakeLength   = 10
	DefaultGrowthFactor  = 2
	SnakeSegmentDistance = 15
	SnakeSegmentRadius   = 100
	MaxPlayerAllowed     = 10
	WorldFactor          = 200
	WorldHeight          = 3 * WorldFactor
	WorldWidth           = 4 * WorldFactor
)

const (
	Left  = iota
	Right = iota
	Up    = iota
	Down  = iota
)

const (
	HelloMessage     = "hello"
	PingMessage      = "ping"
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
	Coordinate Coordinate `json:"coordinate"`
}

type PlayerState struct {
	Color    string       `json:"color"`
	Segments []Coordinate `json:"positions"`
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
