package core

import "github.com/goki/mat32"

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
	TwoVertexDirection() MoveDirection
	GetDirectionOrientation() MoveDirection
	PlaceDoor()
	PlaceHollowHallway()
}

type DataRenderer2D struct {
	Rooms     []Room2D
	Hallways  [][]mat32.Vec2
	GridLines []GridLine
	MSTEdges  []Line2D
}

func (d *DataRenderer2D) Calculate() {
	d.CalculateHallways()
	d.CalculateRooms()
}

func (d *DataRenderer2D) CalculateHallways() {
	//TODO implement me
	panic("implement me")
}

func (d *DataRenderer2D) CalculateRooms() {
	//TODO implement me
	panic("implement me")
}

func (d *DataRenderer2D) TwoVertexDirection() MoveDirection {
	//TODO implement me
	panic("implement me")
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

func NewDataRenderer2D(rooms []Room2D, hallways [][]mat32.Vec2, MSTEdges []Line2D) *DataRenderer2D {
	return &DataRenderer2D{
		Rooms:    rooms,
		Hallways: hallways,
		MSTEdges: MSTEdges,
	}
}
