package handlers

import (
	"fmt"
	"io"
	"net/http"
	"simple_server/application"
	"simple_server/utils"
)

func GetPeopleHandler(w http.ResponseWriter, req *http.Request) {
	people, err := application.DataSource.GetPeople()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error parsing csv to JSON response", 500)
		return
	}
	respContent, parseErr := utils.ParseMapSliceToJsonStr(people)
	if parseErr != nil {
		fmt.Println(parseErr)
		http.Error(w, "Error parsing csv to JSON response", 500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, respContent)
}

func CreatePersonHandler(w http.ResponseWriter, req *http.Request) {
	bodyBytes, bodyErr := io.ReadAll(req.Body)
	if bodyErr != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	body := string(bodyBytes)
	data, jsonErr := utils.ParseJsonStrToMap(body)
	if jsonErr != nil {
		http.Error(w, "Error parsing body to map", 500)
		return
	}
	writeErr := application.DataSource.SavePerson(data)
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
		http.Error(w, "Unsupported HTTP method", http.StatusMethodNotAllowed)
	}
}
