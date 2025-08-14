package core

import (
	"math"

	"github.com/Erickype/DungeonPit/internal/core"
	"github.com/goki/mat32"
)

type IAStar2D interface {
	FindPath() [][]mat32.Vec2
	InitializeHelperData()
	IsInputDataValid() bool
	GetMinimumCostBetweenTwoTiles(tileA, tileB mat32.Vec2) int
}

type AStar2D struct {
	Start                      mat32.Vec2
	Target                     mat32.Vec2
	Grid                       map[mat32.Vec2]core.CellType
	PathFindingData            map[mat32.Vec2]PathFindingData
	DiscoveredTilesSortingCost []int
	DiscoveredTilesIndexes     []mat32.Vec2
	AnalysedTileIndexes        []mat32.Vec2
}

func (a *AStar2D) GetMinimumCostBetweenTwoTiles(tileA, tileB mat32.Vec2) int {
	sub := tileA.Sub(tileB)
	cost := int(math.Abs(float64(sub.X)) + math.Abs(float64(sub.Y)))
	switch a.Grid[tileB] {
	case core.CellTypeRoom:
		cost += 10
		break
	case core.CellTypeHallway:
		cost += 5
		break
	case core.CellTypeNone:
		cost += 1
		break
	}
	return cost
}

func (a *AStar2D) IsInputDataValid() bool {
	if !a.Start.IsEqual(a.Target) {
		return true
	}
	return false
}

func (a *AStar2D) InitializeHelperData() {
	a.PathFindingData = make(map[mat32.Vec2]PathFindingData)
	a.DiscoveredTilesSortingCost = make([]int, 0)
	a.DiscoveredTilesIndexes = make([]mat32.Vec2, 0)
	a.AnalysedTileIndexes = make([]mat32.Vec2, 0)
}

func (a *AStar2D) FindPath() [][]mat32.Vec2 {
	a.InitializeHelperData()
	if !a.IsInputDataValid() {
		return nil
	}
	NewPathFindingData(
		WithIndex(a.Start),
		WithCostToEnterTile(2),
		WithCostFromStart(0),
		WithMinimumCostToTarget(15),
	)
	return nil
}

func NewAStar2D(start, target mat32.Vec2, grid map[mat32.Vec2]core.CellType) *AStar2D {
	return &AStar2D{
		Start:  start,
		Target: target,
		Grid:   grid,
	}
}
