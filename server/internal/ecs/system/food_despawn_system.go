package system

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/component"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/storage"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/types"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/utils"
)

type FoodDespawnSystem struct {
	storage          storage.Storage
	FoodDespawnQueue []*types.FoodDespawnEvent
}

func NewFoodDespawnSystem(storage storage.Storage) System {
	return &FoodDespawnSystem{
		storage:          storage,
		FoodDespawnQueue: make([]*types.FoodDespawnEvent, 0, utils.FoodSpawnThreshold),
	}
}

func (f *FoodDespawnSystem) Name() string {
	return utils.FoodDespawnSystemName
}

func (f *FoodDespawnSystem) Update() {
	if len(f.FoodDespawnQueue) > 0 {
		f.FoodDespawnQueue = f.FoodDespawnQueue[:0]
	}
	foodEntities := f.storage.GetAllEntitiesByType(utils.FoodEntity)
	var expiredFoodEntityIds []types.Id
	for _, entityId := range foodEntities {
		c := f.storage.GetComponentByEntityIdAndName(entityId, utils.ExpiryComponent)
		if c == nil {
			continue
		}
		expiryComponent := c.(*component.Expiry)
		expiryComponent.TicksRemaining -= 1
		if expiryComponent.TicksRemaining <= 0 {
			expiredFoodEntityIds = append(expiredFoodEntityIds, entityId)
		}
	}
	if len(expiredFoodEntityIds) == 0 {
		return
	}
	utils.Logger.Info().Msgf("De-spawning %d food entities", len(expiredFoodEntityIds))
	for _, entityId := range expiredFoodEntityIds {
		f.storage.RemoveEntity(entityId, utils.FoodEntity)
		f.FoodDespawnQueue = append(f.FoodDespawnQueue, &types.FoodDespawnEvent{
			EntityId: entityId,
		})
	}
	f.storage.AddSharedResource(utils.FoodDespawnEventQueue, f.FoodDespawnQueue)
}

func (f *FoodDespawnSystem) Stop() {
}
