package system

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/component"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/storage"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/types"
	gameutils "github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/utils"
	"math"
	"math/rand/v2"
)

const (
	spawnRegionDistanceFromOrigin = gameutils.WorldBoundaryRadius * 0.7
	spawnRegionCount              = 12
)

var (
	spawnRegions []gameutils.Coordinate
)

type SpawnSystem struct {
	storage            storage.Storage
	spawnQueue         <-chan types.Id
	newId              func() types.Id
	lastSpawnRegionIdx int
}

func NewSpawnSystem(storage storage.Storage, spawnQueue <-chan types.Id, newId func() types.Id) *SpawnSystem {
	spawnRegions = GenerateSpawnPoints(spawnRegionCount)
	storage.AddSharedResource(gameutils.SpawnRegions, spawnRegions)
	return &SpawnSystem{
		storage,
		spawnQueue,
		newId,
		0,
	}
}

func (s *SpawnSystem) Update() {
	s.buildQuadTree()
	r := s.storage.GetSharedResource(gameutils.QuadTreeResource)
	if r == nil {
		return
	}
	qt := r.(*storage.QuadTree)
	for {
		select {
		case playerId := <-s.spawnQueue:
			gameutils.Logger.Info().Msgf("Spawning %v player id", playerId)
			var minDensityRegion gameutils.Coordinate
			previousRegionDensity := math.MaxInt
			var playerHeads []storage.Point
			regionIdx := (s.lastSpawnRegionIdx + 1) % spawnRegionCount
			s.lastSpawnRegionIdx = regionIdx
			gameutils.Logger.Info().Msgf("Spawn region: %v", regionIdx)
			for i := 0; i < spawnRegionCount; i++ {
				region := spawnRegions[regionIdx]
				regionIdx = (regionIdx + 1) % spawnRegionCount
				var p []storage.Point
				qt.QueryBCircleByPointType(storage.BCircle{
					X: region.X,
					Y: region.Y,
					R: gameutils.SpawnRegionRadius},
					map[string]bool{gameutils.PlayerHeadPointType: true},
					&p)
				if previousRegionDensity > len(p) {
					minDensityRegion = region
					playerHeads = p
				}
				previousRegionDensity = len(p)
			}
			gameutils.Logger.Debug().Msgf("Spawn region: %v", minDensityRegion)
			var segments []gameutils.Coordinate
			if playerHeads == nil {
				angle := rand.Float64() * 2 * math.Pi
				radius := gameutils.SpawnRegionRadius - 250*math.Sqrt(rand.Float64())
				segments = GenerateSnakeSegments(gameutils.Coordinate{
					X: minDensityRegion.X + radius*math.Cos(angle),
					Y: minDensityRegion.Y + radius*math.Sin(angle),
				}, gameutils.DefaultSnakeLength)
			} else {
				for _, point := range playerHeads {
					angle := math.Atan2(point.Y-minDensityRegion.Y, point.X-minDensityRegion.X)
					if angle > math.Pi {
						angle -= 2 * math.Pi
					}
					segments = GenerateSnakeSegments(gameutils.Coordinate{
						X: point.X + math.Cos(angle),
						Y: point.X + math.Sin(angle),
					}, gameutils.DefaultSnakeLength)
				}
			}
			gameutils.Logger.Info().Msgf("Head: %v", segments[0])
			c := s.storage.GetComponentByEntityIdAndName(playerId, gameutils.SnakeComponent)
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
	playerEntities := s.storage.GetAllEntitiesByType(gameutils.PlayerEntity)
	qt := storage.NewQuadTree(storage.BBox{X: 0, Y: 0, W: gameutils.WorldWeight, H: gameutils.WorldHeight}, 15)
	for _, playerId := range playerEntities {
		comp := s.storage.GetComponentByEntityIdAndName(playerId, gameutils.SnakeComponent)
		if comp == nil {
			continue
		}
		snakeComponent := comp.(*component.Snake)
		if len(snakeComponent.Segments) == 0 {
			continue
		}
		head := snakeComponent.Segments[0]
		qt.Insert(storage.Point{X: head.X, Y: head.Y, EntityId: playerId, PointType: gameutils.PlayerHeadPointType})
		for i := 1; i < len(snakeComponent.Segments); i++ {
			segment := snakeComponent.Segments[i]
			qt.Insert(storage.Point{X: segment.X, Y: segment.Y, EntityId: playerId, PointType: gameutils.PlayerSegmentPointType})
		}
	}
	s.storage.AddSharedResource(gameutils.QuadTreeResource, qt)
}

func GenerateSpawnPoints(count int) []gameutils.Coordinate {
	var points []gameutils.Coordinate
	angleStep := 2 * math.Pi / float64(count)
	for i := 0; i < count; i++ {
		angle := angleStep * float64(i)
		x := int(spawnRegionDistanceFromOrigin * math.Cos(angle))
		y := int(spawnRegionDistanceFromOrigin * math.Sin(angle))
		points = append(points, gameutils.Coordinate{X: float64(x), Y: float64(y)})
	}
	return points
}

func GenerateSnakeSegments(c gameutils.Coordinate, length int) []gameutils.Coordinate {
	theta := math.Atan2(c.Y, c.X)
	x := c.X
	y := c.Y
	segments := make([]gameutils.Coordinate, length)
	segments[0] = gameutils.Coordinate{X: x, Y: y}
	for i := 1; i < len(segments); i++ {
		segments[i].X = x - float64(i)*gameutils.SnakeSegmentDistance*math.Cos(theta)
		segments[i].Y = y - float64(i)*gameutils.SnakeSegmentDistance*math.Sin(theta)
	}
	return segments
}
