package core

import (
	"github.com/goki/mat32"
	"testing"
)

func TestDelaunayTriangulation2D_Calculate(t *testing.T) {
	vertexes := []mat32.Vec2{
		{0.5, 0.5},
		{2.5, 0.5},
		{2, 3},
	}
	d := NewDelaunayTriangulation2D(vertexes)
	d.Calculate()
	if d.Triangles == nil {
		t.Error("Triangles should not be nil")
	}
}
