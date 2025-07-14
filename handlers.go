package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func setContentType(w http.ResponseWriter, contentType string) {
	w.Header().Set("Content-Type", contentType)
}

func GetPeopleHandler(w http.ResponseWriter, req *http.Request) {
	csvContent, readErr := ReadCsvDataAsString()
	if readErr != nil {
		setContentType(w, "text/plain")
		w.WriteHeader(500)
		fmt.Fprint(w, "Error while reading CSV file")
		return
	}
	csvReader := csv.NewReader(strings.NewReader(csvContent))
	records, csvErr := csvReader.ReadAll()
	if csvErr != nil {
		setContentType(w, "text/plain")
		w.WriteHeader(500)
		fmt.Fprint(w, "Error while reading CSV file")
		return
	}
	people, rtmErr := RecordsToMap(records)
	if rtmErr != nil {
		setContentType(w, "text/plain")
		w.WriteHeader(500)
		fmt.Println(rtmErr)
		fmt.Fprint(w, "Error parsing csv to JSON response")
		return
	}

	respContent, parseErr := ParseMapSliceToJsonStr(people)
	if parseErr != nil {
		setContentType(w, "text/plain")
		w.WriteHeader(500)
		fmt.Println(parseErr)
		fmt.Fprint(w, "Error parsing csv to JSON response")
		return
	}

	setContentType(w, "application/json")
	fmt.Fprint(w, respContent)
}

func CreatePersonHandler(w http.ResponseWriter, req *http.Request) {
	bodyBytes, bodyErr := io.ReadAll(req.Body)
	if bodyErr != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	body := string(bodyBytes)
	data, jsonErr := ParseJsonStrToMap(body)
	if jsonErr != nil {
		http.Error(w, "Error parsing body to map", 500)
		return
	}
	writeErr := WritePersonMapToCsv(data)
	if writeErr != nil {
		http.Error(w, "Failed to write to csv file", 500)
		return
	}

	fmt.Fprint(w, "Created person successfully")
}

func PeopleHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		GetPeopleHandler(w, req)
	case "POST":
		CreatePersonHandler(w, req)
	default:
		setContentType(w, "text/plain")
		w.WriteHeader(405)
		fmt.Fprint(w, "Unsupported HTTP method")
	}
}
