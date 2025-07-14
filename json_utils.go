package main

import "encoding/json"

func ParseMapSliceToJsonStr(data []map[string]any) (string, error) {
	contentBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	content := string(contentBytes)
	return content, nil
}
