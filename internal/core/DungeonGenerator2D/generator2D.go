package core

import (
	"math"
	"math/rand"

	"github.com/Erickype/DungeonPit/internal/core"
	"github.com/goki/mat32"
)

type IDungeonGenerator2DSectionData interface {
	InitializeGrid()
	PlaceRooms()
	TruncSize() (int, int)
	TruncRoomMaxSize() (int, int)
	CheckRoomOutBounds(newRoom Room2D) bool
	PlaceRoomPositionsInGrid(newRoom Room2D)
	CalculateRoomsCenters()
	CalculateDelaunayTriangulation()
	CalculateMST2D()
	CalculateHallways()
}

type DungeonGenerator2DSectionData struct {
	Size                  mat32.Vec2
	RoomIterations        int
	RoomMaxSize           mat32.Vec2
	Grid                  map[mat32.Vec2]core.CellType
	Rooms                 []Room2D
	RoomsCenters          []mat32.Vec2
	DelaunayTriangles     []mat32.Triangle
	DelaunayEdgesSet      []Line2D
	MinimumSpanningTree2D *MinimumSpanningTree2D
	Hallways              [][]mat32.Vec2
}

func (d *DungeonGenerator2DSectionData) CalculateHallways() {
	d.Hallways = make([][]mat32.Vec2, 0)
	for _, edge := range d.MinimumSpanningTree2D.MSTEdges {
		aX := float32(math.Floor(float64(edge.A.X)))
		aY := float32(math.Floor(float64(edge.A.Y)))
		bX := float32(math.Floor(float64(edge.B.X)))
		bY := float32(math.Floor(float64(edge.B.Y)))
		start := mat32.NewVec2(aX, aY)
		end := mat32.NewVec2(bX, bY)
		aStar2D := NewAStar2D(start, end, d.Grid)
		hallway := aStar2D.FindPath()
		if hallway != nil {
			d.Hallways = append(d.Hallways, hallway)
		}
	}
}

func (d *DungeonGenerator2DSectionData) CalculateMST2D() {
	d.MinimumSpanningTree2D = NewMinimumSpanningTree2D(d.DelaunayEdgesSet)
	d.MinimumSpanningTree2D.CalculatePrimDistances()
	d.MinimumSpanningTree2D.CalculateMST2D()
}

func (d *DungeonGenerator2DSectionData) CalculateRoomsCenters() {
	d.RoomsCenters = make([]mat32.Vec2, len(d.Rooms))
	for i, room := range d.Rooms {
		center := mat32.Vec2{
			X: (room.Position.X + room.Size.X) - (room.Size.X / 2),
			Y: (room.Position.Y + room.Size.Y) - (room.Size.Y / 2),
		}
		d.RoomsCenters[i] = center
	}
}

func (d *DungeonGenerator2DSectionData) CalculateDelaunayTriangulation() {
	d.CalculateRoomsCenters()
	delaunayTriangulation2D := NewDelaunayTriangulation2D(d.RoomsCenters)
	delaunayTriangulation2D.Calculate()
	d.DelaunayTriangles = delaunayTriangulation2D.Triangles
	delaunayTriangulation2D.GenerateEdgesSet()
	d.DelaunayEdgesSet = delaunayTriangulation2D.EdgesSet
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
	dungeonGenerator2DSectionData.CalculateDelaunayTriangulation()
	dungeonGenerator2DSectionData.CalculateMST2D()
	dungeonGenerator2DSectionData.CalculateHallways()

	return dungeonGenerator2DSectionData
}
