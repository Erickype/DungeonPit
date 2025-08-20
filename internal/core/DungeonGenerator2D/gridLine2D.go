package core

import (
	"github.com/goki/mat32"
)

type LineType int

const (
	LineTypeRoom LineType = iota
	LineTypeDoor
	LineTypeHallway
	LineTypeHallwayPath
)

type Line struct {
	A mat32.Vec2
	B mat32.Vec2
}

type GridLine struct {
	Line     Line
	LineType LineType
}
