package tools

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"strings"
	"time"
)

func LogRaw(keys []string, res map[string]interface{}, indent bool) error {
	outputDir := os.Getenv("OUTPUT_DIR")
	if outputDir == "" {
		return errors.New("output directory is not set")
	}
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, 0755)
	}
	outputDir = path.Join(outputDir, "raw")
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, 0755)
	}
	timestamp := time.Now().Format("20060102150405")
	keys = append(keys, timestamp)
	encodedKeys := strings.Join(keys, "_")
	fileName := path.Join(outputDir, encodedKeys+".json")
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	var encodedResult []byte
	if indent {
		encodedResult, err = json.MarshalIndent(res, "", "  ")
	} else {
		encodedResult, err = json.Marshal(res)
	}
	if err != nil {
		return err
	}
	_, err = file.Write(encodedResult)
	if err != nil {
		return err
	}
	return nil
}

func Log(dir string, keys []string, res map[string]interface{}, indent bool) error {
	outputDir := os.Getenv("OUTPUT_DIR")
	if outputDir == "" {
		return errors.New("output directory is not set")
	}
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, 0755)
	}
	outputDir = path.Join(outputDir, dir)
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, 0755)
	}
	timestamp := time.Now().Format("20060102150405")
	keys = append(keys, timestamp)
	encodedKeys := strings.Join(keys, "_")
	fileName := path.Join(outputDir, encodedKeys+".json")
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	var encodedResult []byte
	if indent {
		encodedResult, err = json.MarshalIndent(res, "", "  ")
	} else {
		encodedResult, err = json.Marshal(res)
	}
	if err != nil {
		return err
	}
	_, err = file.Write(encodedResult)
	if err != nil {
		return err
	}
	return nil
}

func LogOverwrite(dir string, keys []string, res map[string]interface{}, indent bool) error {
	outputDir := os.Getenv("OUTPUT_DIR")
	if outputDir == "" {
		return errors.New("output directory is not set")
	}
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, 0755)
	}
	outputDir = path.Join(outputDir, dir)
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, 0755)
	}
	encodedKeys := strings.Join(keys, "_")
	fileName := path.Join(outputDir, encodedKeys+".json")
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	var encodedResult []byte
	if indent {
		encodedResult, err = json.MarshalIndent(res, "", "  ")
	} else {
		encodedResult, err = json.Marshal(res)
	}
	if err != nil {
		return err
	}
	_, err = file.Write(encodedResult)
	if err != nil {
		return err
	}
	return nil
}
