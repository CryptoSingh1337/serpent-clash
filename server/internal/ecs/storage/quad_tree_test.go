package storage

import (
	"fmt"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/utils"
	"math/rand"
	"testing"
)

func TestQuadTree_Insert(t *testing.T) {
	qt := NewQuadTree(BBox{X: 200, Y: 200, W: 200, H: 200}, 4)
	for i := 1; i <= 10; i++ {
		qt.Insert(Point{X: rand.Float64() * qt.Boundary.W, Y: rand.Float64() * qt.Boundary.H})
	}
	qt.Print(true)
}

func TestQuadTree_QueryBCircle(t *testing.T) {
	qt := NewQuadTree(BBox{X: 200, Y: 200, W: 200, H: 200}, 15)
	for i := 1; i <= 10000; i++ {
		qt.Insert(Point{X: rand.Float64() * qt.Boundary.W, Y: rand.Float64() * qt.Boundary.H})
	}
	qt.Print(true)

	var points []Point
	qt.QueryBCircle(BCircle{X: 150, Y: 150, R: 10}, &points)
	for _, point := range points {
		fmt.Println(point)
	}
}

func TestQuadTree_QueryByPointType(t *testing.T) {
	qt := NewQuadTree(BBox{X: 200, Y: 200, W: 200, H: 200}, 15)
	for i := 1; i <= 10; i++ {
		qt.Insert(Point{X: rand.Float64() * qt.Boundary.W, Y: rand.Float64() * qt.Boundary.H, PointType: utils.FoodPointType})
	}
	for i := 1; i <= 5; i++ {
		qt.Insert(Point{X: rand.Float64() * qt.Boundary.W, Y: rand.Float64() * qt.Boundary.H, PointType: utils.PlayerHeadPointType})
	}
	qt.Print(true)

	var foodPoints []Point
	qt.QueryByPointType(map[string]bool{utils.FoodPointType: true}, &foodPoints)
	fmt.Println("Returned points", len(foodPoints))
	for _, point := range foodPoints {
		fmt.Println(point)
	}
}
