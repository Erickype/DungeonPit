package core

import (
	"github.com/Erickype/DungeonPit/internal/core"
	"github.com/goki/mat32"
	"math"
	"math/rand"
)

type IDungeonGenerator2DSectionData interface {
	InitializeGrid()
	PlaceRooms()
	TruncSize() (int, int)
	TruncRoomMaxSize() (int, int)
	CheckRoomOutBounds(newRoom Room2D) bool
	PlaceRoomPositionsInGrid(newRoom Room2D)
}

type DungeonGenerator2DSectionData struct {
	Size           mat32.Vec2
	RoomIterations int
	RoomMaxSize    mat32.Vec2
	Grid           map[mat32.Vec2]core.CellType
	Rooms          []Room2D
}

func (d *DungeonGenerator2DSectionData) TruncSize() (int, int) {
	sizeX := int(math.Trunc(float64(d.Size.X)))
	sizeY := int(math.Trunc(float64(d.Size.Y)))
	return sizeX, sizeY
}

func (d *DungeonGenerator2DSectionData) TruncRoomMaxSize() (int, int) {
	roomMaxsizeX := int(math.Trunc(float64(d.RoomMaxSize.X)))
	roomMaxsizeY := int(math.Trunc(float64(d.RoomMaxSize.Y)))
	return roomMaxsizeX, roomMaxsizeY
}

func (d *DungeonGenerator2DSectionData) CheckRoomOutBounds(newRoom Room2D) bool {
	c1 := newRoom.Position.X < 0
	c2 := newRoom.Position.X+newRoom.Size.X >= d.Size.X
	c3 := newRoom.Position.Y < 0
	c4 := newRoom.Position.Y+newRoom.Size.Y >= d.Size.Y

	return c1 || c2 || c3 || c4
}

func (d *DungeonGenerator2DSectionData) InitializeGrid() {
	d.Grid = make(map[mat32.Vec2]core.CellType)
	x, y := d.TruncSize()
	for i := 0; i < x; i++ {
		for j := 0; j < y; j++ {
			v := mat32.Vec2{
				X: float32(i),
				Y: float32(j),
			}
			d.Grid[v] = core.CellTypeNone
		}
	}
}

func (d *DungeonGenerator2DSectionData) PlaceRoomPositionsInGrid(newRoom Room2D) {
	positions := newRoom.AllPositionsWithin()
	for _, position := range positions {
		d.Grid[position] = core.CellTypeRoom
	}
}

func (d *DungeonGenerator2DSectionData) PlaceRooms() {
	d.Rooms = []Room2D{}
	sX, sY := d.TruncSize()
	rmsX, rmsY := d.TruncRoomMaxSize()
	for i := 0; i < d.RoomIterations; i++ {
		location := mat32.Vec2{
			X: float32(rand.Intn(sX)),
			Y: float32(rand.Intn(sY)),
		}
		roomSize := mat32.Vec2{
			X: float32(rand.Intn(rmsX) + 1),
			Y: float32(rand.Intn(rmsY) + 1),
		}
		bufferRoom := Room2D{
			Position: location.Sub(mat32.NewVec2Scalar(1)),
			Size:     roomSize.Add(mat32.NewVec2Scalar(2)),
		}
		newRoom := Room2D{
			Position: location,
			Size:     roomSize,
		}
		add := true
		for j := 0; j < len(d.Rooms); j++ {
			room := d.Rooms[j]
			if room.Intersect(bufferRoom) {
				add = false
				break
			}
		}
		if d.CheckRoomOutBounds(newRoom) {
			add = false
		}
		if add {
			d.Rooms = append(d.Rooms, newRoom)
			d.PlaceRoomPositionsInGrid(newRoom)
		}
	}
}

func GenerateDungeon2DSection(size mat32.Vec2, roomIterations int, roomMaxSize mat32.Vec2) *DungeonGenerator2DSectionData {
	dungeonGenerator2DSectionData := &DungeonGenerator2DSectionData{
		Size:           size,
		RoomIterations: roomIterations,
		RoomMaxSize:    roomMaxSize,
	}
	dungeonGenerator2DSectionData.InitializeGrid()
	dungeonGenerator2DSectionData.PlaceRooms()

	return dungeonGenerator2DSectionData
}
