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
	GetComponentByEntityIdAndName(entityId types.Id, componentName string) any
	AddComponent(entityId types.Id, componentName string, component any)
	ReplaceComponent(entityId types.Id, componentName string, component any)
	DeleteComponent(entityId types.Id, componentName string)
	PrintState()
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
		entities:             make([]types.Id, 0, 10),
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
		s.entityGroup[entityType] = make([]types.Id, 0, 5)
	}
	s.entityGroup[entityType] = append(s.entityGroup[entityType], entityId)
}

func (s *SimpleStorage) RemoveEntity(entityId types.Id, entityType string) {
	for idx, id := range s.entities {
		if id == entityId {
			s.entities = utils.RemoveFromSlice(s.entities, idx)
			break
		}
	}

	if entityIds, exists := s.entityGroup[entityType]; exists {
		for idx, id := range entityIds {
			if id == entityId {
				s.entityGroup[entityType] = utils.RemoveFromSlice(entityIds, idx)
				break
			}
		}
	}
	s.inputComponents.Remove(entityId)
	s.networkComponents.Remove(entityId)
	s.playerInfoComponents.Remove(entityId)
	s.snakeComponents.Remove(entityId)
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
	case utils.InputComponent:
		return s.inputComponents.GetAll()
	case utils.NetworkComponent:
		return s.networkComponents.GetAll()
	case utils.PlayerInfoComponent:
		return s.playerInfoComponents.GetAll()
	case utils.SnakeComponent:
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
		case utils.InputComponent:
			c, exists = s.inputComponents.Get(entityId)
		case utils.NetworkComponent:
			c, exists = s.networkComponents.Get(entityId)
		case utils.PlayerInfoComponent:
			c, exists = s.playerInfoComponents.Get(entityId)
		case utils.SnakeComponent:
			c, exists = s.snakeComponents.Get(entityId)
		}
		if exists {
			components[name] = c
		}
	}
	return components
}

func (s *SimpleStorage) GetComponentByEntityIdAndName(entityId types.Id, componentName string) any {
	var c any
	exists := false
	switch componentName {
	case utils.InputComponent:
		c, exists = s.inputComponents.Get(entityId)
	case utils.NetworkComponent:
		c, exists = s.networkComponents.Get(entityId)
	case utils.PlayerInfoComponent:
		c, exists = s.playerInfoComponents.Get(entityId)
	case utils.SnakeComponent:
		c, exists = s.snakeComponents.Get(entityId)
	}
	if exists {
		return c
	}
	return nil
}

func (s *SimpleStorage) AddComponent(entityId types.Id, componentName string, com any) {
	switch componentName {
	case utils.InputComponent:
		c := com.(*component.Input)
		s.inputComponents.Add(entityId, c)
	case utils.NetworkComponent:
		c := com.(*component.Network)
		s.networkComponents.Add(entityId, c)
	case utils.PlayerInfoComponent:
		c := com.(*component.PlayerInfo)
		s.playerInfoComponents.Add(entityId, c)
	case utils.SnakeComponent:
		c := com.(*component.Snake)
		s.snakeComponents.Add(entityId, c)
	}
}

func (s *SimpleStorage) ReplaceComponent(entityId types.Id, componentName string, com any) {
	switch componentName {
	case utils.InputComponent:
		c := com.(*component.Input)
		s.inputComponents.Replace(entityId, c)
	case utils.NetworkComponent:
		c := com.(*component.Network)
		s.networkComponents.Replace(entityId, c)
	case utils.PlayerInfoComponent:
		c := com.(*component.PlayerInfo)
		s.playerInfoComponents.Replace(entityId, c)
	case utils.SnakeComponent:
		c := com.(*component.Snake)
		s.snakeComponents.Replace(entityId, c)
	}
}

func (s *SimpleStorage) DeleteComponent(entityId types.Id, componentName string) {
	switch componentName {
	case utils.InputComponent:
		s.inputComponents.Remove(entityId)
	case utils.NetworkComponent:
		s.networkComponents.Remove(entityId)
	case utils.PlayerInfoComponent:
		s.playerInfoComponents.Remove(entityId)
	case utils.SnakeComponent:
		s.snakeComponents.Remove(entityId)
	}
}

func (s *SimpleStorage) PrintState() {
	utils.Logger.Debug().Msgf("Entities: %v, EntityGroup: %v", s.entities, s.entityGroup)
	s.inputComponents.PrintState(utils.InputComponent)
	s.networkComponents.PrintState(utils.NetworkComponent)
	s.playerInfoComponents.PrintState(utils.PlayerInfoComponent)
	s.snakeComponents.PrintState(utils.SnakeComponent)
}
