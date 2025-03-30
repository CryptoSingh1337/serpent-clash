package component

import "github.com/CryptoSingh1337/serpent-clash/server/internal/utils"

type Snake struct {
	Color    string
	Segments []utils.Coordinate
	IsAlive  bool
	Angle    float64
}
