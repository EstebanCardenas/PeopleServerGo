package datasource

import (
	"encoding/csv"
	"errors"
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

func (ds *CsvDataSource) InitFile() error {
	_, err := os.Stat(fileName)
	if errors.Is(err, os.ErrNotExist) {
		// Create file
		file, err := os.Create(fileName)
		if err != nil {
			return err
		}
		defer file.Close()
		csvWriter := csv.NewWriter(file)
		err = csvWriter.Write(personFields[:])
		if err != nil {
			return err
		}
		csvWriter.Flush()
		if err = csvWriter.Error(); err != nil {
			return err
		}
	}

	return nil
}

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

func (ds *CsvDataSource) DeletePerson(id int) error {
	tempFileName := "temp_data.csv"
	inFile, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.Create(tempFileName)
	if err != nil {
		return err
	}
	defer outFile.Close()

	csvReader := csv.NewReader(inFile)
	csvWriter := csv.NewWriter(outFile)
	foundId := false
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if record[0] == fmt.Sprint(id) {
			foundId = true
			continue
		}

		err = csvWriter.Write(record)
		if err != nil {
			return err
		}
	}
	csvWriter.Flush()
	if err = csvWriter.Error(); err != nil {
		return err
	}

	err = os.Rename(tempFileName, fileName)
	if err != nil {
		return err
	}

	if !foundId {
		return PersonNotFoundError{id}
	}

	return nil
}

func (ds *CsvDataSource) UpdatePerson(id int, data map[string]any) error {
	if len(data) == 0 {
		return errors.New("payload contains no data")
	}

	tempFileName := "temp_data.csv"
	inFile, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.Create(tempFileName)
	if err != nil {
		return err
	}
	defer outFile.Close()

	csvReader := csv.NewReader(inFile)
	csvWriter := csv.NewWriter(outFile)
	foundId := false
	updatedFields := 0
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if record[0] == fmt.Sprint(id) {
			foundId = true
			newRecord := make([]string, len(personFields))
			newRecord[0] = fmt.Sprint(id)
			for i := 1; i < len(personFields); i++ {
				field := personFields[i]
				if value, ok := data[field]; ok {
					newRecord[i] = fmt.Sprint(value)
					updatedFields++
				} else {
					newRecord[i] = record[i]
				}
			}
			record = newRecord
		}

		err = csvWriter.Write(record)
		if err != nil {
			return err
		}
	}
	csvWriter.Flush()
	if err = csvWriter.Error(); err != nil {
		return err
	}

	err = os.Rename(tempFileName, fileName)
	if err != nil {
		return err
	}

	if !foundId {
		return PersonNotFoundError{id}
	}
	if updatedFields == 0 {
		return errors.New("no fields to update")
	}

	return nil
}
