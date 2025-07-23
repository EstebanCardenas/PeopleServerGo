package handlers

import (
	"fmt"
	"net/http"
	app "simple_server/application"
	ds "simple_server/datasource"
	"simple_server/utils"
	"strconv"
)

func UpdatePersonHandler(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(req.PathValue("id"))
	if err != nil {
		errStr := fmt.Sprintf("Could not convert id to integer: %v", err)
		http.Error(w, errStr, 500)
		return
	}
	reqBody, err := utils.GetRequestBodyAsMap(req.Body)
	if err != nil {
		errStr := fmt.Sprintf("Failed to parse request body: %v", err)
		http.Error(w, errStr, 500)
		return
	}
	err = app.DataSource.UpdatePerson(id, reqBody)
	if err != nil {
		errStr := fmt.Sprintf("Failed to update person: %v", err)
		http.Error(w, errStr, 500)
		return
	}

	fmt.Fprint(w, "Succesfully updated person")
}

func DeletePersonHandler(w http.ResponseWriter, req *http.Request) {
	idStr := req.PathValue("id")
	id, atoiErr := strconv.Atoi(idStr)
	if atoiErr != nil {
		errStr := fmt.Sprintf("Failed to convert id %v to integer: %v", idStr, atoiErr)
		http.Error(w, errStr, 500)
		return
	}
	delErr := app.DataSource.DeletePerson(id)
	if delErr != nil {
		errStr := fmt.Sprintf("Failed to delete person: %v", delErr)
		errCode := 500
		if _, ok := delErr.(ds.PersonNotFoundError); ok {
			errCode = 404
		}
		http.Error(w, errStr, errCode)
		return
	}

	fmt.Fprintf(w, "Succesfully deleted person with id: %d", id)
}

func PeopleDetailHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "PUT":
		UpdatePersonHandler(w, req)
	case "DELETE":
		DeletePersonHandler(w, req)
	default:
		http.Error(w, "Unsupported HTTP method", http.StatusMethodNotAllowed)
	}
}
