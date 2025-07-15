package application

import (
	"simple_server/datasource"
)

var DataSource datasource.DataSource

func SetDataSource(ds datasource.DataSource) {
	DataSource = ds
}
