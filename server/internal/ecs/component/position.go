package component

type Position struct {
	X, Y float64
}

func NewPositionComponent() Position {
	return Position{}
}
