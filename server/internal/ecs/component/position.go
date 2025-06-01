package component

type Position struct {
	X, Y float64
}

func NewPositionComponent(x, y float64) Position {
	return Position{x, y}
}
