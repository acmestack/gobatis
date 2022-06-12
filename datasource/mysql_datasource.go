/*
 * Copyright (c) 2022, OpeningO
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

//import _ "github.com/go-sql-driver/mysql"

type MysqlDataSource struct {
	Host     string
	Port     int
	DBName   string
	Username string
	Password string
	Charset  string
}

func (ds *MysqlDataSource) DriverName() string {
	return "mysql"
}

func (ds *MysqlDataSource) DriverInfo() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", ds.Username, ds.Password, ds.Host, ds.Port, ds.DBName, ds.Charset)
}
