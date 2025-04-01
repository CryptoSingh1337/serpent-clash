package component

import (
	"github.com/lesismal/nbio/nbhttp/websocket"
)

type Network struct {
	Connection      *websocket.Conn
	Connected       bool
	PingTimestamp   uint64
	MessageSequence uint64
}

func NewNetworkComponent(connection *websocket.Conn) Network {
	return Network{
		Connection: connection,
		Connected:  true,
	}
}
