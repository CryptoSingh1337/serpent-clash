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
	systems            []system.System
	playerIdToEntityId map[string]types.Id
	JoinQueue          chan *types.JoinEvent
	LeaveQueue         chan string
	SpawnQueue         chan types.Id
}

func NewEngine() *Engine {
	simpleStorage := storage.NewSimpleStorage()
	engine := &Engine{
		storage:            simpleStorage,
		systems:            make([]system.System, 0),
		playerIdToEntityId: make(map[string]types.Id),
		JoinQueue:          make(chan *types.JoinEvent, utils.MaxPlayerAllowed),
		LeaveQueue:         make(chan string, utils.MaxPlayerAllowed),
		SpawnQueue:         make(chan types.Id, utils.MaxPlayerAllowed),
	}
	quadTreeSystem := system.NewQuadTreeSystem(simpleStorage)
	movementSystem := system.NewMovementSystem(simpleStorage)
	playerSpawnSystem := system.NewSpawnSystem(simpleStorage, engine.SpawnQueue)
	collisionSystem := system.NewCollisionSystem(simpleStorage)
	foodSpawnSystem := system.NewFoodSpawnSystem(simpleStorage)
	foodDespawnSystem := system.NewFoodDespawnSystem(simpleStorage)
	networkSystem := system.NewNetworkSystem(simpleStorage)
	engine.systems = append(engine.systems, quadTreeSystem, movementSystem, playerSpawnSystem, collisionSystem,
		foodSpawnSystem, foodDespawnSystem, networkSystem)
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
	e.SpawnQueue <- entityId
	e.storage.AddComponent(entityId, utils.InputComponent, &inputComponent)
	e.storage.AddComponent(entityId, utils.NetworkComponent, &networkComponent)
	e.storage.AddComponent(entityId, utils.PlayerInfoComponent, &playerInfoComponent)
	e.storage.AddComponent(entityId, utils.SnakeComponent, &snakeComponent)
	c := e.storage.GetComponentByEntityIdAndName(entityId, utils.InputComponent)
	if c == nil {
		utils.Logger.Error().Msgf("Input component is nil")
	}
	c = e.storage.GetComponentByEntityIdAndName(entityId, utils.NetworkComponent)
	if c == nil {
		utils.Logger.Error().Msgf("Network component is nil")
	}
	c = e.storage.GetComponentByEntityIdAndName(entityId, utils.PlayerInfoComponent)
	if c == nil {
		utils.Logger.Error().Msgf("Player info component is nil")
	}
	c = e.storage.GetComponentByEntityIdAndName(entityId, utils.SnakeComponent)
	if c == nil {
		utils.Logger.Error().Msgf("Snake component is nil")
	}
	if networkComponent.Connection == nil {
		utils.Logger.Error().Msgf("connection is nil")
	}
	body := fmt.Sprintf(`{"id":%q}`, joinEvent.PlayerId)
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
	networkComponent := e.storage.GetComponentByEntityIdAndName(entityId, utils.NetworkComponent).(*component.Network)
	networkComponent.Connected = false
	e.storage.RemoveEntity(entityId, utils.PlayerEntity)
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
	for _, s := range e.systems {
		s.Update()
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
