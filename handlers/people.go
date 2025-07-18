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
		errorRes := fmt.Sprintf("Error while reading people from database: %v", err)
		http.Error(w, errorRes, 500)
		return
	}
	respContent, parseErr := utils.ParseMapSliceToJsonStr(people)
	if parseErr != nil {
		errorRes := fmt.Sprintf("Error while building JSON response: %v", parseErr)
		http.Error(w, errorRes, 500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, respContent)
}

func CreatePersonHandler(w http.ResponseWriter, req *http.Request) {
	bodyBytes, bodyErr := io.ReadAll(req.Body)
	if bodyErr != nil {
		errorRes := fmt.Sprintf("Error reading request body: %v", bodyErr)
		http.Error(w, errorRes, http.StatusBadRequest)
		return
	}
	body := string(bodyBytes)
	data, jsonErr := utils.ParseJsonStrToMap(body)
	if jsonErr != nil {
		errorRes := fmt.Sprintf("Error parsing body to map: %v", jsonErr)
		http.Error(w, errorRes, 500)
		return
	}
	writeErr := application.DataSource.SavePerson(data)
	if writeErr != nil {
		errorRes := fmt.Sprintf("Failed to write data to database: %v", writeErr)
		http.Error(w, errorRes, 500)
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
