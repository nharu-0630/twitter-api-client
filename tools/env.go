package tools

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	for {
		if _, err := os.Stat(currentDir + "/.env"); err == nil {
			break
		}
		if currentDir == "/" {
			break
		}
		currentDir = currentDir[:len(currentDir)-1]
	}
	if err := godotenv.Load(currentDir + "/.env"); err != nil {
		return err
	}
	return nil
}
