package datasource

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type SqlDataSource struct {
	db *sql.DB
}

var personColumnFields = [...]string{
	"id", "name", "last_name", "profession", "age",
}

func NewSqlDataSource() *SqlDataSource {
	return &SqlDataSource{}
}

func (ds *SqlDataSource) CreateTable() error {
	createQuery := `
		CREATE TABLE IF NOT EXISTS people(
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			profession TEXT NOT NULL,
			age INTEGER NOT NULL
		)
	`
	_, err := ds.db.Exec(createQuery)
	if err != nil {
		return err
	}

	return nil
}

func (ds *SqlDataSource) InitDb() error {
	db, err := sql.Open("sqlite3", "file:data.db")
	if err != nil {
		return err
	}

	ds.db = db
	createErr := ds.CreateTable()
	if createErr != nil {
		return createErr
	}
	return nil
}

func (ds *SqlDataSource) DeleteAllPeople() error {
	_, err := ds.db.Exec("DELETE FROM people")
	if err != nil {
		return err
	}

	return nil
}

func (ds *SqlDataSource) GetPeople() ([]map[string]any, error) {
	getPeopleQuery := "SELECT * FROM people"
	rows, err := ds.db.Query(getPeopleQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	people := []map[string]any{}

	for rows.Next() {
		var (
			id         int
			name       string
			lastName   string
			profession string
			age        int
		)
		err := rows.Scan(
			&id, &name, &lastName, &profession, &age,
		)
		if err != nil {
			return nil, err
		}
		person := map[string]any{
			"id":         id,
			"name":       name,
			"lastName":   lastName,
			"profession": profession,
			"age":        age,
		}
		people = append(people, person)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return people, nil
}

func (ds *SqlDataSource) SavePerson(person map[string]any) error {
	query := `
		INSERT INTO people (name, last_name, profession, age)
		VALUES (?, ?, ?, ?)
	`
	_, err := ds.db.Exec(
		query,
		person["name"],
		person["lastName"],
		person["profession"],
		person["age"],
	)
	if err != nil {
		return err
	}

	return nil
}

func (ds *SqlDataSource) DeletePerson(id int) error {
	query := `
		DELETE FROM people
		WHERE id = ?
	`
	res, err := ds.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsDeleted, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsDeleted == 0 {
		return PersonNotFoundError{id}
	}

	return nil
}

func (ds *SqlDataSource) UpdatePerson(id int, data map[string]any) error {
	statements := []string{}
	values := []any{}
	for i := 1; i < len(personColumnFields); i++ {
		field := personFields[i]
		if value, ok := data[field]; ok {
			statements = append(statements, fmt.Sprintf("%v = ?", personColumnFields[i]))
			values = append(values, value)
		}
	}

	if len(statements) == 0 {
		return errors.New("no fields to update")
	}

	query := `
		UPDATE people
		SET %v
		WHERE id = ?
	`
	query = fmt.Sprintf(query, strings.Join(statements, ", "))
	res, err := ds.db.Exec(query, append(values, id)...)
	if err != nil {
		return err
	}
	rowsUpdated, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsUpdated == 0 {
		return PersonNotFoundError{id}
	}

	return nil
}
