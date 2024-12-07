package services

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/utils"
)

type SpatialHashGrid struct {
	CellSize float64
	Grid     map[int]map[int]map[string]*Player
}

type Collision struct {
	A string
	B string
}

func NewSpatialHashGrid(cellSize float64) *SpatialHashGrid {
	return &SpatialHashGrid{
		CellSize: cellSize,
		Grid:     make(map[int]map[int]map[string]*Player),
	}
}

func (s *SpatialHashGrid) GetCellCoords(x, y float64) (int, int) {
	return int(x / s.CellSize), int(y / s.CellSize)
}

func (s *SpatialHashGrid) InsertPlayer(player *Player) {
	for _, segment := range player.Segments {
		x, y := s.GetCellCoords(segment.X, segment.Y)
		if s.Grid[x] == nil {
			s.Grid[x] = make(map[int]map[string]*Player)
		}
		if s.Grid[x][y] == nil {
			s.Grid[x][y] = make(map[string]*Player)
		}
		s.Grid[x][y][player.Id] = player
	}
}

func (s *SpatialHashGrid) RemovePlayer(player *Player) {
	for _, segment := range player.Segments {
		x, y := s.GetCellCoords(segment.X, segment.Y)
		if s.Grid[x] != nil && s.Grid[x][y] != nil {
			delete(s.Grid[x][y], player.Id)
		}
	}
}

func (s *SpatialHashGrid) RemovePlayerById(id string) {
	for i, row := range s.Grid {
		for j, cell := range row {
			for idKey, _ := range cell {
				if idKey == id {
					utils.Logger.LogInfo().Msgf("Removed player - %v from hash grid", id)
					delete(s.Grid[i][j], id)
					delete(s.Grid[i], j)
					delete(s.Grid, i)
				}
			}
		}
	}
}

func (s *SpatialHashGrid) DetectCollisions() []Collision {
	collisions := make([]Collision, 0)
	checked := make(map[string]bool)

	for _, row := range s.Grid {
		for _, cell := range row {
			for idA, playerA := range cell {
				for idB, playerB := range cell {
					if idA != idB && !checked[idA+idB] && !checked[idB+idA] {
						if detectPlayerCollision(playerA, playerB) {
							collisions = append(collisions, Collision{A: idA, B: idB})
						}
						checked[idA+idB] = true
					}
				}
			}
		}
	}
	return collisions
}

func detectPlayerCollision(playerA, playerB *Player) bool {
	for _, segA := range playerA.Segments {
		for _, segB := range playerB.Segments {
			dx := segA.X - segB.X
			dy := segA.Y - segA.Y
			distance := dx*dx + dy*dy
			if distance < utils.SnakeSegmentDistance*utils.SnakeSegmentDistance {
				return true
			}
		}
	}
	return false
}
