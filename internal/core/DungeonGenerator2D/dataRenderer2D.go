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
	PlaceCornerHallway(v mat32.Vec2, previousDirection MoveDirection, currentDirection MoveDirection)
	PlaceCorner(corner Corner, v mat32.Vec2)
	GridLinesFindIndexByLine(line Line2D) (bool, int)
	GridLinesAddRoomLine(line Line2D)
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
	for _, room := range d.Rooms {
		for i := 0; i < int(room.Size.X); i++ {
			topLine := NewLine2D(
				mat32.NewVec3(float32(i)+room.Position.X, room.Position.Y, 0),
				mat32.NewVec3(float32(i+1)+room.Position.X, room.Position.Y, 0))
			bottomLine := NewLine2D(
				mat32.NewVec3(float32(i)+room.Position.X, room.Position.Y+room.Size.Y, 0),
				mat32.NewVec3(float32(i+1)+room.Position.X, room.Position.Y+room.Size.Y, 0))
			d.GridLinesAddRoomLine(*topLine)
			d.GridLinesAddRoomLine(*bottomLine)
		}
		for i := 0; i < int(room.Size.Y); i++ {
			rightLine := NewLine2D(
				mat32.NewVec3(room.Position.X, float32(i)+room.Position.Y, 0),
				mat32.NewVec3(room.Position.X, float32(i+1)+room.Position.Y, 0))
			leftLine := NewLine2D(
				mat32.NewVec3(room.Position.X+room.Size.X, float32(i)+room.Position.Y, 0),
				mat32.NewVec3(room.Position.X+room.Size.X, float32(i+1)+room.Position.Y, 0))
			d.GridLinesAddRoomLine(*rightLine)
			d.GridLinesAddRoomLine(*leftLine)
		}
	}
}

func (d *DataRenderer2D) GridLinesAddRoomLine(line Line2D) {
	found, index := d.GridLinesFindIndexByLine(line)
	gridLine := NewGridLine(line, GridLineTypeRoom)
	if !found {
		d.GridLines = append(d.GridLines, *gridLine)
	}
	if found && d.GridLines[index].LineType != GridLineTypeDoor {
		d.GridLines = append(d.GridLines, *gridLine)
	}
}

func (d *DataRenderer2D) CalculateHallways() {
	for i, hallway := range d.Hallways {
		mstItem := d.MSTEdges[i]
		vi := mat32.NewVec2(float32(math.Floor(float64(mstItem.A.X))), float32(math.Floor(float64(mstItem.A.Y))))
		vf := hallway[0]
		previousDirection := d.TwoVertexDirection(vi, vf)
		hallway = slices.Insert(hallway, 0, vi)
		for j := 0; j <= len(hallway)-2; j++ {
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
				} else {
					d.PlaceCornerHallway(vi, previousDirection, currentDirection)
					previousDirection = currentDirection
					if d.Grid[vf] == core.CellTypeRoom {
						d.PlaceDoor(currentDirection, vi, vf)
					}
				}
				break
			default:
				break
			}
		}
	}
}

func (d *DataRenderer2D) TwoVertexDirection(vi, vf mat32.Vec2) MoveDirection {
	if vf.X == vi.X && vf.Y > vi.Y {
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

func (d *DataRenderer2D) GridLinesFindIndexByLine(line Line2D) (bool, int) {
	found := false
	foundIndex := -1
	for i, gridLine := range d.GridLines {
		if gridLine.Line.IsSameLine2D(line) || gridLine.Line.IsSameLine2D(*NewLine2D(line.B, line.A)) {
			found = true
			foundIndex = i
			break
		}
	}
	return found, foundIndex
}

func (d *DataRenderer2D) GridLinesAddUniqueDoor(line Line2D) {
	found, foundIndex := d.GridLinesFindIndexByLine(line)
	if !found {
		d.GridLines = append(d.GridLines, *NewGridLine(line, GridLineTypeDoor))
		return
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
		d.GridLinesAddUniqueWall(hallwayPathLine, GridLineTypeHallwayPath)
		return
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
	found, foundIndex := d.GridLinesFindIndexByLine(line)
	if !found {
		d.GridLines = append(d.GridLines, *NewGridLine(line, lineType))
		return
	}
	if d.GridLines[foundIndex].LineType != GridLineTypeDoor {
		if d.GridLines[foundIndex].LineType != lineType {
			d.GridLines[foundIndex].LineType = lineType
		}
	}
}

func (d *DataRenderer2D) PlaceCornerHallway(v mat32.Vec2, previousDirection MoveDirection, currentDirection MoveDirection) {
	switch previousDirection {
	case MoveDirectionUp:
		switch currentDirection {
		case MoveDirectionRight:
			d.PlaceCorner(CornerUpRight, v)
			break
		case MoveDirectionLeft:
			d.PlaceCorner(CornerUpLeft, v)
			break
		default:
			break
		}
		break
	case MoveDirectionDown:
		switch currentDirection {
		case MoveDirectionRight:
			d.PlaceCorner(CornerDownRight, v)
			break
		case MoveDirectionLeft:
			d.PlaceCorner(CornerDownLeft, v)
			break
		default:
			break
		}
		break
	case MoveDirectionRight:
		switch currentDirection {
		case MoveDirectionUp:
			d.PlaceCorner(CornerDownLeft, v)
			break
		case MoveDirectionDown:
			d.PlaceCorner(CornerUpLeft, v)
			break
		default:
			break
		}
		break
	case MoveDirectionLeft:
		switch currentDirection {
		case MoveDirectionUp:
			d.PlaceCorner(CornerDownRight, v)
			break
		case MoveDirectionDown:
			d.PlaceCorner(CornerUpRight, v)
			break
		default:
			break
		}
		break
	case MoveDirectionNone:
		break
	}
}

type Corner int

const (
	CornerDownLeft Corner = iota
	CornerDownRight
	CornerUpLeft
	CornerUpRight
)

// PlaceCorner
/* Places de 4 different corners of type Corner:
╝ == 0
╚ == 1
╗ == 2
╔ == 3
*/
func (d *DataRenderer2D) PlaceCorner(corner Corner, v mat32.Vec2) {
	var lineA Line2D
	var lineB Line2D
	switch corner {
	case CornerDownLeft:
		lineA = *NewLine2D(mat32.NewVec3(v.X+1, v.Y, 0), mat32.NewVec3(v.X+1, v.Y+1, 0))
		lineB = *NewLine2D(mat32.NewVec3(v.X+1, v.Y+1, 0), mat32.NewVec3(v.X, v.Y+1, 0))
		break
	case CornerDownRight:
		lineA = *NewLine2D(mat32.NewVec3(v.X, v.Y, 0), mat32.NewVec3(v.X, v.Y+1, 0))
		lineB = *NewLine2D(mat32.NewVec3(v.X, v.Y+1, 0), mat32.NewVec3(v.X+1, v.Y+1, 0))
		break
	case CornerUpLeft:
		lineA = *NewLine2D(mat32.NewVec3(v.X, v.Y, 0), mat32.NewVec3(v.X+1, v.Y, 0))
		lineB = *NewLine2D(mat32.NewVec3(v.X+1, v.Y, 0), mat32.NewVec3(v.X+1, v.Y+1, 0))
		break
	case CornerUpRight:
		lineA = *NewLine2D(mat32.NewVec3(v.X, v.Y, 0), mat32.NewVec3(v.X+1, v.Y, 0))
		lineB = *NewLine2D(mat32.NewVec3(v.X, v.Y, 0), mat32.NewVec3(v.X, v.Y+1, 0))
		break
	}
	canCreate, _ := d.CanCreateHallwayLine(lineA)
	if canCreate {
		d.GridLinesAddUniqueWall(lineA, GridLineTypeHallway)
	}
	canCreate, _ = d.CanCreateHallwayLine(lineB)
	if canCreate {
		d.GridLinesAddUniqueWall(lineB, GridLineTypeHallway)
	}
}

func NewDataRenderer2D(grid map[mat32.Vec2]core.CellType, rooms []Room2D, hallways [][]mat32.Vec2, MSTEdges []Line2D) *DataRenderer2D {
	return &DataRenderer2D{
		Grid:     grid,
		Rooms:    rooms,
		Hallways: hallways,
		MSTEdges: MSTEdges,
	}
}
