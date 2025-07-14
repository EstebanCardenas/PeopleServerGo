package main

import (
	"log"
	"net/http"
)

func registerHandlers() {
	http.HandleFunc("/people", PeopleHandler)
}

func main() {
	registerHandlers()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
