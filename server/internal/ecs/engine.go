package ecs

import (
	"errors"
	"fmt"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/component"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/storage"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/system"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/types"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/utils"
	"github.com/gorilla/websocket"
	"sync/atomic"
	"time"
)

const (
	MinPlayerEntityId types.Id = 1
	MaxPlayerEntityId types.Id = 1024
)

type Engine struct {
	playerIdCounter    atomic.Uint32
	storage            storage.Storage
	systems            map[string]system.System
	systemUpdateOrder  []string
	playerIdToEntityId map[string]types.Id
	JoinQueue          chan *types.JoinEvent
	LeaveQueue         chan string
	SpawnQueue         chan *types.JoinEvent
	DespawnQueue       chan *types.LeaveEvent
}

func NewEngine() *Engine {
	simpleStorage := storage.NewSimpleStorage()
	engine := &Engine{
		storage:            simpleStorage,
		systems:            make(map[string]system.System),
		playerIdToEntityId: make(map[string]types.Id),
		JoinQueue:          make(chan *types.JoinEvent, utils.MaxPlayerAllowed),
		LeaveQueue:         make(chan string, utils.MaxPlayerAllowed),
		SpawnQueue:         make(chan *types.JoinEvent, utils.MaxPlayerAllowed),
		DespawnQueue:       make(chan *types.LeaveEvent, utils.MaxPlayerAllowed),
	}
	quadTreeSystem := system.NewQuadTreeSystem(simpleStorage)
	movementSystem := system.NewMovementSystem(simpleStorage)
	playerSpawnSystem := system.NewSpawnSystem(simpleStorage, engine.SpawnQueue)
	playerDespawnSystem := system.NewPlayerDespawnSystem(simpleStorage, engine.DespawnQueue)
	collisionSystem := system.NewCollisionSystem(simpleStorage)
	foodSpawnSystem := system.NewFoodSpawnSystem(simpleStorage)
	foodDespawnSystem := system.NewFoodDespawnSystem(simpleStorage)
	networkSystem := system.NewNetworkSystem(simpleStorage)
	engine.systems[quadTreeSystem.Name()] = quadTreeSystem
	engine.systems[movementSystem.Name()] = movementSystem
	engine.systems[playerSpawnSystem.Name()] = playerSpawnSystem
	engine.systems[playerDespawnSystem.Name()] = playerDespawnSystem
	engine.systems[collisionSystem.Name()] = collisionSystem
	engine.systems[foodSpawnSystem.Name()] = foodSpawnSystem
	engine.systems[foodDespawnSystem.Name()] = foodDespawnSystem
	engine.systems[networkSystem.Name()] = networkSystem
	engine.systemUpdateOrder = []string{
		quadTreeSystem.Name(),
		movementSystem.Name(),
		playerSpawnSystem.Name(),
		playerDespawnSystem.Name(),
		collisionSystem.Name(),
		foodSpawnSystem.Name(),
		foodDespawnSystem.Name(),
		networkSystem.Name(),
	}
	return engine
}

func (e *Engine) Start() {
}

func (e *Engine) AddPlayer(joinEvent *types.JoinEvent) error {
	utils.Logger.Info().Msgf("Inside Engine.AddPlayer :: joinEvent: %v", joinEvent)
	// TODO: add max player validation
	if joinEvent.PlayerId == "" {
		return errors.New("invalid player id")
	}
	if joinEvent.Connection == nil {
		return errors.New("invalid websocket connection")
	}
	_, exists := e.playerIdToEntityId[joinEvent.PlayerId]
	if exists {
		return errors.New("player already exists")
	}
	entityId := joinEvent.EntityId
	e.playerIdToEntityId[joinEvent.PlayerId] = entityId
	e.storage.AddEntity(entityId, utils.PlayerEntity)
	inputComponent := component.NewInputComponent()
	networkComponent := component.NewNetworkComponent(joinEvent.Connection)
	playerInfoComponent := component.NewPlayerInfoComponent(joinEvent.PlayerId, joinEvent.Username)
	snakeComponent := component.NewSnakeComponent()
	e.SpawnQueue <- joinEvent
	e.storage.AddComponent(entityId, utils.InputComponent, &inputComponent)
	e.storage.AddComponent(entityId, utils.NetworkComponent, &networkComponent)
	e.storage.AddComponent(entityId, utils.PlayerInfoComponent, &playerInfoComponent)
	e.storage.AddComponent(entityId, utils.SnakeComponent, &snakeComponent)
	c := e.storage.GetComponentByEntityIdAndName(entityId, utils.InputComponent)
	if c == nil {
		utils.Logger.Error().Msgf("Input component is nil")
		return errors.New("input component is nil")
	}
	c = e.storage.GetComponentByEntityIdAndName(entityId, utils.NetworkComponent)
	if c == nil {
		utils.Logger.Error().Msgf("Network component is nil")
		return errors.New("network component is nil")
	}
	c = e.storage.GetComponentByEntityIdAndName(entityId, utils.PlayerInfoComponent)
	if c == nil {
		utils.Logger.Error().Msgf("Player info component is nil")
		return errors.New("player info component is nil")
	}
	c = e.storage.GetComponentByEntityIdAndName(entityId, utils.SnakeComponent)
	if c == nil {
		utils.Logger.Error().Msgf("Snake component is nil")
		return errors.New("snake component is nil")
	}
	if networkComponent.Connection == nil {
		utils.Logger.Error().Msgf("connection is nil")
		return errors.New("connection is nil")
	}
	c = e.storage.GetAllComponentByName(utils.PlayerInfoComponent)
	if c == nil {
		utils.Logger.Error().Msgf("player info components are nil")
		return errors.New("player info components are nil")
	}
	playerInfoComponents := c.([]*component.PlayerInfo)
	playerDetails := map[string]any{
		"id":      joinEvent.PlayerId,
		"players": make([]map[string]any, 0, len(playerInfoComponents)),
	}
	for _, playerInfo := range playerInfoComponents {
		if playerInfo.ID == joinEvent.PlayerId {
			continue // Skip the joining player
		}
		playerDetails["players"] = append(playerDetails["players"].([]map[string]any), map[string]any{
			playerInfo.ID: playerInfo.Username,
		})
	}
	p, err := utils.ToJsonS(playerDetails["players"])
	if err != nil {
		return err
	}
	body := fmt.Sprintf(`{"id": "%s", "players": %s}`, joinEvent.PlayerId, p)
	payload, err := utils.ToJsonB(utils.Payload{Type: utils.HelloMessageType, Body: []byte(body)})
	if err != nil {
		return err
	}
	return networkComponent.Connection.WriteMessage(websocket.TextMessage, payload)
}

func (e *Engine) RemovePlayer(playerId string) error {
	utils.Logger.Info().Msgf("Inside Engine.RemovePlayer :: playerId: %v", playerId)
	entityId, exists := e.playerIdToEntityId[playerId]
	if !exists {
		return errors.New("player does not exists")
	}
	e.DespawnQueue <- &types.LeaveEvent{
		EntityId: entityId,
		PlayerId: playerId,
	}
	delete(e.playerIdToEntityId, playerId)
	return nil
}

func (e *Engine) UpdateSystems() {
	for {
		select {
		case joinEvent := <-e.JoinQueue:
			entityId := e.newPlayerId()
			joinEvent.EntityId = entityId
			if err := e.AddPlayer(joinEvent); err != nil {
				utils.Logger.Error().Msgf("Error adding player %s: %v", joinEvent.PlayerId, err)
			}
		default:
			goto ProcessLeaveQueue
		}
	}
ProcessLeaveQueue:
	for {
		select {
		case playerId := <-e.LeaveQueue:
			if err := e.RemovePlayer(playerId); err != nil {
				utils.Logger.Error().Msgf("Error removing player %s: %v", playerId, err)
			}
		default:
			goto ProcessSystemUpdates
		}
	}
ProcessSystemUpdates:
	for _, s := range e.systemUpdateOrder {
		e.systems[s].Update()
	}
}

func (e *Engine) Stop() {
	for _, s := range e.systems {
		s.Stop()
	}
}

func (e *Engine) ProcessEvent(playerId string, messageType int, data []byte) {
	entityId, exists := e.playerIdToEntityId[playerId]
	if !exists {
		return
	}
	switch messageType {
	case websocket.TextMessage:
		payload, err := utils.FromJsonB[utils.Payload](data)
		if err != nil {
			return
		}
		switch payload.Type {
		case utils.MovementMessageType:
			inputEvent, err := utils.FromJsonB[utils.InputEvent](payload.Body)
			if err != nil {
				return
			}
			inputComponent := e.storage.GetComponentByEntityIdAndName(entityId, utils.InputComponent).(*component.Input)
			networkComponent := e.storage.GetComponentByEntityIdAndName(entityId, utils.NetworkComponent).(*component.Network)
			networkComponent.MessageSequence = inputEvent.Seq
			inputComponent.Coordinates.X = inputEvent.Coordinate.X
			inputComponent.Coordinates.Y = inputEvent.Coordinate.Y
			inputComponent.Boost = inputEvent.Boost
		case utils.PingMessageType:
			pingEvent, err := utils.FromJsonB[types.PingEvent](payload.Body)
			pingEvent.PlayerId = playerId
			if err != nil {
				return
			}
			networkComponent := e.storage.GetComponentByEntityIdAndName(entityId, utils.NetworkComponent).(*component.Network)
			networkComponent.RequestInitiateTimestamp = pingEvent.RequestInitiateTimestamp
			networkComponent.RequestAckTimestamp = uint64(time.Now().UnixMilli())
		}
	}
}

func (e *Engine) newPlayerId() types.Id {
	for {
		val := e.playerIdCounter.Load()
		if e.playerIdCounter.CompareAndSwap(val, val+1) {
			return (types.Id(val) % (MaxPlayerEntityId - MinPlayerEntityId)) + MinPlayerEntityId
		}
	}
}
