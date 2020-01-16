// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description:

package datasource

import "fmt"

//import _ "github.com/lib/pq"

type PostgreDataSource struct {
	Host     string
	Port     int
	DBName   string
	Username string
	Password string
	SslMode  string
}

func (ds *PostgreDataSource) DriverName() string {
	return "postgres"
}

func (ds *PostgreDataSource) DriverInfo() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", ds.Host, ds.Port, ds.Username, ds.Password, ds.DBName, ds.SslMode)
}
