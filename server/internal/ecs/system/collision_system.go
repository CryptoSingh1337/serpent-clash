package system

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/component"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/storage"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/utils"
)

type CollisionSystem struct {
	storage storage.Storage
}

func NewCollisionSystem(storage storage.Storage) *CollisionSystem {
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
	qt := storage.NewQuadTree(storage.BBox{X: 0, Y: 0, W: utils.WorldWeight, H: utils.WorldHeight}, 15)
	for _, playerId := range playerEntities {
		comp := c.storage.GetComponentByEntityIdAndName(playerId, utils.SnakeComponent)
		if comp == nil {
			continue
		}
		snakeComponent := comp.(*component.Snake)
		head := snakeComponent.Segments[0]
		qt.Insert(storage.Point{X: head.X, Y: head.Y, EntityId: playerId, PointType: "head"})
		for i := 1; i < len(snakeComponent.Segments); i++ {
			segment := snakeComponent.Segments[i]
			qt.Insert(storage.Point{X: segment.X, Y: segment.Y, EntityId: playerId, PointType: "segment"})
		}
	}
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
			if playerId == point.EntityId {
				continue
			}
			distance := utils.EuclideanDistance(point.X, point.Y, head.X, head.Y)
			if distance <= utils.SnakeSegmentDiameter {
				utils.Logger.Debug().Msgf("Collision :: playerId: %v, player head: %v, point: %v", playerId, head, point)
				switch point.PointType {
				case "head":
					utils.Logger.Debug().Msgf("Head to head collision")
					compOne := c.storage.GetComponentByEntityIdAndName(playerId, utils.NetworkComponent)
					if compOne == nil {
						continue
					}
					playerOneNetworkComponent := compOne.(*component.Network)
					compTwo := c.storage.GetComponentByEntityIdAndName(playerId, utils.NetworkComponent)
					if compTwo == nil {
						continue
					}
					playerTwoNetworkComponent := compTwo.(*component.Network)
					if err := playerOneNetworkComponent.Connection.Close(); err != nil {
						utils.Logger.Err(err).Msgf("error while closing connection for player")
					}
					if err := playerTwoNetworkComponent.Connection.Close(); err != nil {
						utils.Logger.Err(err).Msgf("error while closing connection for player")
					}
				case "segment":
					utils.Logger.Debug().Msgf("Head to body collision")
					_comp := c.storage.GetComponentByEntityIdAndName(playerId, utils.NetworkComponent)
					if _comp == nil {
						continue
					}
					networkComponent := _comp.(*component.Network)
					if err := networkComponent.Connection.Close(); err != nil {
						utils.Logger.Err(err).Msgf("error while closing connection for player")
					}
				}
				break
			}
		}
	}
}

func (c *CollisionSystem) Stop() {

}
