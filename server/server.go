package main

import (
	"net"
	"os"

	"github.com/Erickype/DungeonPit/internal/log"
	"github.com/Erickype/DungeonPit/internal/logic"
	"github.com/Erickype/DungeonPit/internal/service"
	service2 "github.com/Erickype/DungeonPit/internal/service/dungeon"
	"github.com/Erickype/DungeonPit/internal/util"
	"github.com/Erickype/DungeonPit/pkg/gen/dungeon"
	"github.com/joho/godotenv"

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
			log.GetCoreInstance().Fatal("Error loading .env file")
		}
		log.GetCoreInstance().Info(".env file loaded successfully")
	} else {
		log.GetCoreInstance().Info("Running in production mode, .env file not loaded")
	}
}

func serveApplication() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.GetCoreInstance().Fatal("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	gameWorld := logic.NewGameWorld()
	pb.RegisterDungeonServiceServer(s, &service.Dungeon{World: gameWorld})
	dungeon.RegisterDungeonServiceServer(s, &service2.Dungeon{})

	log.GetCoreInstance().Info("gRPC service listening on port 50051...")
	if err = s.Serve(lis); err != nil {
		log.GetCoreInstance().Fatal("Failed to serve: %v", err)
	}
}
