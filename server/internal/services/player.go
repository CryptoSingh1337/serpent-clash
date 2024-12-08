package services

import (
	"fmt"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/utils"
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
	Seq                 uint64
	angle               float64
	pingTimestamp       uint32
	lastMouseCoordinate *utils.Coordinate
	speedBoost          bool
}

func NewPlayer(conn *websocket.Conn) *Player {
	player := &Player{
		Id:            uuid.NewString(),
		Conn:          conn,
		Color:         fmt.Sprintf("hsl(%v, 100%%, 50%%)", rand.Intn(360)),
		Seq:           0,
		angle:         0,
		pingTimestamp: 0,
	}
	player.GenerateRandomPosition(utils.DefaultSnakeLength)
	player.lastMouseCoordinate = &player.Segments[0]
	return player
}

func (player *Player) GenerateRandomPosition(length int) {
	totalSnakeLength := float64((length - 1) * utils.SnakeSegmentDistance)
	maxRadius := utils.WorldBoundaryRadius - totalSnakeLength

	theta := rand.Float64() * 2 * math.Pi
	radius := rand.Float64() * maxRadius

	x := radius * math.Cos(theta)
	y := radius * math.Sin(theta)

	segments := make([]utils.Coordinate, length)
	segments[0] = utils.Coordinate{X: x, Y: y}
	for i := 1; i < len(segments); i++ {
		segments[i].X = x - float64(i)*utils.SnakeSegmentDistance*math.Cos(theta)
		segments[i].Y = y - float64(i)*utils.SnakeSegmentDistance*math.Sin(theta)
	}
	player.Segments = segments
}

func (player *Player) Move() {
	mouseCoordinate := player.lastMouseCoordinate
	head := player.Segments[0]
	angle := player.angle
	targetAngle := math.Atan2(mouseCoordinate.Y-head.Y, mouseCoordinate.X-head.X)
	angle = utils.LerpAngle(angle, targetAngle, utils.MaxTurnRate)

	// Move the head towards the mouse coordinate
	speed := float64(utils.PlayerSpeed)
	if player.speedBoost {
		speed += utils.PlayerBoostSpeed
	}
	head.X += math.Cos(angle) * speed
	head.Y += math.Sin(angle) * speed

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

func (player *Player) TeleportTo(x, y float64) *[]utils.Coordinate {
	head := player.Segments[0]
	angle := player.angle
	targetAngle := math.Atan2(y-head.Y, x-head.X)
	angle = utils.LerpAngle(angle, targetAngle, utils.MaxTurnRate)
	player.Segments[0].X = x
	player.Segments[0].Y = y
	for i := 1; i < len(player.Segments); i++ {
		player.Segments[i].X = x - float64(i)*utils.SnakeSegmentDistance*math.Cos(angle)
		player.Segments[i].Y = y - float64(i)*utils.SnakeSegmentDistance*math.Sin(angle)
	}
	return &player.Segments
}
