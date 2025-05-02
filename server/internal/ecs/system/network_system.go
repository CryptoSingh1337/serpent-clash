package system

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/component"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/storage"
	gameutils "github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/utils"
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
	body, _ := gameutils.ToJsonB(gameState)
	payload, _ := gameutils.ToJsonB(gameutils.Payload{Type: gameutils.GameStateMessageType, Body: body})
	networkComponents := n.storage.GetAllComponentByName("network").([]*component.Network)
	pongMessage := gameutils.PongMessage{}
	for _, networkComponent := range networkComponents {
		if networkComponent.Connected {
			err := networkComponent.Connection.WriteMessage(websocket.TextMessage, payload)
			if err != nil {
				networkComponent.Connected = false
				err = networkComponent.Connection.Close()
				if err != nil {
					gameutils.Logger.Err(err).Msgf("error while closing connection for player")
				}
			}
			networkComponent.PingCooldown -= 1
			if networkComponent.PingCooldown <= 0 {
				networkComponent.ResponseInitiateTimestamp = uint64(time.Now().UnixMilli())
				pongMessage.RequestInitiateTimestamp = networkComponent.RequestInitiateTimestamp
				pongMessage.RequestAckTimestamp = networkComponent.RequestAckTimestamp
				pongMessage.ResponseInitiateTimestamp = networkComponent.ResponseInitiateTimestamp
				body, _ = gameutils.ToJsonB(pongMessage)
				pingPayload, _ := gameutils.ToJsonB(gameutils.Payload{Type: gameutils.PongMessageType, Body: body})
				err = networkComponent.Connection.WriteMessage(websocket.TextMessage, pingPayload)
				if err != nil {
					networkComponent.Connected = false
					// TODO: call engine player remove method
				}
				networkComponent.PingCooldown = gameutils.PingCooldown
			}
		}
	}
}

func (n *NetworkSystem) Stop() {

}

func (n *NetworkSystem) createGameState() gameutils.GameStateMessage {
	playerEntityIds := n.storage.GetAllEntitiesByType("player")
	gameState := gameutils.GameStateMessage{
		PlayerStates: make(map[string]gameutils.PlayerStateMessage),
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
		playerState := gameutils.PlayerStateMessage{
			Segments: snakeComponent.Segments,
			Seq:      networkComponent.MessageSequence,
		}
		gameState.PlayerStates[playerInfoComponent.ID] = playerState
	}
	return gameState
}
