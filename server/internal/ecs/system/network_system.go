package system

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/component"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/storage"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/utils"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"time"
)

type NetworkSystem struct {
	storage storage.Storage
}

func NewNetworkSystem(storage storage.Storage) *NetworkSystem {
	return &NetworkSystem{
		storage,
	}
}

func (n *NetworkSystem) Update() {
	gameState := n.createGameState()
	body, _ := utils.ToJsonB(gameState)
	payload, _ := utils.ToJsonB(utils.Payload{Type: utils.GameStateMessage, Body: body})
	networkComponents := n.storage.GetAllComponentByName("network").([]*component.Network)
	pingEvent := utils.PingMessageEvent{}
	for _, networkComponent := range networkComponents {
		if networkComponent.Connected {
			err := networkComponent.Connection.WriteMessage(websocket.TextMessage, payload)
			if err != nil {
				networkComponent.Connected = false
				// TODO: call engine player remove method
			}
			networkComponent.ResponseInitiateTimestamp = uint64(time.Now().UnixMilli())
			pingEvent.Timestamp = networkComponent.RequestInitiateTimestamp
			pingEvent.RequestInitiateTimestamp = networkComponent.RequestInitiateTimestamp
			pingEvent.RequestAckTimestamp = networkComponent.RequestAckTimestamp
			pingEvent.ResponseInitiateTimestamp = networkComponent.ResponseInitiateTimestamp
			body, _ = utils.ToJsonB(pingEvent)
			pingPayload, _ := utils.ToJsonB(utils.Payload{Type: utils.PongMessage, Body: body})
			err = networkComponent.Connection.WriteMessage(websocket.TextMessage, pingPayload)
			if err != nil {
				networkComponent.Connected = false
				// TODO: call engine player remove method
			}
		}
	}
}

func (n *NetworkSystem) Stop() {

}

func (n *NetworkSystem) createGameState() utils.GameState {
	playerEntityIds := n.storage.GetAllEntitiesByType("player")
	gameState := utils.GameState{
		PlayerStates: make(map[string]utils.PlayerState),
	}
	for _, entityId := range playerEntityIds {
		c := n.storage.GetComponentByEntityIdAndName(entityId, "playerInfo")
		if c == nil {
			continue
		}
		playerInfoComponent := c.(*component.PlayerInfo)
		c = n.storage.GetComponentByEntityIdAndName(entityId, "snake")
		if c == nil {
			continue
		}
		snakeComponent := c.(*component.Snake)
		c = n.storage.GetComponentByEntityIdAndName(entityId, "network")
		if c == nil {
			continue
		}
		networkComponent := c.(*component.Network)
		playerState := utils.PlayerState{
			Color:    snakeComponent.Color,
			Segments: snakeComponent.Segments,
			Seq:      networkComponent.MessageSequence,
		}
		gameState.PlayerStates[playerInfoComponent.ID] = playerState
	}
	return gameState
}
