package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("Environment variables loaded successfully")
}
