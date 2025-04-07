package main

import (
	"github.com/Erickype/DungeonPit/internal/logic"
	"github.com/Erickype/DungeonPit/internal/service"
	"github.com/Erickype/DungeonPit/internal/util"
	"github.com/joho/godotenv"
	"log"
	"net"
	"os"

	pb "github.com/Erickype/DungeonPit/proto"
	"google.golang.org/grpc"
)

func main() {
	loadEnv()

	util.LoadDatabase()

	serveApplication()
}

func loadEnv() {
	appEnv := os.Getenv("APP_ENV")

	// Skip loading .env file if running in production
	if appEnv != "production" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		log.Println(".env file loaded successfully")
	} else {
		log.Println("Running in production mode, .env file not loaded")
	}
}

func serveApplication() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	gameWorld := logic.NewGameWorld()
	pb.RegisterDungeonServiceServer(s, &service.Dungeon{World: gameWorld})

	log.Println("gRPC service listening on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
