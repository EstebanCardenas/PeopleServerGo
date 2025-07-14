package main

import (
	"os"
)

const fileName = "data.csv"

func ReadCsvDataAsString() (string, error) {
	contentBytes, readErr := os.ReadFile(fileName)
	if readErr != nil {
		return "", readErr
	}
	content := string(contentBytes)
	return content, nil
}
