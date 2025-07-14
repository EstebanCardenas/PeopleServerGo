package main

import (
	"encoding/csv"
	"fmt"
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
		fmt.Fprint(w, "Error parsing csv to JSON response")
		return
	}

	respContent, parseErr := ParseMapSliceToJsonStr(people)
	if parseErr != nil {
		setContentType(w, "text/plain")
		w.WriteHeader(500)
		fmt.Fprint(w, "Error parsing csv to JSON response")
		return
	}

	setContentType(w, "application/json")
	fmt.Fprint(w, respContent)
}

func PeopleHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		GetPeopleHandler(w, req)
	default:
		setContentType(w, "text/plain")
		w.WriteHeader(405)
		fmt.Fprint(w, "Unsupported HTTP method")
	}
}
