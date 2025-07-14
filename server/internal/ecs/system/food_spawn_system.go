package system

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/component"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/storage"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/types"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/utils"
	"math"
	"math/rand/v2"
	"sync/atomic"
)

const (
	MinFoodEntityId types.Id = 2048
	MaxFoodEntityId types.Id = 999999999
)

type FoodSpawnSystem struct {
	storage             storage.Storage
	foodIdCounter       atomic.Uint32
	FoodSpawnEventQueue []*types.FoodSpawnEvent
}

func NewFoodSpawnSystem(storage storage.Storage) System {
	return &FoodSpawnSystem{
		storage: storage,
	}
}

func (f *FoodSpawnSystem) Name() string {
	return utils.FoodSpawnSystemName
}

func (f *FoodSpawnSystem) Update() {
	if len(f.FoodSpawnEventQueue) > 0 {
		f.FoodSpawnEventQueue = nil
	}
	r := f.storage.GetSharedResource(utils.QuadTreeResource)
	if r == nil {
		return
	}
	qt := r.(*storage.QuadTree)
	var foodEntities []storage.Point
	qt.QueryByPointType(map[string]bool{utils.FoodPointType: true}, &foodEntities)
	foodCount := len(foodEntities)
	if foodCount < utils.FoodSpawnThreshold {
		utils.Logger.Info().Msgf("Spawning %d food entities", utils.FoodSpawnThreshold-foodCount)
		for i := utils.FoodSpawnThreshold - foodCount; i > 0; i-- {
			entityId := f.newFoodId()
			f.storage.AddEntity(entityId, utils.FoodEntity)
			angle := rand.Float64() * 2 * math.Pi
			radius := 100 + float64(rand.Uint64N(utils.WorldBoundaryRadius-100))
			positionComponent := component.NewPositionComponent(radius*math.Cos(angle), radius*math.Sin(angle))
			f.storage.AddComponent(entityId, utils.PositionComponent, &positionComponent)
			expiryComponent := component.NewExpiryComponent(uint32(utils.MinFoodEntityExpiry +
				rand.UintN(uint(utils.MaxFoodEntityExpiry-utils.MinFoodEntityExpiry))))
			f.storage.AddComponent(entityId, utils.ExpiryComponent, &expiryComponent)
			f.FoodSpawnEventQueue = append(f.FoodSpawnEventQueue, &types.FoodSpawnEvent{
				EntityId: entityId,
				Coordinate: utils.Coordinate{
					X: positionComponent.X,
					Y: positionComponent.Y,
				},
			})
			f.storage.AddSharedResource(utils.FoodSpawnEventQueue, f.FoodSpawnEventQueue)
		}
	}
}

func (f *FoodSpawnSystem) Stop() {
}

func (f *FoodSpawnSystem) newFoodId() types.Id {
	for {
		val := f.foodIdCounter.Load()
		if f.foodIdCounter.CompareAndSwap(val, val+1) {
			return (types.Id(val) % (MaxFoodEntityId - MinFoodEntityId)) + MinFoodEntityId
		}
	}
}
