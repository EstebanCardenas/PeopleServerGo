package main

import (
	"log"
	"net/http"
	"simple_server/application"
	ds "simple_server/datasource"
	"simple_server/handlers"
)

func registerHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/people", handlers.PeopleHandler)
	mux.HandleFunc("/people/{id}", handlers.PeopleDetailHandler)
}

func main() {
	dataSource := ds.NewSqlDataSource()
	err := dataSource.InitDb()
	if err != nil {
		log.Fatal(err)
	}
	// dataSource.DeleteAllPeople()
	application.SetDataSource(dataSource)

	mux := http.NewServeMux()
	registerHandlers(mux)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
