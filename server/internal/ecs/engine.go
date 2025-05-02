package ecs

import (
	"errors"
	"fmt"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/component"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/storage"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/system"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/types"
	gameutils "github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/utils"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"sync/atomic"
	"time"
)

const (
	firstEntity types.Id = 1
	MaxEntity   types.Id = 10240
)

type Engine struct {
	idCounter          atomic.Uint32
	minId, maxId       types.Id
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
		minId:              firstEntity,
		maxId:              MaxEntity,
		storage:            simpleStorage,
		systems:            make([]system.System, 0),
		playerIdToEntityId: make(map[string]types.Id),
		JoinQueue:          make(chan *types.JoinEvent, gameutils.MaxPlayerAllowed),
		LeaveQueue:         make(chan string, gameutils.MaxPlayerAllowed),
		SpawnQueue:         make(chan types.Id, gameutils.MaxPlayerAllowed),
	}
	var movementSystem system.System = system.NewMovementSystem(simpleStorage)
	var spawnSystem system.System = system.NewSpawnSystem(simpleStorage, engine.SpawnQueue, engine.newId)
	var collisionSystem system.System = system.NewCollisionSystem(simpleStorage)
	var networkSystem system.System = system.NewNetworkSystem(simpleStorage)
	engine.systems = append(engine.systems, movementSystem, spawnSystem, collisionSystem, networkSystem)
	return engine
}

func (e *Engine) Start() {
}

func (e *Engine) AddPlayer(joinEvent *types.JoinEvent) error {
	gameutils.Logger.Info().Msgf("Inside Engine.AddPlayer :: joinEvent: %v", joinEvent)
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
	e.storage.AddEntity(entityId, gameutils.PlayerEntity)
	inputComponent := component.NewInputComponent()
	networkComponent := component.NewNetworkComponent(joinEvent.Connection)
	playerInfoComponent := component.NewPlayerInfoComponent(joinEvent.PlayerId, joinEvent.Username)
	snakeComponent := component.NewSnakeComponent()
	e.SpawnQueue <- entityId
	e.storage.AddComponent(entityId, gameutils.InputComponent, &inputComponent)
	e.storage.AddComponent(entityId, gameutils.NetworkComponent, &networkComponent)
	e.storage.AddComponent(entityId, gameutils.PlayerInfoComponent, &playerInfoComponent)
	e.storage.AddComponent(entityId, gameutils.SnakeComponent, &snakeComponent)
	c := e.storage.GetComponentByEntityIdAndName(entityId, gameutils.InputComponent)
	if c == nil {
		gameutils.Logger.Error().Msgf("Input component is nil")
	}
	c = e.storage.GetComponentByEntityIdAndName(entityId, gameutils.NetworkComponent)
	if c == nil {
		gameutils.Logger.Error().Msgf("Network component is nil")
	}
	c = e.storage.GetComponentByEntityIdAndName(entityId, gameutils.PlayerInfoComponent)
	if c == nil {
		gameutils.Logger.Error().Msgf("Player info component is nil")
	}
	c = e.storage.GetComponentByEntityIdAndName(entityId, gameutils.SnakeComponent)
	if c == nil {
		gameutils.Logger.Error().Msgf("Snake component is nil")
	}
	if networkComponent.Connection == nil {
		gameutils.Logger.Error().Msgf("connection is nil")
	}
	body := fmt.Sprintf(`{"id":%q}`, joinEvent.PlayerId)
	payload, err := gameutils.ToJsonB(gameutils.Payload{Type: gameutils.HelloMessageType, Body: []byte(body)})
	if err != nil {
		return err
	}
	return networkComponent.Connection.WriteMessage(websocket.TextMessage, payload)
}

func (e *Engine) RemovePlayer(playerId string) error {
	gameutils.Logger.Info().Msgf("Inside Engine.RemovePlayer :: playerId: %v", playerId)
	entityId, exists := e.playerIdToEntityId[playerId]
	if !exists {
		return errors.New("player does not exists")
	}
	networkComponent := e.storage.GetComponentByEntityIdAndName(entityId, gameutils.NetworkComponent).(*component.Network)
	networkComponent.Connected = false
	e.storage.RemoveEntity(entityId, gameutils.PlayerEntity)
	delete(e.playerIdToEntityId, playerId)
	return nil
}

func (e *Engine) UpdateSystems() {
	for {
		select {
		case joinEvent := <-e.JoinQueue:
			entityId := e.newId()
			joinEvent.EntityId = entityId
			if err := e.AddPlayer(joinEvent); err != nil {
				gameutils.Logger.Error().Msgf("Error adding player %s: %v", joinEvent.PlayerId, err)
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
				gameutils.Logger.Error().Msgf("Error removing player %s: %v", playerId, err)
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

func (e *Engine) ProcessEvent(playerId string, messageType websocket.MessageType, data []byte) {
	entityId, exists := e.playerIdToEntityId[playerId]
	if !exists {
		return
	}
	switch messageType {
	case websocket.TextMessage:
		payload, err := gameutils.FromJsonB[gameutils.Payload](data)
		if err != nil {
			return
		}
		switch payload.Type {
		case gameutils.MovementMessageType:
			inputEvent, err := gameutils.FromJsonB[gameutils.InputEvent](payload.Body)
			if err != nil {
				return
			}
			inputComponent := e.storage.GetComponentByEntityIdAndName(entityId, gameutils.InputComponent).(*component.Input)
			networkComponent := e.storage.GetComponentByEntityIdAndName(entityId, gameutils.NetworkComponent).(*component.Network)
			networkComponent.MessageSequence = inputEvent.Seq
			inputComponent.Coordinates.X = inputEvent.Coordinate.X
			inputComponent.Coordinates.Y = inputEvent.Coordinate.Y
			inputComponent.Boost = inputEvent.Boost
		case gameutils.PingMessageType:
			pingEvent, err := gameutils.FromJsonB[types.PingEvent](payload.Body)
			pingEvent.PlayerId = playerId
			if err != nil {
				return
			}
			networkComponent := e.storage.GetComponentByEntityIdAndName(entityId, gameutils.NetworkComponent).(*component.Network)
			networkComponent.RequestInitiateTimestamp = pingEvent.RequestInitiateTimestamp
			networkComponent.RequestAckTimestamp = uint64(time.Now().UnixMilli())
		}
	}
}

func (e *Engine) newId() types.Id {
	for {
		val := e.idCounter.Load()
		if e.idCounter.CompareAndSwap(val, val+1) {
			return (types.Id(val) % (e.maxId - e.minId)) + e.minId
		}
	}
}
