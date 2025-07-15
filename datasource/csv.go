package datasource

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const fileName = "data.csv"

func readCsvDataAsString() (string, error) {
	contentBytes, readErr := os.ReadFile(fileName)
	if readErr != nil {
		return "", readErr
	}
	content := string(contentBytes)
	return content, nil
}

func recordsToMap(records [][]string) ([]map[string]any, error) {
	result := make([]map[string]any, len(records)-1)
	for i := 1; i < len(records); i++ {
		row := records[i]
		data := make(map[string]any)
		data["id"], _ = strconv.Atoi(row[0])
		data["name"] = row[1]
		data["lastName"] = row[2]
		data["profession"] = row[3]
		data["age"], _ = strconv.Atoi(row[4])
		result[i-1] = data
	}
	return result, nil
}

func getLastCsvRecord() ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var lastRecord []string

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break // End of file
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV record: %w", err)
		}
		lastRecord = record // Keep updating lastRecord with the current one
	}

	if lastRecord == nil {
		return nil, fmt.Errorf("file is empty or contains no records")
	}

	return lastRecord, nil
}

func writePersonMapToCsv(personMap map[string]any) error {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	lastRecord, lastErr := getLastCsvRecord()
	if lastErr != nil {
		return lastErr
	}
	var lastIdx int
	if lastRecord[0] == "id" {
		lastIdx = 0
	} else {
		lastIdx, _ = strconv.Atoi(lastRecord[0])
	}
	content := []string{
		fmt.Sprint(lastIdx + 1),
		fmt.Sprint(personMap["name"]),
		fmt.Sprint(personMap["lastName"]),
		fmt.Sprint(personMap["profession"]),
		fmt.Sprint(personMap["age"]),
	}
	csvWriter := csv.NewWriter(file)
	csvErr := csvWriter.Write(content)
	if csvErr != nil {
		return csvErr
	}
	csvWriter.Flush()

	return nil
}

type CsvDataSource struct{}

func (ds *CsvDataSource) GetPeople() ([]map[string]any, error) {
	csvContent, readErr := readCsvDataAsString()
	if readErr != nil {
		return nil, readErr
	}
	csvReader := csv.NewReader(strings.NewReader(csvContent))
	records, csvErr := csvReader.ReadAll()
	if csvErr != nil {
		return nil, csvErr
	}
	people, rtmErr := recordsToMap(records)
	if rtmErr != nil {
		return nil, rtmErr
	}

	return people, nil
}

func (ds *CsvDataSource) SavePerson(person map[string]any) error {
	err := writePersonMapToCsv(person)
	if err != nil {
		return err
	}

	return nil
}
