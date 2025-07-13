package system

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/component"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/storage"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/types"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/utils"
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

type PlayerSpawnSystem struct {
	storage            storage.Storage
	spawnQueue         <-chan *types.JoinEvent
	lastSpawnRegionIdx int
	JoinEventQueue     []*types.JoinEvent
}

func NewSpawnSystem(storage storage.Storage, spawnQueue <-chan *types.JoinEvent) System {
	spawnRegions = GenerateSpawnPoints(spawnRegionCount)
	storage.AddSharedResource(utils.SpawnRegions, spawnRegions)
	return &PlayerSpawnSystem{
		storage:            storage,
		spawnQueue:         spawnQueue,
		lastSpawnRegionIdx: 0,
	}
}

func (s *PlayerSpawnSystem) Name() string {
	return utils.PlayerSpawnSystemName
}

func (s *PlayerSpawnSystem) Update() {
	if len(s.JoinEventQueue) > 0 {
		s.JoinEventQueue = nil
	}
	r := s.storage.GetSharedResource(utils.QuadTreeResource)
	if r == nil {
		return
	}
	qt := r.(*storage.QuadTree)
	for {
		select {
		case joinEvent := <-s.spawnQueue:
			playerId := joinEvent.EntityId
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
			s.JoinEventQueue = append(s.JoinEventQueue, joinEvent)
		default:
			goto escape
		}
	}
escape:
}

func (s *PlayerSpawnSystem) Stop() {

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
