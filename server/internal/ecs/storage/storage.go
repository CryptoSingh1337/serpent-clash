package storage

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/component"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/types"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/utils"
)

type Storage interface {
	AddEntity(entityId types.Id, entityType string)
	RemoveEntity(entityId types.Id, entityType string)
	AddSharedResource(resourceName string, resource any)
	GetSharedResource(resourceName string) any
	DeleteSharedResource(resourceName string)
	GetAllEntitiesByType(componentType string) []types.Id
	GetAllComponentByName(componentName string) any
	GetComponentByEntityIdAndName(entityId types.Id, componentName string) any
	AddComponent(entityId types.Id, componentName string, component any)
	DeleteComponent(entityId types.Id, componentName string)
	LogState()
}

type SimpleStorage struct {
	entityGroup          map[string][]types.Id
	sharedResources      map[string]any
	inputComponents      *Pool[component.Input]
	networkComponents    *Pool[component.Network]
	playerInfoComponents *Pool[component.PlayerInfo]
	snakeComponents      *Pool[component.Snake]
}

func NewSimpleStorage() Storage {
	return &SimpleStorage{
		entityGroup:          make(map[string][]types.Id),
		sharedResources:      make(map[string]any),
		inputComponents:      NewPool[component.Input](),
		networkComponents:    NewPool[component.Network](),
		playerInfoComponents: NewPool[component.PlayerInfo](),
		snakeComponents:      NewPool[component.Snake](),
	}
}

func (s *SimpleStorage) AddEntity(entityId types.Id, entityType string) {
	switch entityType {
	case utils.PlayerEntity:
		_, exists := s.entityGroup[entityType]
		if !exists {
			s.entityGroup[entityType] = make([]types.Id, 0, 5)
		}
		s.entityGroup[entityType] = append(s.entityGroup[entityType], entityId)
	default:
		utils.Logger.Error().Msgf("%s: invalid entity type", entityType)
	}
}

func (s *SimpleStorage) RemoveEntity(entityId types.Id, entityType string) {
	switch entityType {
	case utils.PlayerEntity:
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
	default:
		utils.Logger.Error().Msgf("%s: invalid entity type", entityType)
	}
}

func (s *SimpleStorage) AddSharedResource(resourceName string, resource any) {
	s.sharedResources[resourceName] = resource
}

func (s *SimpleStorage) GetSharedResource(resourceName string) any {
	if resource, exists := s.sharedResources[resourceName]; exists {
		return resource
	}
	return nil
}

func (s *SimpleStorage) DeleteSharedResource(resourceName string) {
	if _, exists := s.sharedResources[resourceName]; exists {
		delete(s.sharedResources, resourceName)
	}
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

func (s *SimpleStorage) LogState() {
	utils.Logger.Debug().Msgf("EntityGroup: %v", s.entityGroup)
	utils.Logger.Debug().Msgf(s.inputComponents.String(utils.InputComponent))
	utils.Logger.Debug().Msgf(s.networkComponents.String(utils.NetworkComponent))
	utils.Logger.Debug().Msgf(s.playerInfoComponents.String(utils.PlayerInfoComponent))
	utils.Logger.Debug().Msgf(s.snakeComponents.String(utils.SnakeComponent))
}
