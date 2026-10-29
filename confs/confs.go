package confs

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() error {
	err := godotenv.Load(".env")
	if errors.Is(err, os.ErrNotExist) {
		// Container deployments provide configuration through environment variables.
		return nil
	}
	return err
}
