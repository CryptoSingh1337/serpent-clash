package system

import (
	"math"

	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/component"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/storage"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/utils"
)

type MovementSystem struct {
	storage storage.Storage
}

func NewMovementSystem(storage storage.Storage) System {
	return &MovementSystem{
		storage,
	}
}

func (m *MovementSystem) Name() string {
	return utils.MovementSystemName
}

func (m *MovementSystem) Update() {
	playerEntityIds := m.storage.GetAllEntitiesByType(utils.PlayerEntity)
	for _, playerEntityId := range playerEntityIds {
		c := m.storage.GetComponentByEntityIdAndName(playerEntityId, utils.InputComponent)
		if c == nil {
			continue
		}
		inputComponent := c.(*component.Input)
		c = m.storage.GetComponentByEntityIdAndName(playerEntityId, utils.SnakeComponent)
		if c == nil {
			continue
		}
		snakeComponent := c.(*component.Snake)
		if len(snakeComponent.Segments) == 0 {
			continue
		}
		previousCoordinate := inputComponent.PrevCoordinates
		mouseCoordinate := inputComponent.Coordinates
		head := snakeComponent.Segments[0]
		angle := snakeComponent.Angle
		if previousCoordinate.X != mouseCoordinate.X || previousCoordinate.Y != mouseCoordinate.Y {
			headToInputDistance := utils.EuclideanDistance(head.X, head.Y, mouseCoordinate.X, mouseCoordinate.Y)
			utils.Logger.Info().Msgf("Head to input distance: %f", headToInputDistance)
			if headToInputDistance > utils.SnakeSegmentRadius*1.25 {
				targetAngle := math.Atan2(mouseCoordinate.Y-head.Y, mouseCoordinate.X-head.X)
				angle = utils.LerpAngle(angle, targetAngle, utils.MaxTurnRate)
			}
		}

		// Move the head towards the mouse coordinate
		speed := float64(utils.DefaultPlayerSpeed)
		if inputComponent.Boost {
			if snakeComponent.Stamina > 0 {
				snakeComponent.Stamina--
			} else {
				if len(snakeComponent.Segments) > 5 {
					snakeComponent.Segments = snakeComponent.Segments[:len(snakeComponent.Segments)-1]
					snakeComponent.Stamina = utils.DefaultSnakeStamina
				} else {
					inputComponent.Boost = false
				}
			}
		}
		if inputComponent.Boost {
			speed += utils.DefaultSpeedBoost
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
