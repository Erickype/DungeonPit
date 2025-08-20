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
	CalculateDoors()
	TwoVertexDirection() MoveDirection
	GetDirectionOrientation() MoveDirection
	PlaceDoor()
	PlaceHollowHallway()
}

type DataRenderer2D struct {
	Rooms     []Room2D
	Hallways  [][]mat32.Vec2
	GridLines []GridLine
}

func (d *DataRenderer2D) Calculate() {
	//TODO implement me
	panic("implement me")
}

func (d *DataRenderer2D) CalculateHallways() {
	//TODO implement me
	panic("implement me")
}

func (d *DataRenderer2D) CalculateDoors() {
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

func NewDataRenderer2D(rooms []Room2D, hallways [][]mat32.Vec2) *DataRenderer2D {
	return &DataRenderer2D{
		Rooms:    rooms,
		Hallways: hallways,
	}
}
