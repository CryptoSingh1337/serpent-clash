package system

import (
	"fmt"
	"math"
	"testing"
)

func TestGenerateSpawnPoints(t *testing.T) {
	points := GenerateSpawnPoints(12)
	for _, point := range points {
		fmt.Printf("X: %f, Y: %f\n", point.X, point.Y)
	}
	angle := math.Atan2(-200-0, -1300-1400)
	opposite := angle + math.Pi
	angle = math.Pi * 0.75
	opposite = angle + math.Pi
	if opposite > math.Pi {
		opposite -= 2 * math.Pi
	}
	fmt.Println("Angle:", angle*57.2958, "Opposite:", opposite*57.2958)
}
