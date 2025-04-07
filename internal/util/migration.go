package util

import (
	"github.com/Erickype/DungeonPit/internal/database"
	"github.com/Erickype/DungeonPit/internal/model"
	"github.com/google/uuid"
	"log"
)

func LoadDatabase() {
	database.InitDb()

	err := database.DB.AutoMigrate(
		&model.Room{},
		&model.Player{},
	)

	if err != nil {
		log.Fatal("Error while migrating: ", err.Error())
	}

	//SeedData()
}

func SeedData() {
	var initialRoom = model.Room{
		X: 0,
		Y: 0,
		Z: 0,
	}

	var player = model.Player{
		ID:       uuid.New(),
		Email:    "erickype@hotmail.com",
		Username: "Erickype",
	}

	database.DB.Create(&initialRoom)
	database.DB.Create(&player)
}
