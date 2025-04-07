package logic

import (
	"github.com/Erickype/DungeonPit/internal/model"
	"google.golang.org/grpc/status"
)

type IGameWorld interface {
	GetPlayer(playerID string) (*model.Player, error)
	AddPlayer(player *model.Player, currentRoom *model.Room)
	RemovePlayer(player *model.Player)
	GetRoom(roomID string) (*model.Room, error)
	AddRoom(room *model.Room, player *model.Player)
}

type GameWorld struct {
	Rooms   map[string]*model.Room
	Players map[string]*model.Player
}

func (gameWorld *GameWorld) GetPlayer(playerID string) (*model.Player, error) {
	loggedPlayer, ok := gameWorld.Players[playerID]
	if !ok {
		return nil, status.Errorf(404, "player not found")
	}
	return loggedPlayer, nil
}

func (gameWorld *GameWorld) AddPlayer(player *model.Player, currentRoom *model.Room) {
	playerID := player.ID.String()
	startRoomID := currentRoom.RoomID.String()
	gameWorld.Rooms[startRoomID] = currentRoom
	gameWorld.Players[playerID] = player
}

func (gameWorld *GameWorld) RemovePlayer(player *model.Player) {
	delete(gameWorld.Players, player.ID.String())
}

func (gameWorld *GameWorld) GetRoom(roomID string) (*model.Room, error) {
	room, ok := gameWorld.Rooms[roomID]
	if !ok {
		return nil, status.Errorf(404, "room not found")
	}
	return room, nil
}

func (gameWorld *GameWorld) AddRoom(room *model.Room, player *model.Player) {
	player.CurrentRoomID = room.RoomID
	gameWorld.Rooms[room.RoomID.String()] = room
	gameWorld.Players[player.ID.String()] = player
}

// NewGameWorld Creates a new GameWorld
func NewGameWorld() *GameWorld {
	return &GameWorld{
		Rooms:   make(map[string]*model.Room),
		Players: make(map[string]*model.Player),
	}
}
