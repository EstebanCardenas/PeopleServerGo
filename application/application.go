package application

import (
	"simple_server/datasource"
)

var DataSource datasource.PeopleDataSource

func SetDataSource(ds datasource.PeopleDataSource) {
	DataSource = ds
}
