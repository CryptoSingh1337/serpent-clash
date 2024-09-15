package main

import (
	"github.com/gorilla/websocket"
	"log"
)

const (
	maxMessageSize = 512
)

type Client struct {
	id      string
	session *Session
	conn    *websocket.Conn
	send    chan []string
}

type Payload struct {
	EventType string `json:"eventType"`
	Body      string `json:"body"`
}

func (client *Client) read() {
	defer func() {
		log.Println("Closing the reading pipe for client id - " + client.id)
		client.session.unregister <- client
		_ = client.conn.Close()
	}()
	client.conn.SetReadLimit(maxMessageSize)
	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		log.Println("Received message: " + string(message))
	}
}

func (client *Client) write() {
	defer func() {
		log.Println("Closing the writing pipe for client id - " + client.id)
		_ = client.conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.send:
			if !ok {
				// The hub closed the channel.
				_ = client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := client.conn.WriteJSON(message)
			if err != nil {
				return
			}
		}
	}
}
