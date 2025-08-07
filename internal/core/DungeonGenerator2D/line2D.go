package core

import "github.com/goki/mat32"

type Line2D struct {
	A mat32.Vec3
	B mat32.Vec3
}

func NewLine2D(a, b mat32.Vec3) *Line2D {
	return &Line2D{a, b}
}
