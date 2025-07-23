package datasource

import "fmt"

type PersonNotFoundError struct {
	id int
}

func (err PersonNotFoundError) Error() string {
	return fmt.Sprintf("Could not find person with id %d", err.id)
}
