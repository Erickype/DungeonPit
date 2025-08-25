package service

import (
	"context"

	core "github.com/Erickype/DungeonPit/internal/core/DungeonGenerator2D"
	"github.com/Erickype/DungeonPit/pkg/gen/common"
	"github.com/Erickype/DungeonPit/pkg/gen/dungeon"
	"github.com/Erickype/DungeonPit/pkg/mapper"
	"github.com/goki/mat32"
)

type Dungeon struct {
	dungeon.UnimplementedDungeonServiceServer
}

func (d *Dungeon) GenerateDungeon(_ context.Context, request *dungeon.GenerateDungeonRequest) (*dungeon.GenerateDungeonResponse, error) {
	size := mat32.NewVec2(float32(request.SizeX), float32(request.SizeY))
	roomSize := mat32.NewVec2(float32(request.RoomMaxSizeX), float32(request.RoomMaxSizeY))
	generatedDungeon := core.GenerateDungeon2DSection(size, int(request.RoomIterations), roomSize)
	gridLines := generatedDungeon.DataRenderer.GridLines
	gridLinesResponse := make([]*common.GridLine, len(gridLines))
	for i, gridLine := range gridLines {
		protoGridLine := &common.GridLine{
			Line:     mapper.ToProtoLine2D(gridLine.Line),
			LineType: common.GridLineType(gridLine.LineType),
		}
		gridLinesResponse[i] = protoGridLine
	}
	return &dungeon.GenerateDungeonResponse{GridLines: gridLinesResponse}, nil
}
