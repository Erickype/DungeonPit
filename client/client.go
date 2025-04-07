package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/Erickype/DungeonPit/proto"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewDungeonServiceClient(conn)

	// 1. Login
	var username string
	fmt.Print("Enter username: ")
	fmt.Scan(&username)
	loginResp, err := client.Login(context.Background(), &pb.LoginRequest{Username: username})
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}
	log.Printf("Logged in as %s", loginResp.GetPlayerId())
	playerID := loginResp.GetPlayerId()

	// Logout
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Caught interrupt signal, logging out...")
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		_, err := client.Logout(ctx, &pb.LogoutRequest{PlayerId: playerID})
		if err != nil {
			log.Printf("Logout failed: %v", err)
		} else {
			log.Println("Logged out successfully.")
		}
		os.Exit(0)
	}()

	// 2. GetCurrentRoom
	roomResp, err := client.GetCurrentRoom(context.Background(), &pb.PlayerRequest{PlayerId: playerID})
	if err != nil {
		log.Fatalf("GetCurrentRoom failed: %v", err)
	}
	logRoom(roomResp)

	// 3. Move
	for {
		var moveDirection pb.Direction
		fmt.Print("Enter move direction: ")
		fmt.Scan(&moveDirection)

		moveResp, err := client.Move(context.Background(), &pb.MoveRequest{
			PlayerId:  playerID,
			Direction: moveDirection,
		})
		if err != nil {
			log.Fatalf("Move failed: %v", err)
		}
		logRoom(moveResp)
	}

}

func logRoom(room *pb.RoomResponse) {
	log.Printf("Room ID: %s", room.GetRoomId())
}
