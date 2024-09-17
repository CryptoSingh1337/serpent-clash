package services

import (
	"errors"
	"log"
)

const MaxSessions = 1

type Game struct {
	Sessions []*Session
}

func NewGame() *Game {
	return &Game{
		Sessions: make([]*Session, 0, MaxSessions),
	}
}

func (game *Game) AddClient(client *Client) error {
	log.Println("Server::AddClient - Adding client in the game")
	for _, session := range game.Sessions {
		log.Printf("Server::AddClient - Session Id - %v, client count - %v\n",
			session.Id, len(session.Clients))
		if MaxClientPerSession > len(session.Clients) {
			err := session.addClient(client)
			if err != nil {
				return err
			}
			return nil
		}
	}
	if MaxSessions > len(game.Sessions) {
		session := NewSession()
		game.AddSession(session)
		err := session.addClient(client)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("all Sessions are occupied")
}

func (game *Game) RemoveClient(client *Client) error {
	log.Println("Server::RemoveClient - Removing client in the game")
	err := client.Session.removeClient(client)
	if err != nil {
		return err
	}
	return nil
}

func (game *Game) AddSession(session *Session) {
	game.Sessions = append(game.Sessions, session)
	go session.run()
}

func (game *Game) RemoveSession(_session *Session) {
	for i, session := range game.Sessions {
		if session == _session {
			game.Sessions = append(game.Sessions[:i], game.Sessions[i+1:]...)
		}
	}
}
