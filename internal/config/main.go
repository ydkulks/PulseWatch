package config

import (
	"log"
	"os"
	"strconv"

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

type ResponseInterval struct {
	Interval int
}

func GetResponseInterval() *ResponseInterval {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	val := os.Getenv("RESPONSE_INTERVAL")
	if val == "" {
		log.Fatal("RESPONSE_INTERVAL is not set in .env file")
	}

	interval, err := strconv.Atoi(val)
	if err != nil {
		log.Fatal("Error parsing RESPONSE_INTERVAL from .env file")
	}

	return &ResponseInterval {
		Interval: interval,
	}
}
