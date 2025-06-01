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
	MaxFoodEntityId types.Id = 10240
)

type FoodSpawnSystem struct {
	storage       storage.Storage
	foodIdCounter atomic.Uint32
}

func NewFoodSpawnSystem(storage storage.Storage) System {
	return &FoodSpawnSystem{
		storage: storage,
	}
}

func (f *FoodSpawnSystem) Update() {
	r := f.storage.GetSharedResource(utils.QuadTreeResource)
	if r == nil {
		return
	}
	qt := r.(*storage.QuadTree)
	var foodEntities []storage.Point
	qt.QueryByPointType(map[string]bool{utils.FoodPointType: true}, &foodEntities)
	foodCount := len(foodEntities)
	if foodCount < utils.DefaultFoodSpawnThreshold {
		utils.Logger.Info().Msgf("Spawning %d food entities", utils.DefaultFoodSpawnThreshold-foodCount)
		for i := utils.DefaultFoodSpawnThreshold - foodCount; i > 0; i-- {
			entityId := f.newFoodId()
			f.storage.AddEntity(entityId, utils.FoodEntity)
			angle := rand.Float64() * 2 * math.Pi
			radius := utils.SpawnRegionRadius - 200*math.Sqrt(rand.Float64())
			positionComponent := component.NewPositionComponent(radius*math.Cos(angle), radius*math.Sin(angle))
			f.storage.AddComponent(entityId, utils.PositionComponent, &positionComponent)
			expiryComponent := component.NewExpiryComponent(rand.UintN(utils.MaxFoodEntityExpiry))
			f.storage.AddComponent(entityId, utils.ExpiryComponent, &expiryComponent)
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
