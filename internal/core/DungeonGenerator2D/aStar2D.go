package core

import (
	"math"
	"slices"

	"github.com/Erickype/DungeonPit/internal/core"
	"github.com/goki/mat32"
)

type IAStar2D interface {
	FindPath() [][]mat32.Vec2
	InitializeHelperData()
	IsInputDataValid() bool
	GetMinimumCostBetweenTwoTiles(tileA, tileB mat32.Vec2) int
	DiscoverTile(pathData PathFindingData)
	InsertTileInDiscoveredArray(pathData PathFindingData)
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

func (a *AStar2D) InsertTileInDiscoveredArray(pathData PathFindingData) {
	sortingCost := pathData.CostFromStart + pathData.MinimumCostToTarget
	if len(a.DiscoveredTilesIndexes) == 0 {
		a.DiscoveredTilesSortingCost = append(a.DiscoveredTilesSortingCost, sortingCost)
		a.DiscoveredTilesIndexes = append(a.DiscoveredTilesIndexes, pathData.Index)
	}
	index := a.DiscoveredTilesSortingCost[len(a.DiscoveredTilesSortingCost)-1]
	if sortingCost >= index {
		a.DiscoveredTilesSortingCost = append(a.DiscoveredTilesSortingCost, sortingCost)
		a.DiscoveredTilesIndexes = append(a.DiscoveredTilesIndexes, pathData.Index)
	}
	for i, cost := range a.DiscoveredTilesSortingCost {
		if cost >= sortingCost {
			a.DiscoveredTilesSortingCost = slices.Insert(a.DiscoveredTilesSortingCost, i, cost)
			a.DiscoveredTilesIndexes = slices.Insert(a.DiscoveredTilesIndexes, i, pathData.Index)
			break
		}
	}
}

func (a *AStar2D) DiscoverTile(pathData PathFindingData) {
	a.PathFindingData[pathData.Index] = pathData
	a.InsertTileInDiscoveredArray(pathData)
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
	pathData := NewPathFindingData(
		WithIndex(a.Start),
		WithCostToEnterTile(1),
		WithCostFromStart(0),
		WithMinimumCostToTarget(a.GetMinimumCostBetweenTwoTiles(a.Start, a.Target)),
	)
	a.DiscoverTile(*pathData)

	for {
		if len(a.DiscoveredTilesIndexes) <= 0 {
			break
		}
	}
	return nil
}

func NewAStar2D(start, target mat32.Vec2, grid map[mat32.Vec2]core.CellType) *AStar2D {
	return &AStar2D{
		Start:  start,
		Target: target,
		Grid:   grid,
	}
}
