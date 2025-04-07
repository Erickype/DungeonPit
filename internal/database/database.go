package database

import (
	"fmt"
	"github.com/Erickype/DungeonPit/internal/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func InitDb() *gorm.DB {
	DB = connectDB()
	return DB
}

func connectDB() *gorm.DB {
	var err error
	host := os.Getenv("POSTGRES_HOST")
	username := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")
	sslMode := os.Getenv("SSL_MODE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=America/Guayaquil",
		host, username, password, dbname, port, sslMode)
	log.GetCoreInstance().Info("dsn:", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})

	if err != nil {
		log.GetCoreInstance().Fatal("Error connecting to database :", err)
		return nil
	}
	log.GetCoreInstance().Info("Successfully connected to the database")

	return db
}
