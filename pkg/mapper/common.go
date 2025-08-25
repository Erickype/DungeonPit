package mapper

import (
	core "github.com/Erickype/DungeonPit/internal/core/DungeonGenerator2D"
	"github.com/Erickype/DungeonPit/pkg/gen/common"
	"github.com/goki/mat32"
)

func ToProtoVector2D(v mat32.Vec2) *common.Vector2D {
	return &common.Vector2D{
		X: int64(v.X),
		Y: int64(v.Y),
	}
}

func FromProtoVector2D(v *common.Vector2D) mat32.Vec2 {
	if v == nil {
		return mat32.NewVec2(0, 0) // safe default
	}
	return mat32.NewVec2(float32(v.X), float32(v.Y))
}

func ToProtoLine2D(l core.Line2D) *common.Line2D {
	return &common.Line2D{
		A: ToProtoVector2D(mat32.NewVec2(l.A.X, l.A.Y)),
		B: ToProtoVector2D(mat32.NewVec2(l.B.X, l.B.Y)),
	}
}

func FromProtoLine2D(l *common.Line2D) core.Line2D {
	if l == nil {
		return core.Line2D{
			A: mat32.NewVec3(0, 0, 0),
			B: mat32.NewVec3(0, 0, 0),
		}
	}
	A := FromProtoVector2D(l.A)
	B := FromProtoVector2D(l.B)
	return core.Line2D{
		A: mat32.NewVec3(A.X, A.Y, 0),
		B: mat32.NewVec3(B.X, B.Y, 0),
	}
}
