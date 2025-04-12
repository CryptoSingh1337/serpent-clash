package storage

import (
	"fmt"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/types"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/utils"
	"math"
)

const (
	maxDepth = 5
)

type Point struct {
	X, Y      float64
	EntityId  types.Id
	PointType string
}

type BBox struct {
	X, Y, W, H float64
}

func (r BBox) Contains(p Point) bool {
	return p.X >= r.X-r.W && p.X <= r.X+r.W &&
		p.Y >= r.Y-r.H && p.Y <= r.Y+r.H
}

func (r BBox) IntersectBBox(other BBox) bool {
	return !(other.X-other.W > r.X+r.W ||
		other.X+other.W < r.X-r.W ||
		other.Y-other.H > r.Y+r.H ||
		other.Y+other.H < r.Y-r.H)
}

func (r BBox) IntersectBCircle(c BCircle) bool {
	dx := math.Abs(r.X - c.X)
	dy := math.Abs(r.Y - c.Y)
	if dx > r.W+c.R || dy > r.H+c.R {
		return false
	}
	if dx <= r.W || dy <= r.H {
		return true
	}
	cornerDx := dx - r.W
	cornerDy := dy - r.H
	return (cornerDx*cornerDx + cornerDy*cornerDy) <= c.R*c.R
}

type BCircle struct {
	X, Y, R float64
}

func (c BCircle) Contains(p Point) bool {
	return utils.EuclideanDistance(c.X, c.Y, p.X, p.Y) <= c.R
}

type QuadTree struct {
	Boundary       BBox
	Capacity       int
	Points         []Point
	Divided        bool
	Depth          int
	NW, NE, SW, SE *QuadTree
}

func NewQuadTree(boundary BBox, capacity int) *QuadTree {
	return &QuadTree{Boundary: boundary, Capacity: capacity, Points: make([]Point, 0, capacity)}
}

func (qt *QuadTree) Insert(p Point) bool {
	if !qt.Boundary.Contains(p) {
		return false
	}
	if len(qt.Points) < qt.Capacity || qt.Depth >= maxDepth {
		qt.Points = append(qt.Points, p)
		return true
	}
	if !qt.Divided {
		qt.subDivide()
	}
	return qt.NW.Insert(p) || qt.NE.Insert(p) || qt.SW.Insert(p) || qt.SE.Insert(p)
}

func (qt *QuadTree) QueryBBox(rangeBBox BBox, found *[]Point) {
	if !qt.Boundary.IntersectBBox(rangeBBox) {
		return
	}
	for _, p := range qt.Points {
		if rangeBBox.Contains(p) {
			*found = append(*found, p)
		}
	}
	if qt.Divided {
		qt.NW.QueryBBox(rangeBBox, found)
		qt.NE.QueryBBox(rangeBBox, found)
		qt.SW.QueryBBox(rangeBBox, found)
		qt.SE.QueryBBox(rangeBBox, found)
	}
}

func (qt *QuadTree) QueryBCircle(rangeBCircle BCircle, found *[]Point) {
	if !qt.Boundary.IntersectBCircle(rangeBCircle) {
		return
	}
	for _, p := range qt.Points {
		if rangeBCircle.Contains(p) {
			*found = append(*found, p)
		}
	}
	if qt.Divided {
		qt.NW.QueryBCircle(rangeBCircle, found)
		qt.NE.QueryBCircle(rangeBCircle, found)
		qt.SW.QueryBCircle(rangeBCircle, found)
		qt.SE.QueryBCircle(rangeBCircle, found)
	}
}

func (qt *QuadTree) subDivide() {
	x := qt.Boundary.X
	y := qt.Boundary.Y
	w := qt.Boundary.W
	h := qt.Boundary.H
	qt.NW = NewQuadTree(BBox{X: x - w/2, Y: y + h/2, W: w / 2, H: h / 2}, qt.Capacity)
	qt.NE = NewQuadTree(BBox{X: x + w/2, Y: y + h/2, W: w / 2, H: h / 2}, qt.Capacity)
	qt.SW = NewQuadTree(BBox{X: x - w/2, Y: y - h/2, W: w / 2, H: h / 2}, qt.Capacity)
	qt.SE = NewQuadTree(BBox{X: x + w/2, Y: y - h/2, W: w / 2, H: h / 2}, qt.Capacity)
	qt.NW.Depth = qt.Depth + 1
	qt.NE.Depth = qt.Depth + 1
	qt.SW.Depth = qt.Depth + 1
	qt.SE.Depth = qt.Depth + 1
	qt.Divided = true
}

func (qt *QuadTree) Print(test bool) {
	var queue []*QuadTree
	queue = append(queue, qt)
	nodeId := 1
	for len(queue) > 0 {
		q := queue[0]
		msg := fmt.Sprintf("QT: %d, level: %d, Boundary: %v, points: %v\n", nodeId, q.Depth, q.Boundary, q.Points)
		if test {
			fmt.Printf(msg)
		} else {
			utils.Logger.Info().Msgf(msg)
		}
		if q.Divided {
			queue = append(queue, q.NW)
			queue = append(queue, q.NE)
			queue = append(queue, q.SW)
			queue = append(queue, q.SE)
		}
		queue = queue[1:]
		nodeId++
	}
}
