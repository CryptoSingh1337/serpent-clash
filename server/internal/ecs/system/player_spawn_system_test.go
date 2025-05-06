package system

import (
	"fmt"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/utils"
	"math"
	"math/rand/v2"
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

func TestGenerateSnakeSegments(t *testing.T) {
	x := float64(1662)
	y := float64(962)
	angle := rand.Float64() * 2 * math.Pi
	radius := utils.SpawnRegionRadius - 250*math.Sqrt(rand.Float64())
	segments := GenerateSnakeSegments(utils.Coordinate{
		X: x + radius*math.Cos(angle),
		Y: y + radius*math.Sin(angle),
	}, utils.DefaultSnakeLength)
	fmt.Printf("Head: %v", segments[0])
}
