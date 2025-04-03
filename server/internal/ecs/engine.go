package ecs

import (
	"errors"
	"fmt"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/component"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/storage"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/system"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/types"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/utils"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"sync/atomic"
	"time"
)

const (
	firstEntity types.Id = 1
	MaxEntity   types.Id = 1024
)

type Engine struct {
	idCounter          atomic.Uint32
	minId, maxId       types.Id
	storage            storage.Storage
	systems            []system.System
	playerIdToEntityId map[string]types.Id
	JoinQueue          chan *types.JoinEvent
	LeaveQueue         chan string
}

func NewEngine() *Engine {
	simpleStorage := storage.NewSimpleStorage()
	engine := &Engine{
		minId:              firstEntity + 1,
		maxId:              MaxEntity,
		storage:            simpleStorage,
		systems:            make([]system.System, 0),
		playerIdToEntityId: make(map[string]types.Id),
		JoinQueue:          make(chan *types.JoinEvent, utils.MaxPlayerAllowed),
		LeaveQueue:         make(chan string, utils.MaxPlayerAllowed),
	}
	var networkSystem system.System = system.NewNetworkSystem(simpleStorage)
	var movementSystem system.System = system.NewMovementSystem(simpleStorage)
	var collisionSystem system.System = system.NewCollisionSystem(simpleStorage)
	engine.systems = append(engine.systems, movementSystem, collisionSystem, networkSystem)
	return engine
}

func (e *Engine) AddPlayer(joinEvent *types.JoinEvent) error {
	utils.Logger.Info().Msgf("Inside Engine.AddPlayer :: joinEvent: %v", joinEvent)
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
	payload, err := utils.ToJsonB(utils.Payload{Type: utils.HelloMessage, Body: []byte(body)})
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
	e.storage.DeleteComponent(entityId, utils.InputComponent)
	e.storage.DeleteComponent(entityId, utils.PlayerInfoComponent)
	e.storage.DeleteComponent(entityId, utils.SnakeComponent)
	e.storage.DeleteComponent(entityId, utils.NetworkComponent)
	return nil
}

func (e *Engine) UpdateSystems() {
	// Process all players in JoinQueue
	for {
		select {
		case joinEvent := <-e.JoinQueue:
			entityId := e.newId()
			joinEvent.EntityId = entityId
			if err := e.AddPlayer(joinEvent); err != nil {
				utils.Logger.LogError().Msgf("Error adding player %s: %v", joinEvent.PlayerId, err)
			}
		default:
			// Exit the loop when JoinQueue is empty
			goto ProcessLeaveQueue
		}
	}
ProcessLeaveQueue:
	// Process all players in LeaveQueue
	for {
		select {
		case playerId := <-e.LeaveQueue:
			if err := e.RemovePlayer(playerId); err != nil {
				utils.Logger.LogError().Msgf("Error removing player %s: %v", playerId, err)
			}
		default:
			// Exit the loop when JoinQueue is empty
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
		payload, err := utils.FromJsonB[utils.Payload](data)
		if err != nil {
			return
		}
		switch payload.Type {
		case utils.Movement:
			mouseEvent, err := utils.FromJsonB[utils.MouseEvent](payload.Body)
			if err != nil {
				return
			}
			inputComponent := e.storage.GetComponentByEntityIdAndName(entityId, utils.InputComponent).(*component.Input)
			inputComponent.Coordinates.X = mouseEvent.Coordinate.X
			inputComponent.Coordinates.Y = mouseEvent.Coordinate.Y
			networkComponent := e.storage.GetComponentByEntityIdAndName(entityId, utils.NetworkComponent).(*component.Network)
			networkComponent.MessageSequence = mouseEvent.Seq
		case utils.SpeedBoost:
			speedBoostEvent, err := utils.FromJsonB[utils.SpeedBoostEvent](payload.Body)
			if err != nil {
				return
			}
			inputComponent := e.storage.GetComponentByEntityIdAndName(entityId, utils.InputComponent).(*component.Input)
			inputComponent.Boost = speedBoostEvent.Enabled
		case utils.PingMessage:
			pingEvent, err := utils.FromJsonB[utils.PingEvent](payload.Body)
			if err != nil {
				return
			}
			networkComponent := e.storage.GetComponentByEntityIdAndName(entityId, utils.NetworkComponent).(*component.Network)
			networkComponent.RequestInitiateTimestamp = pingEvent.Timestamp
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
