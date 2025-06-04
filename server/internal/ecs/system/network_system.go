package system

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/component"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/storage"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/utils"
	"github.com/gorilla/websocket"
	"time"
)

type NetworkSystem struct {
	storage storage.Storage
}

func NewNetworkSystem(storage storage.Storage) System {
	return &NetworkSystem{
		storage,
	}
}

func (n *NetworkSystem) Update() {
	gameState := n.createGameState()
	body, _ := utils.ToJsonB(gameState)
	payload, _ := utils.ToJsonB(utils.Payload{Type: utils.GameStateMessageType, Body: body})
	networkComponents := n.storage.GetAllComponentByName(utils.NetworkComponent).([]*component.Network)
	pongMessage := utils.PongMessage{}
	for _, networkComponent := range networkComponents {
		if networkComponent.Connected {
			err := networkComponent.Connection.WriteMessage(websocket.TextMessage, payload)
			if err != nil {
				networkComponent.Connected = false
				err = networkComponent.Connection.Close()
				if err != nil {
					utils.Logger.Err(err).Msgf("error while closing connection for player")
				}
			}
			networkComponent.PingCooldown -= 1
			if networkComponent.PingCooldown <= 0 {
				networkComponent.ResponseInitiateTimestamp = uint64(time.Now().UnixMilli())
				pongMessage.RequestInitiateTimestamp = networkComponent.RequestInitiateTimestamp
				pongMessage.RequestAckTimestamp = networkComponent.RequestAckTimestamp
				pongMessage.ResponseInitiateTimestamp = networkComponent.ResponseInitiateTimestamp
				body, _ = utils.ToJsonB(pongMessage)
				pingPayload, _ := utils.ToJsonB(utils.Payload{Type: utils.PongMessageType, Body: body})
				err = networkComponent.Connection.WriteMessage(websocket.TextMessage, pingPayload)
				if err != nil {
					networkComponent.Connected = false
					// TODO: call engine player remove method
				}
				networkComponent.PingCooldown = utils.PingCooldown
			}
		}
	}
}

func (n *NetworkSystem) Stop() {

}

func (n *NetworkSystem) createGameState() utils.GameStateMessage {
	playerEntityIds := n.storage.GetAllEntitiesByType(utils.PlayerEntity)
	gameState := utils.GameStateMessage{
		PlayerStates: make(map[string]utils.PlayerStateMessage),
	}
	for _, entityId := range playerEntityIds {
		c := n.storage.GetComponentByEntityIdAndName(entityId, utils.PlayerInfoComponent)
		if c == nil {
			continue
		}
		playerInfoComponent := c.(*component.PlayerInfo)
		c = n.storage.GetComponentByEntityIdAndName(entityId, utils.SnakeComponent)
		if c == nil {
			continue
		}
		snakeComponent := c.(*component.Snake)
		c = n.storage.GetComponentByEntityIdAndName(entityId, utils.NetworkComponent)
		if c == nil {
			continue
		}
		networkComponent := c.(*component.Network)
		playerState := utils.PlayerStateMessage{
			Segments: snakeComponent.Segments,
			Seq:      networkComponent.MessageSequence,
		}
		gameState.PlayerStates[playerInfoComponent.ID] = playerState
	}
	return gameState
}
