package core

import (
	"github.com/goki/mat32"
	"math"
	"slices"
)

type IDelaunayTriangulation interface {
	Calculate()
	FindSuperTriangle() mat32.Triangle
	AddVertex(vertex mat32.Vec2, triangles []mat32.Triangle) []mat32.Triangle
	IsInCircumcircle(vertex mat32.Vec2, triangle mat32.Triangle) bool
}
type DelaunayTriangulation2D struct {
	Vertexes  []mat32.Vec2
	Triangles []mat32.Triangle
}

func (d *DelaunayTriangulation2D) FindSuperTriangle() mat32.Triangle {
	minX := 1000000.0
	minY := 1000000.0
	maxX := -1000000.0
	maxY := -1000000.0
	for _, center := range d.Vertexes {
		x := float64(center.X)
		y := float64(center.Y)
		minX = math.Min(minX, x)
		minY = math.Min(minY, y)
		maxX = math.Max(maxX, x)
		maxY = math.Max(maxY, y)
	}
	dx := maxX - minX
	dy := maxY - minY
	deltaMax := math.Max(dx, dy) * 2
	a := mat32.Vec3{
		X: float32(minX - 1),
		Y: float32(minY - 1),
	}
	b := mat32.Vec3{
		X: float32(minX - 1),
		Y: float32(maxY + deltaMax),
	}
	c := mat32.Vec3{
		X: float32(maxX + deltaMax),
		Y: float32(minY - 1),
	}
	return mat32.NewTriangle(a, b, c)
}

func (d *DelaunayTriangulation2D) IsInCircumcircle(vertex mat32.Vec2, triangle mat32.Triangle) bool {
	ax := float64(triangle.A.X)
	ay := float64(triangle.A.Y)
	bx := float64(triangle.B.X)
	by := float64(triangle.B.Y)
	cx := float64(triangle.C.X)
	cy := float64(triangle.C.Y)
	px := float64(vertex.X)
	py := float64(vertex.Y)

	// 1 — Triangle orientation (signed area * 2)
	orient := (bx-ax)*(cy-ay) - (by-ay)*(cx-ax)

	// 2 — Translate so vertex is origin
	ax -= px
	ay -= py
	bx -= px
	by -= py
	cx -= px
	cy -= py

	// 3 — Squared lengths
	aSq := ax*ax + ay*ay
	bSq := bx*bx + by*by
	cSq := cx*cx + cy*cy

	// 4 — Determinant
	det := ax*(by*cSq-bSq*cy) -
		ay*(bx*cSq-bSq*cx) +
		aSq*(bx*cy-by*cx)

	// 5 — Interpret sign based on orientation
	const EPS = 1e-9
	if math.Abs(det) < EPS {
		// On the circumcircle
		return true
	}
	if orient > 0 {
		// CCW triangle → inside if det > 0
		return det > 0
	} else {
		// CW triangle → inside if det < 0
		return det < 0
	}
}

func (d *DelaunayTriangulation2D) AddVertex(vertex mat32.Vec2, triangles []mat32.Triangle) (newTriangles []mat32.Triangle) {
	newTriangles = []mat32.Triangle{}
	newEdges := make([]Line2D, 0)
	for _, triangle := range triangles {
		if d.IsInCircumcircle(vertex, triangle) {
			ab := NewLine2D(triangle.A, triangle.B)
			bc := NewLine2D(triangle.B, triangle.C)
			ac := NewLine2D(triangle.A, triangle.C)
			newEdges = append(newEdges, *ab, *bc, *ac)
		} else {
			newTriangles = append(newTriangles, triangle)
		}
	}
	newEdges2 := make([]Line2D, 0)
	notFUniqueEdges := make([]Line2D, 0)
	for _, edge := range newEdges {
		if slices.Contains(newEdges2, edge) || slices.Contains(newEdges2, *NewLine2D(edge.B, edge.A)) {
			notFUniqueEdges = append(notFUniqueEdges, edge)
			index := slices.Index(newEdges2, *NewLine2D(edge.B, edge.A))
			if index != -1 {
				newEdges2 = slices.Delete(newEdges2, index, index+1)
			}
			index = slices.Index(newEdges2, edge)
			if index != -1 {
				newEdges2 = slices.Delete(newEdges2, index, index+1)
			}
		} else if !(slices.Contains(notFUniqueEdges, edge) || slices.Contains(notFUniqueEdges, *NewLine2D(edge.B, edge.A))) {
			newEdges2 = append(newEdges2, edge)
		}
	}
	for _, edge := range newEdges2 {
		triangle := mat32.NewTriangle(edge.A, edge.B, mat32.NewVec3(vertex.X, vertex.Y, 0))
		newTriangles = append(newTriangles, triangle)
	}
	return newTriangles
}

func (d *DelaunayTriangulation2D) Calculate() {
	superTriangle := d.FindSuperTriangle()
	triangles := []mat32.Triangle{superTriangle}
	newTriangles := make([]mat32.Triangle, 0)
	for _, vertex := range d.Vertexes {
		triangles = d.AddVertex(vertex, triangles)
	}
	for _, triangle := range triangles {
		c1 := triangle.A.IsEqual(superTriangle.A)
		c2 := triangle.A.IsEqual(superTriangle.B)
		c3 := triangle.A.IsEqual(superTriangle.C)
		c4 := triangle.B.IsEqual(superTriangle.A)
		c5 := triangle.B.IsEqual(superTriangle.B)
		c6 := triangle.B.IsEqual(superTriangle.C)
		c7 := triangle.C.IsEqual(superTriangle.A)
		c8 := triangle.C.IsEqual(superTriangle.B)
		c9 := triangle.C.IsEqual(superTriangle.C)
		if !(c1 || c2 || c3 || c4 || c5 || c6 || c7 || c8 || c9) {
			newTriangles = append(newTriangles, triangle)
		}
	}
	d.Triangles = newTriangles
}

func NewDelaunayTriangulation2D(vertexes []mat32.Vec2) *DelaunayTriangulation2D {
	return &DelaunayTriangulation2D{Vertexes: vertexes}
}
