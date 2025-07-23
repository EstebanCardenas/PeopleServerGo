package utils

import (
	"encoding/json"
	"io"
)

func GetRequestBodyAsMap(body io.ReadCloser) (map[string]any, error) {
	bodyBytes, bodyErr := io.ReadAll(body)
	if bodyErr != nil {
		return nil, bodyErr
	}
	bodyStr := string(bodyBytes)
	data, jsonErr := ParseJsonStrToMap(bodyStr)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return data, nil
}

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
