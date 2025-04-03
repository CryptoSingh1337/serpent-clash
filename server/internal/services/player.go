package services

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"math"
)

type Player struct {
	Id                  string
	Username            *string
	Conn                *websocket.Conn
	W                   *echo.Response
	Segments            []utils.Coordinate
	Color               string
	Seq                 uint64
	angle               float64
	pingTimestamp       uint64
	lastMouseCoordinate *utils.Coordinate
	speedBoost          bool
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
