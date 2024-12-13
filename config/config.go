package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var (
	PostgresDb       string
	PostgresUser     string
	PostgresPassword string
	DbPort           string
	DbHost           string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	DbHost = os.Getenv("DB_HOST")
	DbPort = os.Getenv("DB_PORT")
	PostgresDb = os.Getenv("POSTGRES_DB")
	PostgresUser = os.Getenv("POSTGRES_USER")
	PostgresPassword = os.Getenv("POSTGRES_PASSWORD")
}
