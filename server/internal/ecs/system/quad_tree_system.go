package system

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/component"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/storage"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/ecs/utils"
)

type QuadTreeSystem struct {
	storage storage.Storage
}

func NewQuadTreeSystem(storage storage.Storage) System {
	return &QuadTreeSystem{
		storage,
	}
}

func (q *QuadTreeSystem) Name() string {
	return utils.QuadTreeSystemName
}

func (q *QuadTreeSystem) Update() {
	playerEntities := q.storage.GetAllEntitiesByType(utils.PlayerEntity)
	qt := storage.NewQuadTree(storage.BBox{
		X: 0,
		Y: 0,
		W: utils.WorldWeight,
		H: utils.WorldHeight,
	}, utils.QuadTreeSegmentCapacity)
	for _, playerId := range playerEntities {
		comp := q.storage.GetComponentByEntityIdAndName(playerId, utils.SnakeComponent)
		if comp == nil {
			continue
		}
		snakeComponent := comp.(*component.Snake)
		if len(snakeComponent.Segments) == 0 {
			continue
		}
		head := snakeComponent.Segments[0]
		qt.Insert(storage.Point{
			X:         head.X,
			Y:         head.Y,
			EntityId:  playerId,
			PointType: utils.PlayerHeadPointType,
		})
		for i := 1; i < len(snakeComponent.Segments); i++ {
			segment := snakeComponent.Segments[i]
			qt.Insert(storage.Point{
				X:         segment.X,
				Y:         segment.Y,
				EntityId:  playerId,
				PointType: utils.PlayerSegmentPointType,
			})
		}
	}
	foodEntities := q.storage.GetAllEntitiesByType(utils.FoodEntity)
	for _, foodId := range foodEntities {
		comp := q.storage.GetComponentByEntityIdAndName(foodId, utils.PositionComponent)
		if comp == nil {
			continue
		}
		positionComponent := comp.(*component.Position)
		qt.Insert(storage.Point{
			X:         positionComponent.X,
			Y:         positionComponent.Y,
			EntityId:  foodId,
			PointType: utils.FoodPointType,
		})
	}
	q.storage.AddSharedResource(utils.QuadTreeResource, qt)
}

func (q *QuadTreeSystem) Stop() {
}
