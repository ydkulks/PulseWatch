package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string
}

func GetServerPort() *Config {
	// Load the environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Create a new Config struct
	return &Config{
		ServerPort: ":" + os.Getenv("SERVER_PORT"),
	}
}
