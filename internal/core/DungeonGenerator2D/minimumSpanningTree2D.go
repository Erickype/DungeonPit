package core

import "github.com/goki/mat32"

type IMinimumSpanningTree2D interface {
	CalculatePrimDistances()
	CalculateMST2D()
}

type MinimumSpanningTree2D struct {
	EdgesSet      []Line2D
	PrimDistances []float64
	MSTEdges      []Line2D
}

func (m *MinimumSpanningTree2D) CalculateMST2D() {
	m.MSTEdges = make([]Line2D, 0)
	openSet := make(map[mat32.Vec3]struct{})
	for _, edge := range m.EdgesSet {
		openSet[edge.A] = struct{}{}
		openSet[edge.B] = struct{}{}
	}
	closedSet := make(map[mat32.Vec3]struct{})
	closedSet[m.EdgesSet[0].A] = struct{}{}
	for {
		if len(openSet) <= 0 {
			break
		}
		chosen := false
		var chosenEdge Line2D
		minWeight := 100000000.0
		for i, edge := range m.EdgesSet {
			closedVertices := 0
			if _, exists := closedSet[edge.A]; !exists {
				closedVertices++
			}
			if _, exists := closedSet[edge.B]; !exists {
				closedVertices++
			}
			if (closedVertices == 1) && (m.PrimDistances[i] < minWeight) {
				chosenEdge = edge
				chosen = true
				minWeight = m.PrimDistances[i]
			}
		}
		if chosen {
			m.MSTEdges = append(m.MSTEdges, chosenEdge)
			delete(openSet, chosenEdge.A)
			delete(openSet, chosenEdge.B)
			closedSet[chosenEdge.A] = struct{}{}
			closedSet[chosenEdge.B] = struct{}{}
		}
	}
}

func (m *MinimumSpanningTree2D) CalculatePrimDistances() {
	m.PrimDistances = make([]float64, len(m.EdgesSet))
	for i, edge := range m.EdgesSet {
		m.PrimDistances[i] = edge.CalculateVertexDistance()
	}
}

func NewMinimumSpanningTree2D(delaunayEdgesSet []Line2D) *MinimumSpanningTree2D {
	return &MinimumSpanningTree2D{
		EdgesSet: delaunayEdgesSet,
	}
}
