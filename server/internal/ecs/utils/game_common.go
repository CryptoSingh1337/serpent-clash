package utils

import (
	"encoding/json"
	"fmt"
)

// Common constants
const (
	TickRate                = 60
	PingCooldown            = 30
	PlayerSpeed             = 6
	PlayerBoostSpeed        = 3
	MaxTurnRate             = 0.03
	DefaultSnakeLength      = 10
	DefaultGrowthFactor     = 10
	SnakeSegmentDistance    = 20
	SnakeSegmentDiameter    = 50
	FoodConsumeDistance     = 30
	FoodSpawnThreshold      = 50
	MinFoodEntityExpiry     = 4000
	MaxFoodEntityExpiry     = 10000
	MaxPlayerAllowed        = 10
	QuadTreeSegmentCapacity = 50
	WorldBoundaryRadius     = 2850
	WorldHeight             = 3000
	WorldWeight             = 3000
	SpawnRegionRadius       = WorldBoundaryRadius * 0.175
)

// Message types
const (
	HelloMessageType       = "hello"
	PingMessageType        = "ping"
	PongMessageType        = "pong"
	MovementMessageType    = "movement"
	PlayerStateMessageType = "player_state"
	FoodStateMessageType   = "food_state"
)

// Entity types
const (
	PlayerEntity = "player"
	FoodEntity   = "food"
)

// Component names
const (
	InputComponent      = "input"
	ExpiryComponent     = "expiry"
	NetworkComponent    = "network"
	PlayerInfoComponent = "playerInfo"
	PositionComponent   = "position"
	SnakeComponent      = "snake"
)

// QuadTreeResource Shared resources names
const (
	QuadTreeResource      = "quad_tree"
	SpawnRegions          = "spawn_regions"
	FoodSpawnEventQueue   = "food_spawn_event_queue"
	FoodDespawnEventQueue = "food_despawn_event_queue"
)

const (
	PlayerHeadPointType    = "head"
	PlayerSegmentPointType = "segment"
	FoodPointType          = "food"
)

// System names
const (
	CollisionSystemName     = "CollisionSystem"
	FoodDespawnSystemName   = "FoodDespawnSystem"
	FoodSpawnSystemName     = "FoodSpawnSystem"
	MovementSystemName      = "MovementSystem"
	NetworkSystemName       = "NetworkSystem"
	PlayerDespawnSystemName = "PlayerDespawnSystem"
	PlayerSpawnSystemName   = "PlayerSpawnSystem"
	QuadTreeSystemName      = "QuadTreeSystem"
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
