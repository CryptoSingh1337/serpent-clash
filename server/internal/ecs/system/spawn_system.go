package system

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/component"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/storage"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/types"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/utils"
	"math"
	"math/rand/v2"
)

const (
	spawnRegionDistanceFromOrigin = utils.WorldBoundaryRadius * 0.7
	spawnRegionCount              = 12
)

var (
	spawnRegions []utils.Coordinate
)

type SpawnSystem struct {
	storage            storage.Storage
	spawnQueue         <-chan types.Id
	newId              func() types.Id
	lastSpawnRegionIdx int
}

func NewSpawnSystem(storage storage.Storage, spawnQueue <-chan types.Id, newId func() types.Id) *SpawnSystem {
	spawnRegions = GenerateSpawnPoints(spawnRegionCount)
	storage.AddSharedResource(utils.SpawnRegions, spawnRegions)
	return &SpawnSystem{
		storage,
		spawnQueue,
		newId,
		0,
	}
}

func (s *SpawnSystem) Update() {
	s.buildQuadTree()
	r := s.storage.GetSharedResource(utils.QuadTreeResource)
	if r == nil {
		return
	}
	qt := r.(*storage.QuadTree)
	for {
		select {
		case playerId := <-s.spawnQueue:
			utils.Logger.Info().Msgf("Spawning %v player id", playerId)
			var minDensityRegion utils.Coordinate
			previousRegionDensity := math.MaxInt
			var playerHeads []storage.Point
			regionIdx := (s.lastSpawnRegionIdx + 1) % spawnRegionCount
			s.lastSpawnRegionIdx = regionIdx
			utils.Logger.Info().Msgf("Spawn region: %v", regionIdx)
			for i := 0; i < spawnRegionCount; i++ {
				region := spawnRegions[regionIdx]
				regionIdx = (regionIdx + 1) % spawnRegionCount
				var p []storage.Point
				qt.QueryBCircleByPointType(storage.BCircle{
					X: region.X,
					Y: region.Y,
					R: utils.SpawnRegionRadius},
					map[string]bool{utils.PlayerHeadPointType: true},
					&p)
				if previousRegionDensity > len(p) {
					minDensityRegion = region
					playerHeads = p
				}
				previousRegionDensity = len(p)
			}
			utils.Logger.Debug().Msgf("Spawn region: %v", minDensityRegion)
			var segments []utils.Coordinate
			if playerHeads == nil {
				angle := rand.Float64() * 2 * math.Pi
				radius := utils.SpawnRegionRadius - 250*math.Sqrt(rand.Float64())
				segments = GenerateSnakeSegments(utils.Coordinate{
					X: minDensityRegion.X + radius*math.Cos(angle),
					Y: minDensityRegion.Y + radius*math.Sin(angle),
				}, utils.DefaultSnakeLength)
			} else {
				for _, point := range playerHeads {
					angle := math.Atan2(point.Y-minDensityRegion.Y, point.X-minDensityRegion.X)
					if angle > math.Pi {
						angle -= 2 * math.Pi
					}
					segments = GenerateSnakeSegments(utils.Coordinate{
						X: point.X + math.Cos(angle),
						Y: point.X + math.Sin(angle),
					}, utils.DefaultSnakeLength)
				}
			}
			utils.Logger.Info().Msgf("Head: %v", segments[0])
			c := s.storage.GetComponentByEntityIdAndName(playerId, utils.SnakeComponent)
			if c == nil {
				break
			}
			snakeComponent := c.(*component.Snake)
			snakeComponent.Segments = segments
		default:
			goto SpawnFood
		}
	}
SpawnFood:
	// TODO: spawn/de-spawn food and maintain food threshold
}

func (s *SpawnSystem) Stop() {

}

func (s *SpawnSystem) buildQuadTree() {
	playerEntities := s.storage.GetAllEntitiesByType(utils.PlayerEntity)
	qt := storage.NewQuadTree(storage.BBox{X: 0, Y: 0, W: utils.WorldWeight, H: utils.WorldHeight}, 15)
	for _, playerId := range playerEntities {
		comp := s.storage.GetComponentByEntityIdAndName(playerId, utils.SnakeComponent)
		if comp == nil {
			continue
		}
		snakeComponent := comp.(*component.Snake)
		if len(snakeComponent.Segments) == 0 {
			continue
		}
		head := snakeComponent.Segments[0]
		qt.Insert(storage.Point{X: head.X, Y: head.Y, EntityId: playerId, PointType: utils.PlayerHeadPointType})
		for i := 1; i < len(snakeComponent.Segments); i++ {
			segment := snakeComponent.Segments[i]
			qt.Insert(storage.Point{X: segment.X, Y: segment.Y, EntityId: playerId, PointType: utils.PlayerSegmentPointType})
		}
	}
	s.storage.AddSharedResource(utils.QuadTreeResource, qt)
}

func GenerateSpawnPoints(count int) []utils.Coordinate {
	var points []utils.Coordinate
	angleStep := 2 * math.Pi / float64(count)
	for i := 0; i < count; i++ {
		angle := angleStep * float64(i)
		x := int(spawnRegionDistanceFromOrigin * math.Cos(angle))
		y := int(spawnRegionDistanceFromOrigin * math.Sin(angle))
		points = append(points, utils.Coordinate{X: float64(x), Y: float64(y)})
	}
	return points
}

func GenerateSnakeSegments(c utils.Coordinate, length int) []utils.Coordinate {
	theta := math.Atan2(c.Y, c.X)
	x := c.X
	y := c.Y
	segments := make([]utils.Coordinate, length)
	segments[0] = utils.Coordinate{X: x, Y: y}
	for i := 1; i < len(segments); i++ {
		segments[i].X = x - float64(i)*utils.SnakeSegmentDistance*math.Cos(theta)
		segments[i].Y = y - float64(i)*utils.SnakeSegmentDistance*math.Sin(theta)
	}
	return segments
}
