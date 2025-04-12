package system

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/component"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/storage"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/types"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/utils"
)

type SpawnSystem struct {
	storage storage.Storage
	newId   func() types.Id
}

func NewSpawnSystem(storage storage.Storage, newId func() types.Id) *SpawnSystem {
	return &SpawnSystem{
		storage,
		newId,
	}
}

func (s *SpawnSystem) Update() {
	s.buildQuadTree()
	r := s.storage.GetSharedResource(utils.QuadTreeResource)
	if r == nil {
		return
	}
	_ = r.(*storage.QuadTree)
	// TODO: spawn new players
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
		head := snakeComponent.Segments[0]
		qt.Insert(storage.Point{X: head.X, Y: head.Y, EntityId: playerId, PointType: "head"})
		for i := 1; i < len(snakeComponent.Segments); i++ {
			segment := snakeComponent.Segments[i]
			qt.Insert(storage.Point{X: segment.X, Y: segment.Y, EntityId: playerId, PointType: "segment"})
		}
	}
	s.storage.AddSharedResource(utils.QuadTreeResource, qt)
}
