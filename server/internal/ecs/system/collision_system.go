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
	// Get all players
	playerEntities := c.storage.GetAllEntitiesByType(utils.PlayerEntity)
	for _, playerId := range playerEntities {
		comp := c.storage.GetComponentByEntityIdAndName(playerId, utils.SnakeComponent)
		if comp == nil {
			continue
		}
		snakeComponent := comp.(*component.Snake)
		head := snakeComponent.Segments[0]
		distanceFromOrigin := utils.EuclideanDistance(0, 0, head.X, head.Y)
		if distanceFromOrigin >= utils.WorldBoundaryRadius {
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
}

func (c *CollisionSystem) Stop() {

}
