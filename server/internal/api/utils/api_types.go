package utils

import "github.com/gorilla/websocket"

type WSHandler struct {
	OnOpen    func(conn *websocket.Conn)
	OnMessage func(playerId string, messageType int, data []byte)
	OnClose   func(playerId string, err error)
}
