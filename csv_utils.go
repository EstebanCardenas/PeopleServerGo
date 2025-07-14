package main

import "errors"

func RecordsToMap(records [][]string) ([]map[string]any, error) {
	if len(records) <= 1 {
		return nil, errors.New("there are no available records")
	}

	headerRow := records[0]
	result := make([]map[string]any, len(records)-1)
	for i := 1; i < len(records); i++ {
		row := records[i]
		data := make(map[string]any)
		for j, header := range headerRow {
			data[header] = row[j]
		}
		result[i-1] = data
	}
	return result, nil
}
