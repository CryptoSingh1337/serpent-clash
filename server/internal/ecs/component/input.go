package component

import (
	gameutils "github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/utils"
)

type Input struct {
	Coordinates gameutils.Coordinate
	Boost       bool
}

func NewInputComponent() Input {
	return Input{
		Coordinates: gameutils.Coordinate{
			X: 0, Y: 0,
		},
	}
}
