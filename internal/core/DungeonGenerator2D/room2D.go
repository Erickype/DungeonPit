package core

import (
	"github.com/goki/mat32"
	"math"
)

type IRoom2D interface {
	Intersect(Room2D) bool
	AllPositionsWithin() []mat32.Vec2
}

type Room2D struct {
	Position mat32.Vec2
	Size     mat32.Vec2
}

func (r *Room2D) AllPositionsWithin() (cords []mat32.Vec2) {
	iLimit := math.Trunc(float64(r.Size.X - 1))
	jLimit := math.Trunc(float64(r.Size.Y - 1))
	for i := 0; i <= int(iLimit); i++ {
		for j := 0; j <= int(jLimit); j++ {
			cord := mat32.Vec2{
				X: float32(i) + r.Position.X,
				Y: float32(j) + r.Position.Y,
			}
			cords = append(cords, cord)
		}
	}
	return cords
}

func (r *Room2D) Intersect(roomB Room2D) bool {
	AMinX := r.Position.X
	AMaxX := r.Position.X + r.Size.X
	AMinY := r.Position.Y
	AMaxY := r.Position.Y + r.Size.Y

	BMinX := roomB.Position.X
	BMaxX := roomB.Position.X + roomB.Size.X
	BMinY := roomB.Position.Y
	BMaxY := roomB.Position.Y + roomB.Size.Y

	return !(((AMinX >= BMaxX) || (AMaxX <= BMinX)) || ((AMinY >= BMaxY) || (AMaxY <= BMinY)))
}
