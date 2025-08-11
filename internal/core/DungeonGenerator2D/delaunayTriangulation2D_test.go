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
	triangles := d.Triangles
	if triangles == nil {
		t.Error("Triangles should not be nil")
	}
	if len(triangles) != 1 {
		t.Error("Triangles should have length 1")
	}
}

func TestDelaunayTriangulation2D_GenerateEdgesSet(t *testing.T) {
	triangles := []mat32.Triangle{
		mat32.NewTriangle(
			mat32.NewVec3(0.5, 0.5, 0),
			mat32.NewVec3(2.5, 0.5, 0),
			mat32.NewVec3(2, 3, 0),
		),
	}
	delaunayTriangulation2D := DelaunayTriangulation2D{
		Triangles: triangles,
	}
	delaunayTriangulation2D.GenerateEdgesSet()
	if delaunayTriangulation2D.EdgesSet == nil {
		t.Error("Edges set should not be nil")
	}
	if len(delaunayTriangulation2D.EdgesSet) == 1 {
		t.Error("Edges set should have length 1")
	}
	t.Log(delaunayTriangulation2D.EdgesSet)
}
