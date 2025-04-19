package component

import (
	"fmt"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/utils"
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
		Color:   fmt.Sprintf("hsl(%v, 100%%, 50%%)", rand.Intn(360)),
		IsAlive: true,
		Angle:   0,
	}
}
