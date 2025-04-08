package service

import (
	"context"
	"fmt"
	"github.com/Erickype/DungeonPit/internal/log"
	"github.com/Erickype/DungeonPit/internal/logic"
	"github.com/Erickype/DungeonPit/internal/model"
	"github.com/Erickype/DungeonPit/internal/util"
	pb "github.com/Erickype/DungeonPit/proto"
	"google.golang.org/grpc/status"
)

type Dungeon struct {
	pb.UnimplementedDungeonServiceServer
	World *logic.GameWorld
}

// Login Handle login
func (s *Dungeon) Login(_ context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	username := req.GetUsername()
	player, err := logic.LoginPlayer(username)
	if err != nil {
		return nil, err
	}

	initialRoom, err := logic.SetInitialPlayerRoom(player)
	if err != nil {
		return nil, err
	}

	s.World.AddPlayer(player, initialRoom)

	return &pb.LoginResponse{
		PlayerId: player.ID.String(),
	}, nil
}

func (s *Dungeon) Logout(_ context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	loggedPlayer, err := s.World.GetPlayer(req.GetPlayerId())
	if err != nil {
		return nil, err
	}

	err = loggedPlayer.SetCurrentRoomID(loggedPlayer.CurrentRoomID)
	if err != nil {
		return nil, err
	}

	log.GetCoreInstance().Info("Player", loggedPlayer.Username, "logout")
	s.World.RemovePlayer(loggedPlayer)

	return &pb.LogoutResponse{
		PlayerId: loggedPlayer.ID.String(),
	}, nil
}

// GetCurrentRoom Return a static room
func (s *Dungeon) GetCurrentRoom(_ context.Context, req *pb.PlayerRequest) (*pb.RoomResponse, error) {
	player, ok := s.World.Players[req.GetPlayerId()]
	if !ok {
		return nil, status.Errorf(404, "Player not found")
	}

	room := s.World.Rooms[player.CurrentRoomID.String()]
	return &pb.RoomResponse{
		RoomId: room.RoomID.String(),
	}, nil
}

// Move Handle movement
func (s *Dungeon) Move(_ context.Context, req *pb.MoveRequest) (*pb.RoomResponse, error) {
	player, err := s.World.GetPlayer(req.PlayerId)
	if err != nil {
		return nil, err
	}

	room, err := s.World.GetRoom(player.CurrentRoomID.String())
	if err != nil {
		return nil, err
	}

	direction := req.GetDirection()
	x, y, z := util.MoveCoordinates(room.X, room.Y, room.Z, direction)

	var newRoom = &model.Room{
		X: x,
		Y: y,
		Z: z,
	}
	err = newRoom.GetRoom()

	if err != nil {
		err = newRoom.Create()
		if err != nil {
			return nil, err
		}
	}

	message := fmt.Sprintf("Player %s moving to Room (%d, %d, %d)",
		player.Username, newRoom.X, newRoom.Y, newRoom.Z)
	log.GetCoreInstance().Info(message)

	s.World.AddRoom(newRoom, player)

	return &pb.RoomResponse{
		RoomId: newRoom.RoomID.String(),
	}, nil
}
