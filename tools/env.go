package tools

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	for {
		envPath := filepath.Join(currentDir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			if err := godotenv.Load(envPath); err != nil {
				return err
			}
			return nil
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			break
		}
		currentDir = parentDir
	}

	return nil
}
