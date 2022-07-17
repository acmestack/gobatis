/*
 * Copyright (c) 2022, AcmeStack
 * All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
	if ds.SslMode == "" {
		ds.SslMode = "disable"
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", ds.Host, ds.Port, ds.Username, ds.Password, ds.DBName, ds.SslMode)
}
