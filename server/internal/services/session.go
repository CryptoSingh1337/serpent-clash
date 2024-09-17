package services

import (
	"errors"
	"github.com/google/uuid"
	"log"
)

const MaxClientPerSession = 10

type Session struct {
	Id         string
	Clients    map[*Client]bool
	Broadcast  chan []string
	Register   chan *Client
	Unregister chan *Client
}

func NewSession() *Session {
	return &Session{
		Id:         uuid.NewString(),
		Broadcast:  make(chan []string),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
}

func (session *Session) addClient(client *Client) error {
	log.Println("Session::AddClient - Adding client: ", client.Id)
	if MaxClientPerSession <= len(session.Clients) {
		return errors.New("max Clients reached")
	}
	if session.Clients[client] {
		return errors.New("client already exists")
	}
	client.Session = session
	session.Register <- client
	return nil
}

func (session *Session) removeClient(client *Client) error {
	log.Println("Session::removeClient - Removing client: ", client.Id)
	if len(session.Clients) == 0 {
		return errors.New("session is empty")
	}
	if session.Clients[client] {
		session.Unregister <- client
		return nil
	}
	return errors.New("client does not exist")
}

func (session *Session) run() {
	for {
		select {
		case client := <-session.Register:
			log.Println("Session::run - client Id - " + client.Id + " registered")
			session.Clients[client] = true
		case client := <-session.Unregister:
			if _, ok := session.Clients[client]; ok {
				log.Println("Session::run - client Id - " + client.Id + " unregistered")
				close(client.Send)
				delete(session.Clients, client)
			}
		case message := <-session.Broadcast:
			for client := range session.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(session.Clients, client)
				}
			}
		}
	}
}
