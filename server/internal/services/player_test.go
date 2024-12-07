package services

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/utils"
	"math"
	"testing"
)

func TestPlayer_GenerateRandomPosition(t *testing.T) {
	for i := 0; i < 1000; i++ {
		player := NewPlayer(nil)
		player.GenerateRandomPosition(5)
		p, q := player.Segments[0].X, player.Segments[0].Y
		distance := math.Sqrt(p*p + q*q)
		if i < 4 {
			t.Logf("Segments: %v", player.Segments)
		}
		if distance > utils.WorldBoundaryRadius {
			t.Fatalf("Interation: %v - Point is out of world radius, p, q: %v, %v, distance: %v", i, p,
				q, distance)
		}
	}
}
