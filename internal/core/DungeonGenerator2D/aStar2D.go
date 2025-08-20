package core

import (
	"math"
	"slices"

	"github.com/Erickype/DungeonPit/internal/core"
	"github.com/goki/mat32"
)

type IAStar2D interface {
	FindPath() []mat32.Vec2
	InitializeHelperData()
	IsInputDataValid() bool
	GetMinimumCostBetweenTwoTiles(tileA, tileB mat32.Vec2) int
	DiscoverTile(pathData PathFindingData)
	InsertTileInDiscoveredArray(pathData PathFindingData)
	AnalyzeNextDiscoveredTile() bool
	GeneratePath() []mat32.Vec2
	PullCheapestTileOutOfDiscoveredList() PathFindingData
	GetValidTileNeighbours(tile mat32.Vec2) []PathFindingData
	DiscoverNextNeighbour() bool
}

type AStar2D struct {
	Start                      mat32.Vec2
	Target                     mat32.Vec2
	Grid                       map[mat32.Vec2]core.CellType
	PathFindingData            map[mat32.Vec2]PathFindingData
	DiscoveredTilesSortingCost []int
	DiscoveredTilesIndexes     []mat32.Vec2
	AnalysedTileIndexes        []mat32.Vec2
	CurrentDiscoveredTile      PathFindingData
	CurrentNeighbours          []PathFindingData
	CurrentNeighbour           PathFindingData
}

func (a *AStar2D) DiscoverNextNeighbour() bool {
	a.CurrentNeighbour = a.CurrentNeighbours[0]
	a.CurrentNeighbours = slices.Delete(a.CurrentNeighbours, 0, 1)
	if !slices.Contains(a.AnalysedTileIndexes, a.CurrentNeighbour.Index) {
		costFromStart := a.CurrentDiscoveredTile.CostFromStart + a.CurrentNeighbour.CostToEnterTile
		indexInDiscovered := slices.Index(a.DiscoveredTilesIndexes, a.CurrentNeighbour.Index)
		if !(indexInDiscovered == -1) {
			a.CurrentNeighbour = a.PathFindingData[a.CurrentNeighbour.Index]
			if !(costFromStart < a.CurrentNeighbour.CostFromStart) {
				return false
			}
			a.DiscoveredTilesIndexes = slices.Delete(a.DiscoveredTilesIndexes, indexInDiscovered, 1)
			a.DiscoveredTilesSortingCost = slices.Delete(a.DiscoveredTilesSortingCost, indexInDiscovered, 1)
		}
		pathData := NewPathFindingData(
			WithIndex(a.CurrentNeighbour.Index),
			WithCostToEnterTile(a.CurrentNeighbour.CostToEnterTile),
			WithCostFromStart(costFromStart),
			WithMinimumCostToTarget(a.GetMinimumCostBetweenTwoTiles(a.CurrentNeighbour.Index, a.Target)),
			WithPreviousIndex(a.CurrentDiscoveredTile.Index),
		)
		a.DiscoverTile(*pathData)
		return a.CurrentNeighbour.Index == a.Target
	}
	return false
}

func (a *AStar2D) GetValidTileNeighbours(tile mat32.Vec2) []PathFindingData {
	neighbours := make([]PathFindingData, 4)
	neighbours[0] = *NewPathFindingData(WithIndex(tile.Sub(mat32.NewVec2(0, 1))))
	neighbours[1] = *NewPathFindingData(WithIndex(tile.Add(mat32.NewVec2(1, 0))))
	neighbours[2] = *NewPathFindingData(WithIndex(tile.Add(mat32.NewVec2(0, 1))))
	neighbours[3] = *NewPathFindingData(WithIndex(tile.Sub(mat32.NewVec2(1, 0))))
	return neighbours
}

func (a *AStar2D) PullCheapestTileOutOfDiscoveredList() PathFindingData {
	tileIndex := a.DiscoveredTilesIndexes[0]
	a.DiscoveredTilesSortingCost = slices.Delete(a.DiscoveredTilesSortingCost, 0, 1)
	a.DiscoveredTilesIndexes = slices.Delete(a.DiscoveredTilesIndexes, 0, 1)
	a.AnalysedTileIndexes = append(a.AnalysedTileIndexes, tileIndex)
	return a.PathFindingData[tileIndex]
}

func (a *AStar2D) GeneratePath() []mat32.Vec2 {
	current := a.Target
	invertedPath := make([]mat32.Vec2, 0)
	for {
		if current == a.Start {
			break
		}
		invertedPath = append(invertedPath, current)
		current = a.PathFindingData[current].PreviousIndex
	}
	reversed := make([]mat32.Vec2, 0)
	for i := len(invertedPath) - 1; i >= 0; i-- {
		reversed = append(reversed, invertedPath[i])
	}
	return reversed
}

func (a *AStar2D) AnalyzeNextDiscoveredTile() bool {
	a.CurrentDiscoveredTile = a.PullCheapestTileOutOfDiscoveredList()
	a.CurrentNeighbours = a.GetValidTileNeighbours(a.CurrentDiscoveredTile.Index)
	for {
		if len(a.CurrentNeighbours) <= 0 {
			break
		}
		if a.DiscoverNextNeighbour() {
			return true
		}
	}
	return false
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

func (a *AStar2D) FindPath() []mat32.Vec2 {
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
		if a.AnalyzeNextDiscoveredTile() {
			return a.GeneratePath()
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
