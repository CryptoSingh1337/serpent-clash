package services

import (
	"encoding/json"
	"github.com/CryptoSingh1337/multiplayer-snake-game/server/internal/utils"
	"github.com/lesismal/nbio/nbhttp/websocket"
)

const (
	maxMessageSize = 512
)

type Client struct {
	Id      string
	Session *Session
	Conn    *websocket.Conn
	Player  *utils.Player
	Send    chan []string
}

type Payload struct {
	Type string          `json:"type"`
	Body json.RawMessage `json:"body"`
}

//func (client *Client) Read() {
//	defer func() {
//		log.Println("Closing the reading pipe for client Id - " + client.Id)
//		client.Session.Unregister <- client
//		_ = client.Conn.Close()
//	}()
//	client.Conn.SetReadLimit(maxMessageSize)
//	for {
//		//message := Payload{}
//		_, message, err := client.Conn.ReadMessage()
//		//err := client.Conn.ReadJSON(message)
//		if err != nil {
//			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
//				log.Printf("error: %v", err)
//			}
//			break
//		}
//		log.Println("Received a message:", string(message))
//		//log.Printf("Received message: type - %v, body - %v\n", message.Type, message.Body)
//	}
//}
//
//func (client *Client) Write() {
//	defer func() {
//		log.Println("Closing the writing pipe for client Id - " + client.Id)
//		_ = client.Conn.Close()
//	}()
//
//	for {
//		select {
//		case message, ok := <-client.Send:
//			if !ok {
//				// The hub closed the channel.
//				_ = client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
//				return
//			}
//
//			err := client.Conn.WriteJSON(message)
//			if err != nil {
//				return
//			}
//		}
//	}
//}
