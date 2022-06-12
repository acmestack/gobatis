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

package test

import (
	//_ "github.com/go-sql-driver/mysql"
	"github.com/xfali/gobatis"
	"github.com/xfali/gobatis/datasource"
	"github.com/xfali/gobatis/factory"
	"testing"
)

type TestTable struct {
	TestTable gobatis.ModelName "test_table"
	Id        int64             `xfield:"id"`
	Username  string            `xfield:"username"`
	Password  string            `xfield:"password"`
}

func connect() factory.Factory {
	return gobatis.NewFactory(
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
}

func TestSelectgobatis(t *testing.T) {
	var testV TestTable

	mgr := gobatis.NewSessionManager(connect())
	gobatis.RegisterSql("queryTest", "select * from test_table where id = #{0}")
	gobatis.RegisterModel(&testV)
	mgr.NewSession().Select("queryTest").Param(1, 200, 300).Result(&testV)

	t.Logf("%v %v", testV.Username, testV.Password)
}

func TestSelectgobatis2(t *testing.T) {
	var testV TestTable

	mgr := gobatis.NewSessionManager(connect())
	gobatis.RegisterSql("queryTest", "select * from test_table limit 10")
	gobatis.RegisterModel(&testV)
	testList := []TestTable{}
	mgr.NewSession().Select("queryTest").Param().Result(&testList)

	for _, v := range testList {
		t.Logf("%v %v", v.Username, v.Password)
	}
}

func TestSelectgobatis3(t *testing.T) {
	var testV TestTable

	mgr := gobatis.NewSessionManager(connect())
	gobatis.RegisterSql("queryTest", "select count(*) from test_table")
	gobatis.RegisterModel(&testV)
	i := 0
	mgr.NewSession().Select("queryTest").Param().Result(&i)
}

func TestInsertgobatis(t *testing.T) {
	var testV TestTable

	mgr := gobatis.NewSessionManager(connect())
	gobatis.RegisterSql("insertTest", "insert into test_table (username, password) value(#{username}, #{password})")
	gobatis.RegisterModel(&testV)
	testV.Username = "test_user"
	testV.Password = "test_pw"
	i := 0
	mgr.NewSession().Insert("insertTest").Param(testV).Result(&i)
	t.Logf("insert %d\n", i)
}

func TestUpdategobatis(t *testing.T) {
	var testV TestTable

	mgr := gobatis.NewSessionManager(connect())
	gobatis.RegisterSql("updateTest", "update test_table set username = #{username}, password = #{password} where id = 1")
	gobatis.RegisterModel(&testV)
	testV.Username = "test_user"
	testV.Password = "test_pw"
	i := 0
	mgr.NewSession().Update("updateTest").Param(testV).Result(&i)
	t.Logf("update %d\n", i)
}

func TestDeletegobatis(t *testing.T) {
	var testV TestTable

	mgr := gobatis.NewSessionManager(connect())
	gobatis.RegisterSql("deleteTest", "delete from test_table where username = #{username}, password = #{password}")
	gobatis.RegisterModel(&testV)
	testV.Username = "test_user"
	testV.Password = "test_pw"
	i := 0
	mgr.NewSession().Delete("deleteTest").Param(testV).Result(&i)
	t.Logf("delete %d\n", i)
}

func TestDynamicSelectgobatis(t *testing.T) {
	var testV TestTable

	mgr := gobatis.NewSessionManager(connect())
	gobatis.RegisterSql("deleteTest", `select id from test_table 
        <where> 
            <if test="">
            username = #{username}, password = #{password}
        </where>`)
	gobatis.RegisterModel(&testV)
	testV.Username = "test_user"
	testV.Password = "test_pw"
	i := 0
	mgr.NewSession().Delete("deleteTest").Param(testV).Result(&i)
	t.Logf("delete %d\n", i)
}

func TestTx1(t *testing.T) {
	mgr := gobatis.NewSessionManager(connect())
	testV := TestTable{
		Username: "testuser",
		Password: "testpw",
	}
	gobatis.RegisterModel(&testV)
	gobatis.RegisterSql("insert_id", "insert into test_table (username, password) value (#{username}, #{password})")
	gobatis.RegisterSql("select_id", "select * from test_table")

	var testList []TestTable

	mgr.NewSession().Tx(func(session *gobatis.Session) error {
		ret := 0
		err := session.Insert("insert_id").Param(testV).Result(&ret)
		if err != nil {
			return err
		}
		t.Logf("ret %d\n", ret)
		session.Select("select_id").Param().Result(&testList)
		for _, v := range testList {
			t.Logf("data: %v", v)
		}
		//commit
		return nil
	})
}

func TestTx2(t *testing.T) {
	mgr := gobatis.NewSessionManager(connect())
	testV := TestTable{
		Username: "testuser",
		Password: "testpw",
	}
	gobatis.RegisterModel(&testV)
	gobatis.RegisterSql("insert_id", "insert into test_table (username, password) value (#{username}, #{password})")
	gobatis.RegisterSql("select_id", "select * from test_table")

	var testList []TestTable

	mgr.NewSession().Tx(func(session *gobatis.Session) error {
		ret := 0
		err := session.Insert("insert_id").Param(testV).Result(&ret)
		if err != nil {
			t.Logf("error %v\n", err)
			return err
		}
		t.Logf("ret %d\n", ret)
		session.Select("select_id").Param().Result(&testList)
		for _, v := range testList {
			t.Logf("data: %v", v)
		}
		//rollback
		panic("ROLLBACK panic!!")

		//rollback
		return nil
	})
}
