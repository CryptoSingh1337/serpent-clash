package system

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/component"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/storage"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/types"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/utils"
	"github.com/gorilla/websocket"
	"time"
)

type PlayerState struct {
	Username string             `json:"username"`
	Segments []utils.Coordinate `json:"positions"`
	Seq      uint64             `json:"seq"`
}

type PlayerStatePayload struct {
	Players map[string]PlayerState `json:"players"`
}

type FoodState struct {
	Coordinate utils.Coordinate `json:"coordinate"`
}

type FoodStatePayload struct {
	Foods map[types.Id]FoodState `json:"foods"`
}

type NetworkSystem struct {
	storage storage.Storage
}

func NewNetworkSystem(storage storage.Storage) System {
	return &NetworkSystem{
		storage: storage,
	}
}

func (n *NetworkSystem) Name() string {
	return utils.NetworkSystemName
}

func (n *NetworkSystem) Update() {
	playerState := n.createPlayerState()
	foodState := n.createFoodState()

	// Broadcast game state to all players
	n.broadcast(utils.PlayerStateMessageType, playerState)
	n.broadcast(utils.FoodStateMessageType, foodState)

	// Send ping message to all players
	networkComponents := n.storage.GetAllComponentByName(utils.NetworkComponent).([]*component.Network)
	pongMessage := utils.PongMessage{}
	for _, networkComponent := range networkComponents {
		if networkComponent.Connected {
			networkComponent.PingCooldown -= 1
			if networkComponent.PingCooldown <= 0 {
				networkComponent.ResponseInitiateTimestamp = uint64(time.Now().UnixMilli())
				pongMessage.RequestInitiateTimestamp = networkComponent.RequestInitiateTimestamp
				pongMessage.RequestAckTimestamp = networkComponent.RequestAckTimestamp
				pongMessage.ResponseInitiateTimestamp = networkComponent.ResponseInitiateTimestamp
				body, _ := utils.ToJsonB(pongMessage)
				pingPayload, _ := utils.ToJsonB(utils.Payload{Type: utils.PongMessageType, Body: body})
				err := networkComponent.Connection.WriteMessage(websocket.TextMessage, pingPayload)
				if err != nil {
					networkComponent.Connected = false
					err = networkComponent.Connection.Close()
				}
				networkComponent.PingCooldown = utils.PingCooldown
			}
		}
	}
}

func (n *NetworkSystem) Stop() {

}

func (n *NetworkSystem) broadcast(msgType string, data any) {
	body, _ := utils.ToJsonB(data)
	payload, _ := utils.ToJsonB(utils.Payload{Type: msgType, Body: body})
	networkComponents := n.storage.GetAllComponentByName(utils.NetworkComponent).([]*component.Network)
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
		}
	}
}

func (n *NetworkSystem) createPlayerState() PlayerStatePayload {
	playerEntityIds := n.storage.GetAllEntitiesByType(utils.PlayerEntity)
	playerStatePayload := PlayerStatePayload{
		Players: make(map[string]PlayerState),
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
		playerState := PlayerState{
			Username: playerInfoComponent.Username,
			Segments: snakeComponent.Segments,
			Seq:      networkComponent.MessageSequence,
		}
		playerStatePayload.Players[playerInfoComponent.ID] = playerState
	}
	return playerStatePayload
}

func (n *NetworkSystem) createFoodState() FoodStatePayload {
	foodState := FoodStatePayload{
		Foods: make(map[types.Id]FoodState),
	}
	foodEntityIds := n.storage.GetAllEntitiesByType(utils.FoodEntity)
	for _, entityId := range foodEntityIds {
		c := n.storage.GetComponentByEntityIdAndName(entityId, utils.PositionComponent)
		if c == nil {
			continue
		}
		positionComponent := c.(*component.Position)
		foodState.Foods[entityId] = FoodState{
			Coordinate: utils.Coordinate{
				X: positionComponent.X,
				Y: positionComponent.Y,
			},
		}
	}
	return foodState
}
