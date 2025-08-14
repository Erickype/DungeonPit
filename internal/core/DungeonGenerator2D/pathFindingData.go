package core

import "github.com/goki/mat32"

type PathFindingData struct {
	Index               mat32.Vec2
	CostToEnterTile     int
	CostFromStart       int
	MinimumCostToTarget int
	PreviousIndex       mat32.Vec2
}

type PathFindingOption func(*PathFindingData)

func WithIndex(v mat32.Vec2) PathFindingOption {
	return func(p *PathFindingData) { p.Index = v }
}

func WithCostToEnterTile(v int) PathFindingOption {
	return func(p *PathFindingData) { p.CostToEnterTile = v }
}

func WithCostFromStart(v int) PathFindingOption {
	return func(p *PathFindingData) { p.CostFromStart = v }
}

func WithMinimumCostToTarget(v int) PathFindingOption {
	return func(p *PathFindingData) { p.MinimumCostToTarget = v }
}

func WithPreviousIndex(v mat32.Vec2) PathFindingOption {
	return func(p *PathFindingData) { p.PreviousIndex = v }
}

func NewPathFindingData(opts ...PathFindingOption) *PathFindingData {
	// Default values
	p := &PathFindingData{
		Index:               mat32.NewVec2(-999, 999),
		CostToEnterTile:     1,
		CostFromStart:       999999,
		MinimumCostToTarget: 999999,
		PreviousIndex:       mat32.NewVec2(-999, 999),
	}
	// Apply options
	for _, opt := range opts {
		opt(p)
	}
	return p
}
