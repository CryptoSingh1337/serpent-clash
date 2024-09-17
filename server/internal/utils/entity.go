package utils

import "math/rand"

const (
	worldFactor         = 120
	worldWidth          = 16 * worldFactor
	worldHeight         = 9 * worldFactor
	playerSpeed         = 10
	defaultSnakeLength  = 10
	defaultGrowthFactor = 2
)

type Position struct {
	X int
	Y int
}

type Player struct {
	Positions []Position
}

type Food struct {
	Point Position
}

func SpawnPlayer() *Player {
	x := 100 + rand.Intn(worldWidth-100)
	y := 100 + rand.Intn(worldHeight-100)

	var positions []Position
	positions[0] = Position{x, y}
	for i := 1; i <= len(positions); i++ {
		positions[i].X = x + defaultGrowthFactor
		positions[i].Y = y + defaultSnakeLength
	}
	return &Player{
		Positions: positions,
	}
}
