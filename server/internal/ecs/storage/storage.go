package storage

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/component"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/types"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/utils"
)

type Storage interface {
	AddEntity(entityId types.Id, entityType string)
	RemoveEntity(entityId types.Id, entityType string)
	GetAllEntities() map[string][]types.Id
	GetAllEntitiesByType(componentType string) []types.Id
	GetAllComponentByName(componentName string) any
	GetComponentByEntityIdAndNames(entityId types.Id, componentNames ...string) map[string]any
	GetComponentByEntityIdAndName(entityId types.Id, componentName string) (any, bool)
	AddComponent(entityId types.Id, entityType string, component any)
	DeleteComponent(entityId types.Id, entityType string)
	UpdateComponent(entityId types.Id, entityType string, component any)
}

type SimpleStorage struct {
	entities             []types.Id
	entityGroup          map[string][]types.Id
	inputComponents      *Pool[component.Input]
	networkComponents    *Pool[component.Network]
	playerInfoComponents *Pool[component.PlayerInfo]
	snakeComponents      *Pool[component.Snake]
}

func NewSimpleStorage() Storage {
	return &SimpleStorage{
		entities:             make([]types.Id, 0),
		entityGroup:          make(map[string][]types.Id),
		inputComponents:      NewPool[component.Input](),
		networkComponents:    NewPool[component.Network](),
		playerInfoComponents: NewPool[component.PlayerInfo](),
		snakeComponents:      NewPool[component.Snake](),
	}
}

func (s *SimpleStorage) AddEntity(entityId types.Id, entityType string) {
	s.entities = append(s.entities, entityId)
	_, exists := s.entityGroup[entityType]
	if !exists {
		s.entityGroup[entityType] = make([]types.Id, 5)
	}
	s.entityGroup[entityType] = append(s.entityGroup[entityType], entityId)
}

func (s *SimpleStorage) RemoveEntity(entityId types.Id, entityType string) {
	entityIdIdx := -1
	for idx, id := range s.entities {
		if id == entityId {
			entityIdIdx = idx
		}
	}
	if entityIdIdx != -1 {
		utils.RemoveFromSlice(s.entities, entityIdIdx)
		entityIds := s.entityGroup[entityType]
		entityIdIdx = -1
		for idx, id := range entityIds {
			if id == entityId {
				entityIdIdx = idx
			}
		}
		if entityIdIdx != -1 {
			utils.RemoveFromSlice(entityIds, entityIdIdx)
		}
	}
}

func (s *SimpleStorage) GetAllEntities() map[string][]types.Id {
	return s.entityGroup
}

func (s *SimpleStorage) GetAllEntitiesByType(t string) []types.Id {
	if entities, exists := s.entityGroup[t]; exists {
		return entities
	}
	return []types.Id{}
}

func (s *SimpleStorage) GetAllComponentByName(componentName string) any {
	switch componentName {
	case "input":
		return s.inputComponents.GetAll()
	case "network":
		return s.networkComponents.GetAll()
	case "playerInfo":
		return s.playerInfoComponents.GetAll()
	case "snake":
		return s.snakeComponents.GetAll()
	}
	return nil
}

func (s *SimpleStorage) GetComponentByEntityIdAndNames(entityId types.Id, componentNames ...string) map[string]any {
	components := make(map[string]any)
	for _, name := range componentNames {
		var c any
		exists := false
		switch name {
		case "input":
			c, exists = s.inputComponents.Get(entityId)
		case "network":
			c, exists = s.networkComponents.Get(entityId)
		case "playerInfo":
			c, exists = s.playerInfoComponents.Get(entityId)
		case "snake":
			c, exists = s.snakeComponents.Get(entityId)
		}
		if exists {
			components[name] = c
		}
	}
	return components
}

func (s *SimpleStorage) GetComponentByEntityIdAndName(entityId types.Id, componentName string) (any, bool) {
	var c any
	exists := false
	switch componentName {
	case "input":
		c, exists = s.inputComponents.Get(entityId)
	case "network":
		c, exists = s.networkComponents.Get(entityId)
	case "playerInfo":
		c, exists = s.playerInfoComponents.Get(entityId)
	case "snake":
		c, exists = s.snakeComponents.Get(entityId)
	}
	if exists {
		return c, true
	}
	return nil, false
}

func (s *SimpleStorage) AddComponent(entityId types.Id, entityType string, com any) {
	switch entityType {
	case "input":
		c := com.(component.Input)
		s.inputComponents.Add(entityId, c)
	case "network":
		c := com.(component.Network)
		s.networkComponents.Add(entityId, c)
	case "playerInfo":
		c := com.(component.PlayerInfo)
		s.playerInfoComponents.Add(entityId, c)
	case "snake":
		c := com.(component.Snake)
		s.snakeComponents.Add(entityId, c)
	}
}

func (s *SimpleStorage) DeleteComponent(entityId types.Id, entityType string) {
	//TODO implement me
	panic("implement me")
}

func (s *SimpleStorage) UpdateComponent(entityId types.Id, entityType string, component any) {
	//TODO implement me
	panic("implement me")
}
