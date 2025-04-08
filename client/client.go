package main

import (
	"context"
	"fmt"
	"github.com/Erickype/DungeonPit/internal/log"
	pb "github.com/Erickype/DungeonPit/proto"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.GetClientInstance().Fatal("Could not connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		if err != nil {
			log.GetClientInstance().Fatal("Could not close connection: %v", err)
		}
	}(conn)

	client := pb.NewDungeonServiceClient(conn)

	// 1. Login
	var username string
	fmt.Print("Enter username: ")
	_, err = fmt.Scan(&username)
	if err != nil {
		log.GetClientInstance().Fatal("Could not read username:", err)
	}
	loginResp, err := client.Login(context.Background(), &pb.LoginRequest{Username: username})
	if err != nil {
		log.GetClientInstance().Fatal("Login failed:", err)
	}
	log.GetClientInstance().Info("Logged in as:", loginResp.GetPlayerId())
	playerID := loginResp.GetPlayerId()

	// 2. GetCurrentRoom
	roomResp, err := client.GetCurrentRoom(context.Background(), &pb.PlayerRequest{PlayerId: playerID})
	if err != nil {
		log.GetClientInstance().Fatal("GetCurrentRoom failed:", err)
	}
	logRoom(roomResp)

	// 3. Move
	for {
		var moveDirection pb.Direction
		fmt.Print("Enter move direction: ")
		_, err = fmt.Scan(&moveDirection)
		if err != nil {
			log.GetClientInstance().Fatal("Could not read move direction:", err)
		}

		if moveDirection == pb.Direction_UNKNOWN {
			log.GetClientInstance().Warn("Invalid move direction, logging out...")
			break
		}

		moveResp, err := client.Move(context.Background(), &pb.MoveRequest{
			PlayerId:  playerID,
			Direction: moveDirection,
		})
		if err != nil {
			log.GetClientInstance().Fatal("Move failed:", err)
		}
		logRoom(moveResp)
	}

	// 4. Logout
	_, err = client.Logout(context.Background(), &pb.LogoutRequest{PlayerId: playerID})
	if err != nil {
		log.GetClientInstance().Trace("Logout failed:", err)
	} else {
		log.GetClientInstance().Info("Logged out successfully.")
	}
}

func logRoom(room *pb.RoomResponse) {
	log.GetClientInstance().Info("Room ID:", room.GetRoomId())
}
