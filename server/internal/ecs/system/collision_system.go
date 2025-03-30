package system

import "github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/storage"

type CollisionSystem struct {
	storage storage.Storage
}

func NewCollisionSystem(storage storage.Storage) *CollisionSystem {
	return &CollisionSystem{
		storage,
	}
}

func (c *CollisionSystem) Update() {
	//TODO implement me
	panic("implement me")
}

func (c *CollisionSystem) Stop() {

}
