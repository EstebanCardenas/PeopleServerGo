package handlers

import (
	"fmt"
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
	data, err := utils.GetRequestBodyAsMap(req.Body)
	if err != nil {
		errorRes := fmt.Sprintf("Error parsing request body: %v", err)
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
