package storage

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestQuadTree_Insert(t *testing.T) {
	qt := NewQuadTree(BBox{X: 200, Y: 200, W: 200, H: 200}, 4)
	for i := 1; i <= 10; i++ {
		qt.Insert(Point{X: rand.Float64() * qt.Boundary.W, Y: rand.Float64() * qt.Boundary.H})
	}
	qt.Print()
}

func TestQuadTree_QueryBCircle(t *testing.T) {
	qt := NewQuadTree(BBox{X: 200, Y: 200, W: 200, H: 200}, 4)
	for i := 1; i <= 1000; i++ {
		qt.Insert(Point{X: rand.Float64() * qt.Boundary.W, Y: rand.Float64() * qt.Boundary.H})
	}
	qt.Print()

	var points []Point
	qt.QueryBCircle(BCircle{X: 150, Y: 150, R: 10}, &points)
	for _, point := range points {
		fmt.Println(point)
	}
}
