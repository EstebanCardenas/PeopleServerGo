package main

import (
	"encoding/json"
)

func ParseMapSliceToJsonStr(data []map[string]any) (string, error) {
	contentBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	content := string(contentBytes)
	return content, nil
}

func ParseJsonStrToMap(str string) (map[string]any, error) {
	result := make(map[string]any)
	err := json.Unmarshal([]byte(str), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
