package component

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/utils"
)

type Input struct {
	PrevCoordinates utils.Coordinate
	Coordinates     utils.Coordinate
	Boost           bool
}

func NewInputComponent() Input {
	return Input{
		Coordinates: utils.Coordinate{
			X: 0, Y: 0,
		},
	}
}
