package core

import (
	"testing"

	"github.com/goki/mat32"
)

func TestDataRenderer2D_CalculateHallways(t *testing.T) {
	rooms := []Room2D{
		{Position: mat32.Vec2{X: 0, Y: 0}, Size: mat32.Vec2{X: 1, Y: 1}},
		{Position: mat32.Vec2{X: 2, Y: 0}, Size: mat32.Vec2{X: 1, Y: 1}},
		{Position: mat32.Vec2{X: 1, Y: 2}, Size: mat32.Vec2{X: 2, Y: 2}},
	}
	dungeon := DungeonGenerator2DSectionData{
		Size:  mat32.Vec2{X: 5, Y: 5},
		Rooms: rooms,
	}
	dungeon.InitializeGrid()
	for _, room := range dungeon.Rooms {
		dungeon.PlaceRoomPositionsInGrid(room)
	}
	dungeon.CalculateDelaunayTriangulation()
	dungeon.CalculateMST2D()
	dungeon.CalculateHallways()
	dungeon.GenerateGridRenderData()
	if len(dungeon.DataRenderer.GridLines) != 8 {
		t.Errorf("dungeon.DataRenderer.GridLines should be 8, calculated: %d", len(dungeon.DataRenderer.GridLines))
	}
}
