package system

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/storage"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/types"
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
