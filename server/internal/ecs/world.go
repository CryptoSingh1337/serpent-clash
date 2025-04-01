package ecs

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/component"
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

type World struct {
	idCounter          atomic.Uint32
	minId, maxId       types.Id
	engine             *Engine
	playerIdToEntityId map[string]types.Id
	JoinQueue          chan *utils.JoinEvent
	LeaveQueue         chan string
	pingQueue          chan string
}

func NewWorld() *World {
	return &World{
		minId:              firstEntity + 1,
		maxId:              MaxEntity,
		engine:             NewEngine(),
		playerIdToEntityId: make(map[string]types.Id),
	}
}

func (w *World) ProcessEvent(playerId string, messageType websocket.MessageType, data []byte) {
	entityId, exists := w.playerIdToEntityId[playerId]
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
			inputComponent := w.engine.GetComponent(entityId, utils.InputComponent).(component.Input)
			inputComponent.Coordinates.X = mouseEvent.Coordinate.X
			inputComponent.Coordinates.Y = mouseEvent.Coordinate.Y
			networkComponent := w.engine.GetComponent(entityId, utils.NetworkComponent).(component.Network)
			networkComponent.MessageSequence = mouseEvent.Seq
		case utils.SpeedBoost:
			speedBoostEvent, err := utils.FromJsonB[utils.SpeedBoostEvent](payload.Body)
			if err != nil {
				return
			}
			inputComponent := w.engine.GetComponent(entityId, utils.InputComponent).(component.Input)
			inputComponent.Boost = speedBoostEvent.Enabled
		case utils.PingMessage:
			pingEvent, err := utils.FromJsonB[utils.PingEvent](payload.Body)
			if err != nil {
				return
			}
			networkComponent := w.engine.GetComponent(entityId, utils.NetworkComponent).(component.Network)
			networkComponent.PingTimestamp = pingEvent.Timestamp
			w.pingQueue <- playerId
		}
	}
}

func (w *World) Update(dt time.Duration) {
	w.engine.UpdateSystems()
}

func (w *World) Stop() {
	w.engine.Stop()
}

func (w *World) processJoinAndLeaveQueue() {
	// Process all players in JoinQueue
	for {
		select {
		case joinEvent := <-w.JoinQueue:
			if err := w.addPlayer(joinEvent); err != nil {
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
		case playerId := <-w.LeaveQueue:
			if err := w.removePlayer(playerId); err != nil {
				utils.Logger.LogError().Msgf("Error removing player %s: %v", playerId, err)
			}
		}
	}
}

func (w *World) newId() types.Id {
	for {
		val := w.idCounter.Load()
		if w.idCounter.CompareAndSwap(val, val+1) {
			return (types.Id(val) % (w.maxId - w.minId)) + w.minId
		}
	}
}

func (w *World) addPlayer(joinEvent *utils.JoinEvent) error {
	// TODO: create entity id and all player components
	return nil
}

func (w *World) removePlayer(playerId string) error {
	// TODO: delete all the components related to the player id
	return nil
}
