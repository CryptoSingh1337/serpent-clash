package storage

import (
	"fmt"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/types"
)

type Pool[T types.Component] struct {
	pool          []*T
	entityToIndex map[types.Id]int
	indexToEntity map[int]types.Id
}

func NewPool[T types.Component]() *Pool[T] {
	return &Pool[T]{
		pool:          make([]*T, 0, 10),
		entityToIndex: make(map[types.Id]int),
		indexToEntity: make(map[int]types.Id),
	}
}

func (p *Pool[T]) Add(entityId types.Id, component *T) {
	if idx, exists := p.entityToIndex[entityId]; exists {
		p.pool[idx] = component
		return
	}
	p.entityToIndex[entityId] = len(p.pool)
	p.indexToEntity[len(p.pool)] = entityId
	p.pool = append(p.pool, component)
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
	lastIdx := len(p.pool) - 1
	lastEntityId := p.indexToEntity[lastIdx]
	if idx != lastIdx {
		p.pool[idx] = p.pool[lastIdx]
		p.entityToIndex[lastEntityId] = idx
		p.indexToEntity[idx] = lastEntityId
	}
	delete(p.entityToIndex, entityId)
	delete(p.indexToEntity, lastIdx)
	p.pool = p.pool[:lastIdx]
}

func (p *Pool[T]) String(name string) string {
	return fmt.Sprintf("%v: Entity to index: %v, Index to Entity: %v", name, p.entityToIndex, p.indexToEntity)
}
