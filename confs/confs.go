package confs

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadConfig() error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return nil
}
