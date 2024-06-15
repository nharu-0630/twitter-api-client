package tools

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"
)

func LogRaw(keys []string, res map[string]interface{}) {
	outputDir := os.Getenv("OUTPUT_DIR")
	if outputDir == "" {
		return
	}
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, 0755)
	}
	timestamp := time.Now().Format("20060102150405")
	keys = append(keys, timestamp)
	encodedKeys := strings.Join(keys, "_")
	fileName := outputDir + "/" + encodedKeys + ".json"
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Failed to create file: %s", err)
	}
	defer file.Close()
	encodedResult, err := json.Marshal(res)
	if err != nil {
		log.Fatalf("Failed to encode result: %s", err)
	}
	_, err = file.Write(encodedResult)
	if err != nil {
		log.Fatalf("Failed to write to file: %s", err)
	}
}
