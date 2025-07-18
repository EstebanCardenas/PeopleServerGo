package main

import (
	"log"
	"net/http"
	"simple_server/application"
	ds "simple_server/datasource"
	"simple_server/handlers"
)

func registerHandlers() {
	http.HandleFunc("/people", handlers.PeopleHandler)
}

func main() {
	dataSource := ds.NewSqlDataSource()
	err := dataSource.InitDb()
	if err != nil {
		log.Fatal(err)
	}
	// dataSource.DeleteAllPeople()
	application.SetDataSource(dataSource)

	registerHandlers()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
