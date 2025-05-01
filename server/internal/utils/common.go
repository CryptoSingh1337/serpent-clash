package utils

import (
	"encoding/json"
	"fmt"
)

// Common constants
const (
	TickRate                  = 60
	PingCooldown              = 30
	PlayerSpeed               = 6
	PlayerBoostSpeed          = 3
	MaxTurnRate               = 0.03
	DefaultSnakeLength        = 10
	DefaultGrowthFactor       = 2
	SnakeSegmentDistance      = 20
	SnakeSegmentDiameter      = 50
	DefaultFoodSpawnThreshold = 10
	MaxPlayerAllowed          = 10
	WorldBoundaryRadius       = 2850
	WorldHeight               = 3000
	WorldWeight               = 3000
	SpawnRegionRadius         = WorldBoundaryRadius * 0.175
)

// Message types
const (
	HelloMessageType     = "hello"
	PingMessageType      = "ping"
	PongMessageType      = "pong"
	GameStateMessageType = "game_state"
	MovementMessageType  = "movement"
	KillMessageType      = "kill"
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

// QuadTreeResource Shared resources names
const (
	QuadTreeResource = "quad_tree"
	SpawnRegions     = "spawn_regions"
)

const (
	PlayerHeadPointType    = "head"
	PlayerSegmentPointType = "segment"
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

type InputEvent struct {
	Seq        uint64     `json:"seq"`
	Coordinate Coordinate `json:"coordinate"`
	Boost      bool       `json:"boost"`
}

type PongMessage struct {
	RequestInitiateTimestamp  uint64 `json:"reqInit"`
	RequestAckTimestamp       uint64 `json:"reqAck"`
	ResponseInitiateTimestamp uint64 `json:"resInit"`
}

type PlayerStateMessage struct {
	Color    string       `json:"color"`
	Segments []Coordinate `json:"positions"`
	Seq      uint64       `json:"seq"`
}

type GameStateMessage struct {
	PlayerStates map[string]PlayerStateMessage `json:"playerStates"`
}
