package system

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/component"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/storage"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/types"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/utils"
	"math"
)

type CollisionSystem struct {
	storage storage.Storage
}

func NewCollisionSystem(storage storage.Storage) System {
	return &CollisionSystem{
		storage,
	}
}

func (c *CollisionSystem) Name() string {
	return utils.CollisionSystemName
}

func (c *CollisionSystem) Update() {
	playerEntities := c.storage.GetAllEntitiesByType(utils.PlayerEntity)
	for _, playerId := range playerEntities {
		comp := c.storage.GetComponentByEntityIdAndName(playerId, utils.SnakeComponent)
		if comp == nil {
			continue
		}
		snakeComponent := comp.(*component.Snake)
		head := snakeComponent.Segments[0]
		distanceFromOrigin := utils.EuclideanDistance(0, 0, head.X, head.Y)
		if distanceFromOrigin+utils.SnakeSegmentRadius >= utils.WorldBoundaryRadius {
			comp = c.storage.GetComponentByEntityIdAndName(playerId, utils.NetworkComponent)
			if comp == nil {
				continue
			}
			networkComponent := comp.(*component.Network)
			if err := networkComponent.Connection.Close(); err != nil {
				utils.Logger.Err(err).Msgf("error while closing connection for player")
			}
		}
	}
	q := c.storage.GetSharedResource(utils.QuadTreeResource)
	if q == nil {
		return
	}
	qt := q.(*storage.QuadTree)
	for _, playerId := range playerEntities {
		comp := c.storage.GetComponentByEntityIdAndName(playerId, utils.SnakeComponent)
		if comp == nil {
			continue
		}
		snakeComponent := comp.(*component.Snake)
		head := snakeComponent.Segments[0]
		var points []storage.Point
		qt.QueryBCircle(storage.BCircle{X: head.X, Y: head.Y, R: utils.SnakeSegmentRadius * 2}, &points)
		for _, point := range points {
			switch point.PointType {
			case utils.FoodPointType:
				c.handlePlayerToFoodCollision(playerId, snakeComponent, head, point)
			default:
				c.handlePlayerToPlayerCollision(playerId, head, point)
			}
		}
	}
}

func (c *CollisionSystem) Stop() {
}

func (c *CollisionSystem) handlePlayerToPlayerCollision(playerId types.Id, head utils.Coordinate, point storage.Point) {
	if playerId == point.EntityId {
		return
	}
	distance := utils.EuclideanDistance(point.X, point.Y, head.X, head.Y)
	if distance <= utils.SnakeSegmentRadius*2 {
		switch point.PointType {
		case utils.PlayerHeadPointType:
			utils.Logger.Debug().Msgf("Head to Head collision :: player id: %v, player head: %v, point: %v",
				playerId, head, point)
			compOne := c.storage.GetComponentByEntityIdAndName(playerId, utils.NetworkComponent)
			if compOne == nil {
				return
			}
			playerOneNetworkComponent := compOne.(*component.Network)
			compTwo := c.storage.GetComponentByEntityIdAndName(point.EntityId, utils.NetworkComponent)
			if compTwo == nil {
				return
			}
			playerTwoNetworkComponent := compTwo.(*component.Network)
			if err := playerOneNetworkComponent.Connection.Close(); err != nil {
				utils.Logger.Err(err).Msgf("error while closing connection for player")
			}
			if err := playerTwoNetworkComponent.Connection.Close(); err != nil {
				utils.Logger.Err(err).Msgf("error while closing connection for player")
			}
		case utils.PlayerSegmentPointType:
			utils.Logger.Debug().Msgf("Head to body collision :: player id: %v, player head: %v, point: %v",
				playerId, head, point)
			comp := c.storage.GetComponentByEntityIdAndName(playerId, utils.NetworkComponent)
			if comp == nil {
				return
			}
			networkComponent := comp.(*component.Network)
			if err := networkComponent.Connection.Close(); err != nil {
				utils.Logger.Err(err).Msgf("error while closing connection for player")
			}
		}
	}
}

func (c *CollisionSystem) handlePlayerToFoodCollision(playerId types.Id, snakeComponent *component.Snake,
	head utils.Coordinate, point storage.Point) {
	distance := utils.EuclideanDistance(point.X, point.Y, head.X, head.Y)
	if distance <= utils.FoodConsumeDistance+utils.FoodRadius {
		utils.Logger.Debug().Msgf("Head to food collision :: player id: %v, player head: %v, point: %v",
			playerId, head, point)
		consumeFood(playerId, snakeComponent)
		c.storage.RemoveEntity(point.EntityId, utils.FoodEntity)
	}
}

func consumeFood(playerId types.Id, snakeComponent *component.Snake) {
	snakeComponent.FoodConsumed++
	snakeComponent.GrowthThreshold--
	if snakeComponent.GrowthThreshold < 0 {
		snakeComponent.GrowthThreshold = 0
	}
	if snakeComponent.GrowthThreshold == 0 {
		lastSegmentIdx := len(snakeComponent.Segments) - 1
		lastSegment := snakeComponent.Segments[lastSegmentIdx]
		theta := math.Atan2(lastSegment.Y, lastSegment.X)
		segment := utils.Coordinate{
			X: lastSegment.X - float64(lastSegmentIdx)*math.Cos(theta),
			Y: lastSegment.Y - float64(lastSegmentIdx)*math.Sin(theta),
		}
		snakeComponent.Segments = append(snakeComponent.Segments, segment)
		snakeComponent.GrowthThreshold = utils.DefaultGrowthFactor
		utils.Logger.Debug().Msgf("Added snake segment, player id: %v", playerId)
	}
}
