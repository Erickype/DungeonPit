package core

import (
	"github.com/goki/mat32"
	"testing"
)

func TestMinimumSpanningTree2D_CalculateMST2D(t *testing.T) {
	vertexes := []mat32.Vec2{
		{0.5, 0.5},
		{2.5, 0.5},
		{2, 3},
	}
	d := NewDelaunayTriangulation2D(vertexes)
	d.Calculate()
	d.GenerateEdgesSet()
	m := NewMinimumSpanningTree2D(d.EdgesSet)
	m.CalculatePrimDistances()
	m.CalculateMST2D()
	if len(m.PrimDistances) == 0 {
		t.Errorf("Prim distances length is zero")
	}
	if len(m.MSTEdges) == 0 {
		t.Errorf("MST edges length is zero")
	}
	if len(m.MSTEdges) != 2 {
		t.Errorf("MST should be 2, actual value: %d", len(m.MSTEdges))
	}
}
