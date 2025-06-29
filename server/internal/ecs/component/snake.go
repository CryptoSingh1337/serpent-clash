package component

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/utils"
)

type Snake struct {
	Segments        []utils.Coordinate
	Angle           float64
	FoodConsumed    uint64
	GrowthThreshold uint
}

func NewSnakeComponent() Snake {
	return Snake{
		Angle:           0,
		FoodConsumed:    0,
		GrowthThreshold: utils.DefaultGrowthFactor,
	}
}
