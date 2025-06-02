package system

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/component"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/storage"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/types"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/utils"
)

type CollisionSystem struct {
	storage storage.Storage
}

func NewCollisionSystem(storage storage.Storage) System {
	return &CollisionSystem{
		storage,
	}
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
		if distanceFromOrigin+utils.SnakeSegmentDiameter/2 >= utils.WorldBoundaryRadius {
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
		qt.QueryBCircle(storage.BCircle{X: head.X, Y: head.Y, R: utils.SnakeSegmentDiameter}, &points)
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
	if distance <= utils.SnakeSegmentDiameter {
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
	if distance <= utils.FoodConsumeDistance {
		utils.Logger.Debug().Msgf("Head to food collision :: player id: %v, player head: %v, point: %v",
			playerId, head, point)
		snakeComponent.FoodConsumed++
	}
}
