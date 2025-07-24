package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"simple_server/application"
	ds "simple_server/datasource"
	"simple_server/handlers"
	"strconv"
	"strings"
)

func registerHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/people", handlers.PeopleHandler)
	mux.HandleFunc("/people/{id}", handlers.PeopleDetailHandler)
}

func initSqlDataSource() {
	dataSource := ds.NewSqlDataSource()
	err := dataSource.InitDb()
	if err != nil {
		log.Fatal(err)
	}
	application.SetDataSource(dataSource)
}

func initCsvDataSource() {
	dataSource := &ds.CsvDataSource{}
	err := dataSource.InitFile()
	if err != nil {
		log.Fatal(err)
	}
	application.SetDataSource(dataSource)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Select datasource:")
	fmt.Println("  1. SQL")
	fmt.Println("  2. CSV")
	scanner.Scan()
	op, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil {
		fmt.Println("Error: option must be a number")
		return
	}
	switch op {
	case 1:
		initSqlDataSource()
	case 2:
		initCsvDataSource()
	default:
		fmt.Println("Invalid option")
		return
	}

	mux := http.NewServeMux()
	registerHandlers(mux)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
