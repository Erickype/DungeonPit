package core

type GridLineType int

const (
	GridLineTypeRoom GridLineType = iota
	GridLineTypeDoor
	GridLineTypeHallway
	GridLineTypeHallwayPath
)

type GridLine struct {
	Line     Line2D
	LineType GridLineType
}

func NewGridLine(line Line2D, lineType GridLineType) *GridLine {
	return &GridLine{
		Line:     line,
		LineType: lineType,
	}
}
