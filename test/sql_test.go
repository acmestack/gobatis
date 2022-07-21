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

package test

import (
	"github.com/acmestack/gobatis"
	"github.com/acmestack/gobatis/datasource"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

type SqlTest struct {
	Id       int64  `column:"id"`
	Username string `column:"username"`
	Password string `column:"password"`
}

var sessionMgr *gobatis.SessionManager

func init() {
	fac := gobatis.NewFactory(
		gobatis.SetMaxConn(100),
		gobatis.SetMaxIdleConn(50),
		gobatis.SetDataSource(&datasource.MysqlDataSource{
			Host:     "localhost",
			Port:     3306,
			DBName:   "test",
			Username: "root",
			Password: "123",
			Charset:  "utf8",
		}))
	var testV TestTable
	gobatis.RegisterModel(&testV)
	sessionMgr = gobatis.NewSessionManager(fac)
}

func TestSelectWithSimpleType(t *testing.T) {
	sql := `SELECT username, password FROM test_table WHERE id = #{0}`
	ret := SqlTest{}
	sessionMgr.NewSession().Select(sql).Param(1).Result(&ret)
	t.Logf("%v %v", ret.Username, ret.Password)
}

func TestSelectWithMap(t *testing.T) {
	sql := `SELECT username, password FROM test_table WHERE id = #{id}`
	ret := SqlTest{}
	sessionMgr.NewSession().Select(sql).Param(map[string]interface{}{"id": 1}).Result(&ret)
	t.Logf("%v %v", ret.Username, ret.Password)
}

func TestSelectWithStruct(t *testing.T) {
	sql := `SELECT username, password FROM test_table WHERE id = #{SqlTest.id}`
	ret := SqlTest{}
	sessionMgr.NewSession().Select(sql).Param(SqlTest{Id: 1}).Result(&ret)
	t.Logf("%v %v", ret.Username, ret.Password)
}

func TestSelectWithComplex(t *testing.T) {
	sql := `SELECT username, password FROM test_table WHERE id = #{SqlTest.id} AND username = #{0} AND password = #{pw}`
	ret := SqlTest{}
	sessionMgr.NewSession().Select(sql).Param(SqlTest{Id: 1}, "test_user", map[string]interface{}{"pw": "test_pw"}).Result(&ret)
	t.Logf("%v %v", ret.Username, ret.Password)
}
