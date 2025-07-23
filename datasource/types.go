package datasource

type PeopleDataSource interface {
	GetPeople() ([]map[string]any, error)
	SavePerson(map[string]any) error
	DeletePerson(id int) error
	UpdatePerson(int, map[string]any) error
}
