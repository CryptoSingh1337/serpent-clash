package ecs

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestQuadTree_Insert(t *testing.T) {
	qt := NewQuadTree(Rect{X: 200, Y: 200, W: 200, H: 200}, 4)
	for i := 1; i <= 5; i++ {
		qt.Insert(Point{X: rand.Float64() * qt.Boundary.W, Y: rand.Float64() * qt.Boundary.H})
		fmt.Printf("QT after %v insert: %v\n", i, qt)
	}
}
