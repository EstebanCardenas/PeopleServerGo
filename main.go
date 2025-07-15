package main

import (
	"log"
	"net/http"
	"simple_server/application"
	"simple_server/datasource"
	"simple_server/handlers"
)

func registerHandlers() {
	http.HandleFunc("/people", handlers.PeopleHandler)
}

func main() {
	registerHandlers()
	application.SetDataSource(&datasource.CsvDataSource{})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
