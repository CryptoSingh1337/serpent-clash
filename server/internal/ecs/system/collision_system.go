package system

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/component"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/storage"
	gameutils "github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/utils"
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
	playerEntities := c.storage.GetAllEntitiesByType(gameutils.PlayerEntity)
	for _, playerId := range playerEntities {
		comp := c.storage.GetComponentByEntityIdAndName(playerId, gameutils.SnakeComponent)
		if comp == nil {
			continue
		}
		snakeComponent := comp.(*component.Snake)
		head := snakeComponent.Segments[0]
		distanceFromOrigin := gameutils.EuclideanDistance(0, 0, head.X, head.Y)
		if distanceFromOrigin+gameutils.SnakeSegmentDiameter/2 >= gameutils.WorldBoundaryRadius {
			comp = c.storage.GetComponentByEntityIdAndName(playerId, gameutils.NetworkComponent)
			if comp == nil {
				continue
			}
			networkComponent := comp.(*component.Network)
			if err := networkComponent.Connection.Close(); err != nil {
				gameutils.Logger.Err(err).Msgf("error while closing connection for player")
			}
		}
	}
	q := c.storage.GetSharedResource(gameutils.QuadTreeResource)
	if q == nil {
		return
	}
	qt := q.(*storage.QuadTree)
	for _, playerId := range playerEntities {
		comp := c.storage.GetComponentByEntityIdAndName(playerId, gameutils.SnakeComponent)
		if comp == nil {
			continue
		}
		snakeComponent := comp.(*component.Snake)
		head := snakeComponent.Segments[0]
		var points []storage.Point
		qt.QueryBCircle(storage.BCircle{X: head.X, Y: head.Y, R: gameutils.SnakeSegmentDiameter}, &points)
		for _, point := range points {
			if playerId == point.EntityId {
				continue
			}
			distance := gameutils.EuclideanDistance(point.X, point.Y, head.X, head.Y)
			if distance <= gameutils.SnakeSegmentDiameter {
				gameutils.Logger.Debug().Msgf("Collision :: playerId: %v, player head: %v, point: %v",
					playerId, head, point)
				switch point.PointType {
				case "head":
					gameutils.Logger.Debug().Msgf("Head to head collision")
					compOne := c.storage.GetComponentByEntityIdAndName(playerId, gameutils.NetworkComponent)
					if compOne == nil {
						continue
					}
					playerOneNetworkComponent := compOne.(*component.Network)
					compTwo := c.storage.GetComponentByEntityIdAndName(playerId, gameutils.NetworkComponent)
					if compTwo == nil {
						continue
					}
					playerTwoNetworkComponent := compTwo.(*component.Network)
					if err := playerOneNetworkComponent.Connection.Close(); err != nil {
						gameutils.Logger.Err(err).Msgf("error while closing connection for player")
					}
					if err := playerTwoNetworkComponent.Connection.Close(); err != nil {
						gameutils.Logger.Err(err).Msgf("error while closing connection for player")
					}
				case "segment":
					gameutils.Logger.Debug().Msgf("Head to body collision")
					_comp := c.storage.GetComponentByEntityIdAndName(playerId, gameutils.NetworkComponent)
					if _comp == nil {
						continue
					}
					networkComponent := _comp.(*component.Network)
					if err := networkComponent.Connection.Close(); err != nil {
						gameutils.Logger.Err(err).Msgf("error while closing connection for player")
					}
				}
				break
			}
		}
	}
}

func (c *CollisionSystem) Stop() {

}
