package main

import (
	"errors"
	"log"
)

const MaxSessions = 1

type Game struct {
	sessions []*Session
}

func NewGame() *Game {
	return &Game{
		sessions: make([]*Session, 0, MaxSessions),
	}
}

func (game *Game) addClient(client *Client) error {
	log.Println("Server::addClient - Adding client in the game")
	for _, session := range game.sessions {
		log.Printf("Server::addClient - Session id - %v, client count - %v\n",
			session.id, len(session.clients))
		if MaxClientPerSession > len(session.clients) {
			err := session.addClient(client)
			if err != nil {
				return err
			}
			return nil
		}
	}
	if MaxSessions > len(game.sessions) {
		session := NewSession()
		game.addSession(session)
		err := session.addClient(client)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("all sessions are occupied")
}

func (game *Game) addSession(session *Session) {
	game.sessions = append(game.sessions, session)
	go session.run()
}

func (game *Game) removeSession(_session *Session) {
	for i, session := range game.sessions {
		if session == _session {
			game.sessions = append(game.sessions[:i], game.sessions[i+1:]...)
		}
	}
}
