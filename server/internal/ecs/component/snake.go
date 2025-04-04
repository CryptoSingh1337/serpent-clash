package component

import (
	"fmt"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/utils"
	"math"
	"math/rand"
)

type Snake struct {
	Color    string
	Segments []utils.Coordinate
	IsAlive  bool
	Angle    float64
}

func NewSnakeComponent() Snake {
	return Snake{
		Color:    fmt.Sprintf("hsl(%v, 100%%, 50%%)", rand.Intn(360)),
		Segments: generateRandomPosition(utils.DefaultSnakeLength),
		IsAlive:  true,
		Angle:    0,
	}
}

func generateRandomPosition(length int) []utils.Coordinate {
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
	return segments
}
