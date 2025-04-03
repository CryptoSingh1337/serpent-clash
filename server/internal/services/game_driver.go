package services

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/utils"
	"math"
)

type GameDriver struct {
	HashGrid   *SpatialHashGrid
	Players    map[*Player]bool
	Broadcast  chan *string
	JoinQueue  chan *Player
	LeaveQueue chan *Player
	PingQueue  chan *Player
	Done       chan bool
}

func (game *GameDriver) TeleportPlayer(playerId string, coordinate *utils.Coordinate) (*[]utils.Coordinate, bool) {
	player, ok := game.getPlayerById(playerId)
	if ok {
		return player.TeleportTo(coordinate.X, coordinate.Y), true
	}
	return nil, false
}

func (game *GameDriver) handleCollisions(collisions []Collision) {
	for _, collision := range collisions {
		a, ok1 := game.getPlayerById(collision.A)
		b, ok2 := game.getPlayerById(collision.B)
		if ok1 && ok2 {
			// Head-to-body collision
			if game.isHeadToBodyCollision(a, b) {
				utils.Logger.LogInfo().Msgf("Player %v collide with Player %v", a.Id, b.Id)
				game.killPlayer(a)
			} else if game.isHeadToBodyCollision(b, a) {
				utils.Logger.LogInfo().Msgf("Player %v collide with Player %v", a.Id, b.Id)
				game.killPlayer(b)
			} else {
				utils.Logger.LogInfo().Msgf("Player %v and player %v had head to head collision", a.Id, b.Id)
				game.killPlayer(a)
				game.killPlayer(b)
			}
		}
	}
}

func (game *GameDriver) getPlayerById(id string) (*Player, bool) {
	for player := range game.Players {
		if player.Id == id {
			return player, true
		}
	}
	return nil, false
}

func (game *GameDriver) isHeadToBodyCollision(a, b *Player) bool {
	head := a.Segments[0]
	for i := 1; i < len(b.Segments); i++ {
		distance := utils.EuclideanDistance(head.X, head.Y, b.Segments[i].X, b.Segments[i].Y) +
			utils.SnakeSegmentDiameter/2
		if math.Floor(distance) == 0 {
			return true
		}
	}
	return false
}

func (game *GameDriver) killPlayer(player *Player) {
	game.LeaveQueue <- player
}
