package types

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/component"
	"github.com/gorilla/websocket"
)

type Id uint32

type Component interface {
	component.Expiry | component.Input | component.Network | component.PlayerInfo | component.Position | component.Snake
}

type JoinEvent struct {
	Connection *websocket.Conn
	EntityId   Id
	PlayerId   string
	Username   string
}

type LeaveEvent struct {
	EntityId Id
	PlayerId string
}

type PingEvent struct {
	PlayerId                 string
	RequestInitiateTimestamp uint64 `json:"reqInit"`
}
