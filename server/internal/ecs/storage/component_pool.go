package storage

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/types"
)

type Pool[T types.Component] struct {
	pool          []*T
	entityToIndex map[types.Id]int
	indexToEntity map[int]types.Id
	nextIndex     int
}

func NewPool[T types.Component]() *Pool[T] {
	return &Pool[T]{
		pool:          make([]*T, 0, 10),
		entityToIndex: make(map[types.Id]int),
		indexToEntity: make(map[int]types.Id),
		nextIndex:     0,
	}
}

func (p *Pool[T]) Add(entityId types.Id, component *T) {
	if idx, exists := p.entityToIndex[entityId]; exists {
		p.pool[idx] = component
		return
	}
	p.pool = append(p.pool, component)
	p.entityToIndex[entityId] = p.nextIndex
	p.indexToEntity[p.nextIndex] = entityId
	p.nextIndex++
}

func (p *Pool[T]) Get(entityId types.Id) (*T, bool) {
	idx, exists := p.entityToIndex[entityId]
	if !exists {
		var c T
		return &c, false
	}
	return p.pool[idx], true
}

func (p *Pool[T]) GetAll() []*T {
	return p.pool
}

func (p *Pool[T]) Replace(entityId types.Id, component *T) {
	idx, exists := p.entityToIndex[entityId]
	if !exists {
		return
	}
	p.pool[idx] = component
}

func (p *Pool[T]) Remove(entityId types.Id) {
	idx, exists := p.entityToIndex[entityId]
	if !exists {
		return
	}
	delete(p.entityToIndex, entityId)
	delete(p.indexToEntity, idx)
	if idx < len(p.pool)-1 {
		p.pool[idx] = p.pool[len(p.pool)-1]
		lastEntityId := p.indexToEntity[len(p.pool)-1]
		p.entityToIndex[lastEntityId] = idx
		p.indexToEntity[idx] = lastEntityId
	}
	p.pool = p.pool[:len(p.pool)-1]
	p.nextIndex--
}
