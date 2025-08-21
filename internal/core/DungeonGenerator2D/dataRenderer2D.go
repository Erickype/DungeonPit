package core

import (
	"math"
	"slices"

	"github.com/Erickype/DungeonPit/internal/core"
	"github.com/goki/mat32"
)

type MoveDirection int

const (
	MoveDirectionUp MoveDirection = iota
	MoveDirectionDown
	MoveDirectionRight
	MoveDirectionLeft
	MoveDirectionNone
)

type IDataRenderer2D interface {
	Calculate()
	CalculateHallways()
	CalculateRooms()
	TwoVertexDirection(vi, vf mat32.Vec2) MoveDirection
	GetDirectionOrientation() MoveDirection
	PlaceDoor()
	PlaceHollowHallway()
}

type DataRenderer2D struct {
	Grid      map[mat32.Vec2]core.CellType
	Rooms     []Room2D
	Hallways  [][]mat32.Vec2
	GridLines []GridLine
	MSTEdges  []Line2D
}

func (d *DataRenderer2D) Calculate() {
	d.GridLines = make([]GridLine, 0)
	d.CalculateHallways()
	d.CalculateRooms()
}

func (d *DataRenderer2D) CalculateHallways() {
	//TODO implement me
	panic("implement me")
}

func (d *DataRenderer2D) CalculateRooms() {
	for i, hallway := range d.Hallways {
		mstItem := d.MSTEdges[i]
		vi := mat32.NewVec2(float32(math.Floor(float64(mstItem.A.X))), float32(math.Floor(float64(mstItem.A.Y))))
		vf := hallway[0]
		previousDirections := d.TwoVertexDirection(vi, vf)
		slices.Insert(hallway, 0, vi)
		for j := 0; j < len(hallway)-2; j++ {
			vi = hallway[j]
			vf = hallway[j+1]
			currentDirections := d.TwoVertexDirection(vi, vf)
			switch d.Grid[vi] {
			case core.CellTypeRoom:
				break
			case core.CellTypeHallway:
				break
			case core.CellTypeNone:
				break
			}
		}
	}
}

func (d *DataRenderer2D) TwoVertexDirection(vi, vf mat32.Vec2) MoveDirection {
	if vf.X == vi.X && vf.Y == vi.Y {
		return MoveDirectionDown
	}
	if vf.X < vi.X && vf.Y == vi.Y {
		return MoveDirectionLeft
	}
	if vf.X == vi.X && vf.Y < vi.Y {
		return MoveDirectionUp
	}
	if vf.X > vi.X && vf.Y == vi.Y {
		return MoveDirectionRight
	}
	return MoveDirectionNone
}

func (d *DataRenderer2D) GetDirectionOrientation() MoveDirection {
	//TODO implement me
	panic("implement me")
}

func (d *DataRenderer2D) PlaceDoor() {
	//TODO implement me
	panic("implement me")
}

func (d *DataRenderer2D) PlaceHollowHallway() {
	//TODO implement me
	panic("implement me")
}

func NewDataRenderer2D(grid map[mat32.Vec2]core.CellType, rooms []Room2D, hallways [][]mat32.Vec2, MSTEdges []Line2D) *DataRenderer2D {
	return &DataRenderer2D{
		Grid:     grid,
		Rooms:    rooms,
		Hallways: hallways,
		MSTEdges: MSTEdges,
	}
}
