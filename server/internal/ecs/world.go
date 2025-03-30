package ecs

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/types"
	"sync/atomic"
	"time"
)

const (
	firstEntity types.Id = 1
	MaxEntity   types.Id = 1024
)

type World struct {
	idCounter    atomic.Uint32
	minId, maxId types.Id
	engine       *Engine
}

func NewWorld() *World {
	return &World{
		minId:  firstEntity + 1,
		maxId:  MaxEntity,
		engine: NewEngine(),
	}
}

func (w *World) NewId() types.Id {
	for {
		val := w.idCounter.Load()
		if w.idCounter.CompareAndSwap(val, val+1) {
			return (types.Id(val) % (w.maxId - w.minId)) + w.minId
		}
	}
}

func (w *World) Update(dt time.Duration) {
	w.engine.Update()
}

func (w *World) Stop() {
	w.engine.Stop()
}
