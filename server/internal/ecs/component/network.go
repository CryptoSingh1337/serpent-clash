package component

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/utils"
	"github.com/gorilla/websocket"
)

type Network struct {
	Connection                *websocket.Conn
	Connected                 bool
	RequestInitiateTimestamp  uint64
	RequestAckTimestamp       uint64
	ResponseInitiateTimestamp uint64
	MessageSequence           uint64
	PingCooldown              uint
}

func NewNetworkComponent(connection *websocket.Conn) Network {
	return Network{
		Connection:   connection,
		Connected:    true,
		PingCooldown: utils.PingCooldown,
	}
}
