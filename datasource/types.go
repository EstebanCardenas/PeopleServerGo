package datasource

type DataSource interface {
	GetPeople() ([]map[string]any, error)
	SavePerson(map[string]any) error
}
