package services

import (
	"errors"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"log"
	"strconv"
	"time"
)

const (
	tickRate            = 1
	playerSpeed         = 10
	defaultSnakeLength  = 10
	defaultGrowthFactor = 2
	maxPlayerAllowed    = 10
	worldFactor         = 120
	worldHeight         = 9 * worldFactor
	worldWidth          = 16 * worldFactor
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
		JoinQueue:  make(chan *Player),
		LeaveQueue: make(chan *Player),
		PingQueue:  make(chan *Player),
		Done:       make(chan bool),
	}
	game.init()
	return game
}

func (game *Game) init() {
	ticker := time.NewTicker(1000 / tickRate * time.Millisecond)
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

func (game *Game) Close() {
	log.Println("Clearing off resources...")
	// Stop ticker
	game.Done <- true

	// Close all ws connections
	for player, _ := range game.Players {
		if err := player.Conn.Close(); err != nil {
			return
		}
	}

	// Close all channels
	close(game.Broadcast)
	close(game.JoinQueue)
	close(game.LeaveQueue)
	close(game.PingQueue)
	close(game.Done)
}

func (game *Game) AddPlayer(player *Player) error {
	if maxPlayerAllowed <= len(game.Players) {
		return errors.New("max players reached")
	}
	if game.Players[player] {
		return errors.New("player already exists")
	}
	log.Println("Server::AddPlayer - player Id - " + player.Id + " joined")
	game.Players[player] = true
	return nil
}

func (game *Game) RemovePlayer(player *Player) error {
	if len(game.Players) == 0 {
		return errors.New("no players left")
	}
	if _, ok := game.Players[player]; ok {
		log.Println("Server::RemovePlayer - player Id - " + player.Id + " left")
		delete(game.Players, player)
		return nil
	}
	return errors.New("player not exists")
}

func (game *Game) processTick() {
	log.Println("Server::processTick")
	timestamp := time.Now().Unix()

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
			if err := player.Conn.Close(); err != nil {
				log.Printf("Error closing connection for player %s: %v", player.Id, err)
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
			err := player.Conn.WriteMessage(websocket.PongMessage, []byte(strconv.FormatInt(timestamp, 10)))
			if err != nil {
				log.Printf("Error pinging player %s: %v", player.Id, err)
			}
		default:
			return
		}
	}
}
