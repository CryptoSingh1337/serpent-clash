package system

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/component"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/storage"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/utils"
	"math"
)

type MovementSystem struct {
	storage storage.Storage
}

func NewMovementSystem(storage storage.Storage) *MovementSystem {
	return &MovementSystem{
		storage,
	}
}

func (m *MovementSystem) Update() {
	playerEntityIds := m.storage.GetAllEntitiesByType("player")
	for _, playerEntityId := range playerEntityIds {
		c := m.storage.GetComponentByEntityIdAndName(playerEntityId, "input")
		if c == nil {
			continue
		}
		inputComponent := c.(*component.Input)
		c = m.storage.GetComponentByEntityIdAndName(playerEntityId, "snake")
		if c == nil {
			continue
		}
		snakeComponent := c.(*component.Snake)
		mouseCoordinate := inputComponent.Coordinates
		head := snakeComponent.Segments[0]
		angle := snakeComponent.Angle
		targetAngle := math.Atan2(mouseCoordinate.Y-head.Y, mouseCoordinate.X-head.X)
		angle = utils.LerpAngle(angle, targetAngle, utils.MaxTurnRate)

		// Move the head towards the mouse coordinate
		speed := float64(utils.PlayerSpeed)
		if inputComponent.Boost {
			speed += utils.PlayerBoostSpeed
		}
		head.X += math.Cos(angle) * speed
		head.Y += math.Sin(angle) * speed

		// Update the head position and angle
		snakeComponent.Segments[0] = head
		snakeComponent.Angle = angle

		// Move the rest of the snake to follow the head
		for i := 1; i < len(snakeComponent.Segments); i++ {
			prevSegment := snakeComponent.Segments[i-1]
			currentSegment := snakeComponent.Segments[i]

			angleToPrev := math.Atan2(prevSegment.Y-currentSegment.Y, prevSegment.X-currentSegment.X)

			// Keep a fixed distance between segments
			currentSegment.X = prevSegment.X - math.Cos(angleToPrev)*utils.SnakeSegmentDistance
			currentSegment.Y = prevSegment.Y - math.Sin(angleToPrev)*utils.SnakeSegmentDistance
			snakeComponent.Segments[i] = currentSegment
		}
	}
}

func (m *MovementSystem) Stop() {

}
