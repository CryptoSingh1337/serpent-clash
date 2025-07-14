package system

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/storage"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/types"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/utils"
)

type PlayerDespawnSystem struct {
	storage         storage.Storage
	despawnQueue    <-chan *types.LeaveEvent
	LeaveEventQueue []*types.LeaveEvent
}

func NewPlayerDespawnSystem(storage storage.Storage, despawnQueue <-chan *types.LeaveEvent) System {
	return &PlayerDespawnSystem{
		storage:      storage,
		despawnQueue: despawnQueue,
	}
}

func (s *PlayerDespawnSystem) Name() string {
	return utils.PlayerDespawnSystemName
}

func (s *PlayerDespawnSystem) Update() {
	if len(s.LeaveEventQueue) > 0 {
		s.LeaveEventQueue = nil
	}
	for {
		select {
		case leaveEvent := <-s.despawnQueue:
			utils.Logger.Info().Msgf("Despawn player with id: %v", leaveEvent.EntityId)
			s.storage.RemoveEntity(leaveEvent.EntityId, utils.PlayerEntity)
			s.LeaveEventQueue = append(s.LeaveEventQueue, leaveEvent)
		default:
			goto escape
		}
	}
escape:
}

func (s *PlayerDespawnSystem) Stop() {

}
