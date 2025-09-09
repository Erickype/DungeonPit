package core

type GridLinePosition int

const (
	GridLinePositionTop GridLinePosition = iota
	GridLinePositionBottom
	GridLinePositionLeft
	GridLinePositionRight
	GridLinePositionDoor
	GridLinePositionHallwayPath
)

type GridLineType int

const (
	GridLineTypeRoom GridLineType = iota
	GridLineTypeDoor
	GridLineTypeHallway
	GridLineTypeHallwayPath
)

type GridLine struct {
	Line         Line2D
	LineType     GridLineType
	LinePosition GridLinePosition
}

func NewGridLine(line Line2D, lineType GridLineType, linePosition *GridLinePosition) *GridLine {
	if linePosition != nil {
		return &GridLine{
			Line:         line,
			LineType:     lineType,
			LinePosition: *linePosition,
		}
	}
	return &GridLine{
		Line:     line,
		LineType: lineType,
	}
}
