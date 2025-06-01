package system

import "github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/storage"

type FoodDespawnSystem struct {
	storage storage.Storage
}

func NewFoodDespawnSystem(storage storage.Storage) System {
	return &FoodDespawnSystem{
		storage,
	}
}

func (f FoodDespawnSystem) Update() {
}

func (f FoodDespawnSystem) Stop() {
}
