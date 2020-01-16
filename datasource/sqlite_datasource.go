// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description:

package datasource

//import _ "github.com/mattn/go-sqlite3"

type SqliteDataSource struct {
	Path string
}

func (ds *SqliteDataSource) DriverName() string {
	return "sqlite3"
}

func (ds *SqliteDataSource) DriverInfo() string {
	return ds.Path
}
