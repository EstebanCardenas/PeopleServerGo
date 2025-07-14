package main

type DataSource interface {
	GetPeople() []map[string]any
	SavePerson(map[string]any)
}
