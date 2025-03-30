package ecs

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/storage"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/system"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/types"
)

type Engine struct {
	storage storage.Storage
	systems []system.System
}

func NewEngine() *Engine {
	simpleStorage := storage.NewSimpleStorage()
	engine := &Engine{
		storage: simpleStorage,
		systems: make([]system.System, 3),
	}
	var networkSystem system.System = system.NewNetworkSystem(simpleStorage)
	var movementSystem system.System = system.NewMovementSystem(simpleStorage)
	var collisionSystem system.System = system.NewCollisionSystem(simpleStorage)
	engine.systems = append(engine.systems, movementSystem, collisionSystem, networkSystem)
	return engine
}

func (e *Engine) AddEntity(entityId types.Id, entityType string) {
	e.storage.AddEntity(entityId, entityType)
}

func (e *Engine) RemoveEntity(entityId types.Id) {
	e.storage.RemoveEntity(entityId)
}

func (e *Engine) Update() {
	for _, s := range e.systems {
		s.Update()
	}
}

func (e *Engine) Stop() {
	for _, s := range e.systems {
		s.Stop()
	}
}
