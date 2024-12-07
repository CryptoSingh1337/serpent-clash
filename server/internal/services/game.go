package services

import (
	"errors"
	"fmt"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/utils"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"time"
)

type Game struct {
	HashGrid   *SpatialHashGrid
	Players    map[*Player]bool
	Broadcast  chan string
	JoinQueue  chan *Player
	LeaveQueue chan *Player
	PingQueue  chan *Player
	Done       chan bool
}

func NewGame() *Game {
	game := &Game{
		HashGrid:   NewSpatialHashGrid(utils.SnakeSegmentDiameter * 2),
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
				utils.Logger.LogError().Msgf("Error adding player %s: %v", player.Id, err)
				if err := player.Conn.Close(); err != nil {
					utils.Logger.LogError().Msgf("Error closing connection for player %s: %v", player.Id, err)
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
				utils.Logger.LogError().Msgf("Error removing player %s: %v", player.Id, err)
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
				body, _ := utils.ToJsonB(utils.PingEvent{Timestamp: player.pingTimestamp})
				payload, _ := utils.ToJsonB(utils.Payload{Type: utils.PongMessage, Body: body})
				err := player.Conn.WriteMessage(websocket.TextMessage, payload)
				if err != nil {
					utils.Logger.LogError().Msgf("Error pinging player %s: %v", player.Id, err)
				}
			}
		default:
			goto MoveAllPlayers
		}
	}

MoveAllPlayers:
	for player := range game.Players {
		val, ok := game.Players[player]
		if ok && val {
			game.HashGrid.RemovePlayer(player)
			player.Move()
			game.HashGrid.InsertPlayer(player)
		}
	}

	//utils.Logger.LogInfo().Msgf("HashGrid: %v", game.HashGrid.Grid)

	collisions := game.HashGrid.DetectCollisions()
	//utils.Logger.LogInfo().Msgf("Collisions: %v", collisions)
	game.handleCollisions(collisions)

	// form players data in json
	gameState := utils.GameState{
		PlayerStates: make(map[string]utils.PlayerState),
	}
	for player, flag := range game.Players {
		if flag {
			playerState := utils.PlayerState{
				Color:    player.Color,
				Segments: player.Segments,
				Seq:      player.Seq,
			}
			gameState.PlayerStates[player.Id] = playerState
		}
	}
	// Send game state to every player
	body, _ := utils.ToJsonB(gameState)
	payload, _ := utils.ToJsonB(utils.Payload{Type: utils.GameStateMessage, Body: body})
	for player, flag := range game.Players {
		if flag {
			err := player.Conn.WriteMessage(websocket.TextMessage, payload)
			if err != nil {
				if err := player.Conn.Close(); err != nil {
					utils.Logger.LogError().Msgf("Error closing connection for player %s: %v", player.Id, err)
				}
			}
		}
	}
}

func (game *Game) ProcessEvent(player *Player, messageType websocket.MessageType, data []byte) {
	switch messageType {
	case websocket.TextMessage:
		payload, err := utils.FromJsonB[utils.Payload](data)
		if err != nil {
			return
		}
		switch payload.Type {
		case utils.Movement:
			mouseEvent, err := utils.FromJsonB[utils.MouseEvent](payload.Body)
			if err != nil {
				return
			}
			player.lastMouseCoordinate = &mouseEvent.Coordinate
			player.Seq = mouseEvent.Seq
		case utils.PingMessage:
			pingEvent, err := utils.FromJsonB[utils.PingEvent](payload.Body)
			if err != nil {
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
	utils.Logger.LogInfo().Msg("Server::AddPlayer - player Id - " + player.Id + " joined")
	game.Players[player] = true
	//player.GenerateRandomPosition(utils.DefaultSnakeLength)
	body := fmt.Sprintf(`{"id":%q}`, player.Id)
	payload, err := utils.ToJsonB(utils.Payload{Type: utils.HelloMessage, Body: []byte(body)})
	if err != nil {
		return err
	}
	return player.Conn.WriteMessage(websocket.TextMessage, payload)
}

func (game *Game) RemovePlayer(player *Player) error {
	if len(game.Players) == 0 {
		return errors.New("no players left")
	}
	// TODO: break snake into food particles
	if _, ok := game.Players[player]; ok {
		if err := player.Conn.Close(); err != nil {
			utils.Logger.LogError().Msgf("Error closing connection for player %s: %v", player.Id, err)
		}
		game.HashGrid.RemovePlayerById(player.Id)
		delete(game.Players, player)
		utils.Logger.LogInfo().Msgf("Server::RemovePlayer - player Id - %v left", player.Id)
		return nil
	}
	return errors.New("player not exists")
}

func (game *Game) handleCollisions(collisions []Collision) {
	for _, collision := range collisions {
		player1, ok1 := game.getPlayerById(collision.A)
		player2, ok2 := game.getPlayerById(collision.B)

		if ok1 && ok2 {
			// Head-to-body collision
			if game.isHeadToBodyCollision(player1, player2) {
				utils.Logger.LogInfo().Msgf("Player %v collide with Player %v", player1.Id, player2.Id)
				game.killPlayer(player1)
			} else if game.isHeadToBodyCollision(player2, player1) {
				utils.Logger.LogInfo().Msgf("Player %v collide with Player %v", player1.Id, player2.Id)
				game.killPlayer(player2)
			} else {
				utils.Logger.LogInfo().Msgf("Player %v and player %v had head to head collision", player1.Id,
					player2.Id)
				game.killPlayer(player1)
				game.killPlayer(player2)
			}
		}
	}
}

func (game *Game) getPlayerById(id string) (*Player, bool) {
	for player := range game.Players {
		if player.Id == id {
			return player, true
		}
	}
	return nil, false
}

func (game *Game) isHeadToBodyCollision(player1, player2 *Player) bool {
	head1 := player1.Segments[0]
	for i := 1; i < len(player2.Segments); i++ {
		if game.distanceSquared(head1, player2.Segments[i]) < utils.SnakeSegmentDistance*utils.SnakeSegmentDistance {
			return true
		}
	}
	return false
}

func (game *Game) distanceSquared(pos1, pos2 utils.Coordinate) float64 {
	dx := pos1.X - pos2.X
	dy := pos1.Y - pos2.Y
	return dx*dx + dy*dy
}

func (game *Game) killPlayer(player *Player) {
	game.LeaveQueue <- player
}

func (game *Game) Close() {
	utils.Logger.LogInfo().Msgf("Clearing off resources...")
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
