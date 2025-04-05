package ecs

import (
	"github.com/CryptoSingh1337/serpent-clash/server/internal/types"
)

type Point struct {
	X, Y     float64
	EntityId types.Id
}

type Rect struct {
	X, Y, W, H float64
}

func (r Rect) Contains(p Point) bool {
	return p.X >= r.X-r.W && p.X <= r.X+r.W &&
		p.Y >= r.Y-r.H && p.Y <= r.Y+r.H
}

type QuadTree struct {
	Boundary       Rect
	Capacity       int
	Points         []Point
	Divided        bool
	NW, NE, SW, SE *QuadTree
}

func NewQuadTree(boundary Rect, capacity int) *QuadTree {
	return &QuadTree{Boundary: boundary, Capacity: capacity, Points: make([]Point, 0, capacity)}
}

func (qt *QuadTree) Insert(p Point) bool {
	if !qt.Boundary.Contains(p) {
		return false
	}
	if len(qt.Points) < qt.Capacity {
		qt.Points = append(qt.Points, p)
		return true
	}
	if !qt.Divided {
		qt.subDivide()
	}
	return qt.NW.Insert(p) || qt.NE.Insert(p) || qt.SW.Insert(p) || qt.SE.Insert(p)
}

func (qt *QuadTree) subDivide() {
	x := qt.Boundary.X
	y := qt.Boundary.Y
	w := qt.Boundary.W
	h := qt.Boundary.H
	qt.NW = NewQuadTree(Rect{X: x - w/2, Y: y + h/2, W: w / 2, H: h / 2}, qt.Capacity)
	qt.NE = NewQuadTree(Rect{X: x + w/2, Y: y + h/2, W: w / 2, H: h / 2}, qt.Capacity)
	qt.SW = NewQuadTree(Rect{X: x - w/2, Y: y - h/2, W: w / 2, H: h / 2}, qt.Capacity)
	qt.SE = NewQuadTree(Rect{X: x + w/2, Y: y - h/2, W: w / 2, H: h / 2}, qt.Capacity)
	qt.Divided = true
}
