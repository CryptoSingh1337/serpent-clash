package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/CryptoSingh1337/multiplayer-snake-game/server/internal/utils"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"log"
	"time"
)

type Game struct {
	Players    map[*Player]bool
	Broadcast  chan string
	JoinQueue  chan *Player
	LeaveQueue chan *Player
	PingQueue  chan *Player
	Done       chan bool
}

func NewGame() *Game {
	game := &Game{
		Players:    make(map[*Player]bool),
		Broadcast:  make(chan string),
		JoinQueue:  make(chan *Player, utils.MaxPlayerAllowed),
		LeaveQueue: make(chan *Player, utils.MaxPlayerAllowed),
		PingQueue:  make(chan *Player, utils.MaxPlayerAllowed),
		Done:       make(chan bool),
	}
	game.init()
	return game
}

func (game *Game) init() {
	ticker := time.NewTicker(1000 / utils.TickRate * time.Millisecond)
	go func() {
		for {
			select {
			case <-game.Done:
				ticker.Stop()
				return
			case _ = <-ticker.C:
				game.processTick()
			}
		}
	}()
}

func (game *Game) processTick() {
	// Process all players in JoinQueue
	for {
		select {
		case player := <-game.JoinQueue:
			if err := game.AddPlayer(player); err != nil {
				log.Printf("Error adding player %s: %v", player.Id, err)
				if err := player.Conn.Close(); err != nil {
					log.Printf("Error closing connection for player %s: %v", player.Id, err)
				}
			}
		default:
			// Exit the loop when JoinQueue is empty
			goto ProcessLeaveQueue
		}
	}

ProcessLeaveQueue:
	// Process all players in LeaveQueue
	for {
		select {
		case player := <-game.LeaveQueue:
			if err := game.RemovePlayer(player); err != nil {
				log.Printf("Error removing player %s: %v", player.Id, err)
			}
		default:
			// Exit the loop when LeaveQueue is empty
			goto ProcessPingQueue
		}
	}

ProcessPingQueue:
	// Process all players in PingQueue
	for {
		select {
		case player := <-game.PingQueue:
			val, ok := game.Players[player]
			if ok && val {
				log.Println("ping timestamp", player.pingTimestamp)
				body, _ := json.Marshal(utils.PingEvent{Timestamp: player.pingTimestamp})
				payload, _ := json.Marshal(utils.Payload{Type: utils.PongMessage, Body: body})
				err := player.Conn.WriteMessage(websocket.TextMessage, payload)
				if err != nil {
					log.Printf("Error pinging player %s: %v", player.Id, err)
				}
			}
		default:
			goto BroadcastGameState
		}
	}

BroadcastGameState:
	// form players data in json
	gameState := utils.GameState{
		PlayerStates: make(map[string]utils.PlayerState),
	}
	for player, flag := range game.Players {
		if flag {
			playerState := utils.PlayerState{
				Color:    player.Color,
				Segments: player.Segments,
			}
			gameState.PlayerStates[player.Id] = playerState
		}
	}
	// Send game state to every player
	body, _ := json.Marshal(gameState)
	payload, _ := json.Marshal(utils.Payload{Type: utils.GameStateMessage, Body: body})
	for player, flag := range game.Players {
		if flag {
			err := player.Conn.WriteMessage(websocket.TextMessage, payload)
			if err != nil {
				if err := player.Conn.Close(); err != nil {
					log.Printf("Error closing connection for player %s: %v", player.Id, err)
				}
			}
		}
	}
}

func (game *Game) ProcessEvent(player *Player, messageType websocket.MessageType, data []byte) {
	switch messageType {
	case websocket.TextMessage:
		payload := utils.Payload{}
		if err := json.Unmarshal(data, &payload); err != nil {
			return
		}
		switch payload.Type {
		case utils.Movement:
			mouseEvent := utils.MouseEvent{}
			if err := json.Unmarshal(payload.Body, &mouseEvent); err != nil {
				return
			}
			player.Move(&mouseEvent.Coordinate)
		case utils.PingMessage:
			pingEvent := utils.PingEvent{}
			if err := json.Unmarshal(payload.Body, &pingEvent); err != nil {
				return
			}
			player.pingTimestamp = pingEvent.Timestamp
			game.PingQueue <- player
		}
	}
}

func (game *Game) AddPlayer(player *Player) error {
	if utils.MaxPlayerAllowed <= len(game.Players) {
		return errors.New("max players reached")
	}
	if game.Players[player] {
		return errors.New("player already exists")
	}
	log.Println("Server::AddPlayer - player Id - " + player.Id + " joined")
	game.Players[player] = true
	player.generateRandomPosition()
	body := fmt.Sprintf(`{"id":%q}`, player.Id)
	payload, _ := json.Marshal(utils.Payload{Type: utils.HelloMessage, Body: []byte(body)})
	return player.Conn.WriteMessage(websocket.TextMessage, payload)
}

func (game *Game) RemovePlayer(player *Player) error {
	if len(game.Players) == 0 {
		return errors.New("no players left")
	}
	if _, ok := game.Players[player]; ok {
		log.Println("Server::RemovePlayer - player Id - " + player.Id + " left")
		if err := player.Conn.Close(); err != nil {
			log.Printf("Error closing connection for player %s: %v", player.Id, err)
		}
		delete(game.Players, player)
		return nil
	}
	return errors.New("player not exists")
}

func (game *Game) Close() {
	log.Println("Clearing off resources...")
	// Stop ticker
	game.Done <- true

	// Close all ws connections
	for player := range game.Players {
		if err := player.Conn.Close(); err != nil {
			return
		}
	}

	time.Sleep(2 * time.Second)

	// Close all channels
	close(game.Done)
	close(game.JoinQueue)
	close(game.LeaveQueue)
	close(game.PingQueue)
	close(game.Broadcast)
}
