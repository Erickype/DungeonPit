package core

import "github.com/goki/mat32"

type ILine2D interface {
	IsSameLine2D(testLine Line2D) bool
}

type Line2D struct {
	A mat32.Vec3
	B mat32.Vec3
}

func (l *Line2D) IsSameLine2D(testLine Line2D) bool {
	if l.A.X == testLine.A.X && l.A.Y == testLine.A.Y && l.B.X == testLine.B.X && l.B.Y == testLine.B.Y {
		return true
	}
	return false
}

func NewLine2D(a, b mat32.Vec3) *Line2D {
	return &Line2D{a, b}
}
