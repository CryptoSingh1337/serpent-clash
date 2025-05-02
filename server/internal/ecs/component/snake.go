package component

import (
	gameutils "github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/utils"
)

type Snake struct {
	Segments []gameutils.Coordinate
	IsAlive  bool
	Angle    float64
}

func NewSnakeComponent() Snake {
	return Snake{
		IsAlive: true,
		Angle:   0,
	}
}
