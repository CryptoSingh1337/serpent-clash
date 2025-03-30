package system

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/component"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/storage"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/utils"
	"github.com/lesismal/nbio/nbhttp/websocket"
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
	networkComponents := n.storage.GetAllComponentByName("network").([]component.Network)
	for _, networkComponent := range networkComponents {
		if networkComponent.Connected {
			err := networkComponent.Connection.WriteMessage(websocket.TextMessage, payload)
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
		playerInfoComponent := n.storage.GetComponentByEntityIdAndName(entityId, "playerInfo").(component.PlayerInfo)
		snakeComponent := n.storage.GetComponentByEntityIdAndName(entityId, "snake").(component.Snake)
		networkComponent := n.storage.GetComponentByEntityIdAndName(entityId, "network").(component.Network)
		playerState := utils.PlayerState{
			Color:    snakeComponent.Color,
			Segments: snakeComponent.Segments,
			Seq:      networkComponent.MessageSequence,
		}
		gameState.PlayerStates[playerInfoComponent.ID] = playerState
	}
	return gameState
}
