package tools

import (
	"os"
	"path"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	currentDir, err := os.Getwd()
	envPath := path.Join(currentDir, ".env")
	if err != nil {
		return err
	}
	for {
		if _, err := os.Stat(envPath); err == nil {
			break
		}
		if currentDir == "/" {
			break
		}
		currentDir = path.Dir(currentDir)
		envPath = path.Join(currentDir, ".env")
	}
	if err := godotenv.Load(envPath); err != nil {
		return err
	}
	return nil
}
