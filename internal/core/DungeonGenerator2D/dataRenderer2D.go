package core

import (
	"math"
	"slices"

	"github.com/Erickype/DungeonPit/internal/core"
	"github.com/goki/mat32"
)

type MoveDirectionOrientation int

const (
	MoveDirectionOrientationHorizontal MoveDirectionOrientation = iota
	MoveDirectionOrientationVertical
	MoveDirectionOrientationUnknown
)

type IMoveDirection interface {
	GetOrientation() MoveDirectionOrientation
}

type MoveDirection int

const (
	MoveDirectionUp MoveDirection = iota
	MoveDirectionDown
	MoveDirectionRight
	MoveDirectionLeft
	MoveDirectionNone
)

func (m *MoveDirection) GetOrientation() MoveDirectionOrientation {
	switch *m {
	case MoveDirectionUp:
		return MoveDirectionOrientationVertical
	case MoveDirectionDown:
		return MoveDirectionOrientationVertical
	case MoveDirectionRight:
		return MoveDirectionOrientationHorizontal
	case MoveDirectionLeft:
		return MoveDirectionOrientationHorizontal
	default:
		return MoveDirectionOrientationUnknown
	}
}

type IDataRenderer2D interface {
	Calculate()
	CalculateHallways()
	CalculateRooms()
	TwoVertexDirection(vi, vf mat32.Vec2) MoveDirection
	PlaceDoor(moveDirection MoveDirection, vi, vf mat32.Vec2)
	GridLinesAddUniqueDoor(line Line2D)
	PlaceHollowHallway(v mat32.Vec2, orientation MoveDirectionOrientation)
	CanCreateHallwayLine(line Line2D) (bool, int)
	GridLinesAddUniqueWall(line Line2D, lineType GridLineType)
	PlaceHollowHallwayLinesCheck(hallwayLine Line2D, hallwayPathLine Line2D)
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

func (d *DataRenderer2D) CalculateRooms() {
	//TODO implement me
	panic("implement me")
}

func (d *DataRenderer2D) CalculateHallways() {
	for i, hallway := range d.Hallways {
		mstItem := d.MSTEdges[i]
		vi := mat32.NewVec2(float32(math.Floor(float64(mstItem.A.X))), float32(math.Floor(float64(mstItem.A.Y))))
		vf := hallway[0]
		previousDirection := d.TwoVertexDirection(vi, vf)
		slices.Insert(hallway, 0, vi)
		for j := 0; j < len(hallway)-2; j++ {
			vi = hallway[j]
			vf = hallway[j+1]
			currentDirection := d.TwoVertexDirection(vi, vf)
			switch d.Grid[vi] {
			case core.CellTypeRoom:
				if d.Grid[vf] == core.CellTypeHallway {
					d.PlaceDoor(currentDirection, vi, vf)
				}
				previousDirection = currentDirection
				break
			case core.CellTypeHallway:
				if previousDirection == currentDirection {
					d.PlaceHollowHallway(vi, currentDirection.GetOrientation())
					previousDirection = currentDirection
					if d.Grid[vf] == core.CellTypeRoom {
						d.PlaceDoor(currentDirection, vi, vf)
					}
				}
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

func (d *DataRenderer2D) PlaceDoor(moveDirection MoveDirection, vi, vf mat32.Vec2) {
	viVec3 := mat32.NewVec3(vi.X, vi.Y, 0)
	vfVec3 := mat32.NewVec3(vf.X, vf.Y, 0)
	switch moveDirection {
	case MoveDirectionUp:
		d.GridLinesAddUniqueDoor(*NewLine2D(viVec3, mat32.NewVec3(vi.X+1, vi.Y, 0)))
		break
	case MoveDirectionDown:
		d.GridLinesAddUniqueDoor(*NewLine2D(vfVec3, mat32.NewVec3(vf.X+1, vf.Y, 0)))
		break
	case MoveDirectionRight:
		d.GridLinesAddUniqueDoor(*NewLine2D(vfVec3, mat32.NewVec3(vf.X, vf.Y+1, 0)))
		break
	case MoveDirectionLeft:
		d.GridLinesAddUniqueDoor(*NewLine2D(viVec3, mat32.NewVec3(vi.X, vi.Y+1, 0)))
		break
	default:
		break
	}
}

func (d *DataRenderer2D) GridLinesAddUniqueDoor(line Line2D) {
	found := false
	foundIndex := -1
	for i, gridLine := range d.GridLines {
		if gridLine.Line.IsSameLine2D(line) || gridLine.Line.IsSameLine2D(*NewLine2D(line.B, line.A)) {
			found = true
			foundIndex = i
			break
		}
	}
	if !found {
		d.GridLines = append(d.GridLines, *NewGridLine(line, GridLineTypeDoor))
	}
	foundLineType := d.GridLines[foundIndex].LineType
	if foundLineType == GridLineTypeHallway || foundLineType == GridLineTypeHallwayPath {
		d.GridLines[foundIndex].LineType = GridLineTypeDoor
	}
}

func (d *DataRenderer2D) PlaceHollowHallway(v mat32.Vec2, orientation MoveDirectionOrientation) {
	switch orientation {
	case MoveDirectionOrientationHorizontal:
		for i := 0; i < 2; i++ {
			hallwayLine := NewLine2D(mat32.NewVec3(v.X, v.Y+float32(i), 0), mat32.NewVec3(v.X+1, v.Y+float32(i), 0))
			hallwayPathLine := NewLine2D(mat32.NewVec3(v.X+float32(i), v.Y, 0), mat32.NewVec3(v.X+float32(i), v.Y+1, 0))
			d.PlaceHollowHallwayLinesCheck(*hallwayLine, *hallwayPathLine)
		}
	case MoveDirectionOrientationVertical:
		for i := 0; i < 2; i++ {
			hallwayLine := NewLine2D(mat32.NewVec3(v.X+float32(i), v.Y, 0), mat32.NewVec3(v.X+float32(i), v.Y+1, 0))
			hallwayPathLine := NewLine2D(mat32.NewVec3(v.X, v.Y+float32(i), 0), mat32.NewVec3(v.X+1, v.Y+float32(i), 0))
			d.PlaceHollowHallwayLinesCheck(*hallwayLine, *hallwayPathLine)
		}
	case MoveDirectionOrientationUnknown:
	}
}

func (d *DataRenderer2D) PlaceHollowHallwayLinesCheck(hallwayLine Line2D, hallwayPathLine Line2D) {
	canCreate, _ := d.CanCreateHallwayLine(hallwayLine)
	if canCreate {
		d.GridLinesAddUniqueWall(hallwayLine, GridLineTypeHallway)
	}
	canCreate, index := d.CanCreateHallwayLine(hallwayPathLine)
	if canCreate {
		d.GridLinesAddUniqueWall(hallwayLine, GridLineTypeHallwayPath)
	}
	if d.GridLines[index].LineType == GridLineTypeHallway {
		d.GridLines[index].LineType = GridLineTypeHallwayPath
	}
}

func (d *DataRenderer2D) CanCreateHallwayLine(line Line2D) (bool, int) {
	candidates := []int{
		slices.Index(d.GridLines, *NewGridLine(line, GridLineTypeHallway)),
		slices.Index(d.GridLines, *NewGridLine(*NewLine2D(line.B, line.A), GridLineTypeHallway)),
		slices.Index(d.GridLines, *NewGridLine(line, GridLineTypeHallwayPath)),
		slices.Index(d.GridLines, *NewGridLine(*NewLine2D(line.B, line.A), GridLineTypeHallwayPath)),
		slices.Index(d.GridLines, *NewGridLine(line, GridLineTypeDoor)),
		slices.Index(d.GridLines, *NewGridLine(*NewLine2D(line.B, line.A), GridLineTypeDoor)),
	}
	// Check if no candidate exists → line can be created
	for _, c := range candidates {
		if c != -1 {
			return true, -1
		}
	}
	// Otherwise return the first match (including reversed line)
	for _, c := range candidates {
		if c != -1 {
			return false, c
		}
	}
	// Fallback, shouldn't happen
	return true, -1

}

func (d *DataRenderer2D) GridLinesAddUniqueWall(line Line2D, lineType GridLineType) {
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
