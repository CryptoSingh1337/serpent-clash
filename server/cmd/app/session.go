package main

import (
	"errors"
	"github.com/google/uuid"
	"log"
)

const MaxClientPerSession = 10

type Session struct {
	id         string
	clients    map[*Client]bool
	broadcast  chan []string
	register   chan *Client
	unregister chan *Client
}

func NewSession() *Session {
	return &Session{
		id:         uuid.NewString(),
		broadcast:  make(chan []string),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (session *Session) addClient(client *Client) error {
	log.Println("Session::addClient - Adding client: ", client.id)
	if MaxClientPerSession <= len(session.clients) {
		return errors.New("max clients reached")
	}
	if session.clients[client] {
		return errors.New("client already exists")
	}
	client.session = session
	session.register <- client
	return nil
}

func (session *Session) removeClient(client *Client) error {
	log.Println("Session::removeClient - Removing client: ", client.id)
	if len(session.clients) == 0 {
		return errors.New("session is empty")
	}
	if session.clients[client] {
		session.unregister <- client
		return nil
	}
	return errors.New("client does not exist")
}

func (session *Session) run() {
	for {
		select {
		case client := <-session.register:
			log.Println("Session::run - client id - " + client.id + " registered")
			session.clients[client] = true
		case client := <-session.unregister:
			if _, ok := session.clients[client]; ok {
				log.Println("Session::run - client id - " + client.id + " unregistered")
				close(client.send)
				delete(session.clients, client)
			}
		case message := <-session.broadcast:
			for client := range session.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(session.clients, client)
				}
			}
		}
	}
}
