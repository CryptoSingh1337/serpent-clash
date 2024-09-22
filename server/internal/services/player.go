package services

import (
	"fmt"
	"github.com/CryptoSingh1337/multiplayer-snake-game/server/internal/utils"
	"github.com/google/uuid"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"math"
	"math/rand"
)

type Player struct {
	Id                  string
	Conn                *websocket.Conn
	Segments            []utils.Coordinate
	Color               string
	angle               float64
	pingTimestamp       uint32
	lastMouseCoordinate *utils.Coordinate
}

func NewPlayer(conn *websocket.Conn) *Player {
	player := &Player{
		Id:            uuid.NewString(),
		Conn:          conn,
		Color:         fmt.Sprintf("hsl(%v, 100%%, 50%%)", rand.Intn(360)),
		angle:         0,
		pingTimestamp: 0,
	}
	player.generateRandomPosition()
	player.lastMouseCoordinate = &player.Segments[0]
	return player
}

func (player *Player) generateRandomPosition() {
	x := float64(100 + rand.Intn(utils.WorldWidth-100))
	y := float64(100 + rand.Intn(utils.WorldHeight-100))

	segments := make([]utils.Coordinate, utils.DefaultSnakeLength)
	segments[0] = utils.Coordinate{X: x, Y: y}
	for i := 1; i < len(segments); i++ {
		segments[i].X = x - float64(i)*utils.SnakeSegmentDistance
		segments[i].Y = y
	}
	player.Segments = segments
}

func (player *Player) Move() {
	mouseCoordinate := player.lastMouseCoordinate
	//log.Println("Mouse coordinate", mouseCoordinate)
	head := player.Segments[0]
	angle := player.angle
	targetAngle := math.Atan2(mouseCoordinate.Y-head.Y, mouseCoordinate.X-head.X)
	angle = utils.LerpAngle(angle, targetAngle, utils.MaxTurnRate)

	// Move the head towards the mouse coordinate
	head.X += math.Cos(angle) * utils.PlayerSpeed
	head.Y += math.Sin(angle) * utils.PlayerSpeed

	// Update the head position and angle
	player.Segments[0] = head
	player.angle = angle

	// Move the rest of the snake to follow the head
	for i := 1; i < len(player.Segments); i++ {
		prevSegment := player.Segments[i-1]
		currentSegment := player.Segments[i]

		angleToPrev := math.Atan2(prevSegment.Y-currentSegment.Y, prevSegment.X-currentSegment.X)

		// Keep a fixed distance between segments
		currentSegment.X = prevSegment.X - math.Cos(angleToPrev)*utils.SnakeSegmentDistance
		currentSegment.Y = prevSegment.Y - math.Sin(angleToPrev)*utils.SnakeSegmentDistance
		player.Segments[i] = currentSegment
	}
}
