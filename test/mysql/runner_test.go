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
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/acmestack/gobatis"
	"github.com/acmestack/gobatis/datasource"
	"github.com/acmestack/gobatis/factory"
	_ "github.com/go-sql-driver/mysql"
)

type TestTable struct {
	TestTable gobatis.TableName "test_table"
	Id        int64             `column:"id"`
	Username  string            `column:"username"`
	Password  string            `column:"password"`
}

func (t *TestTable) String() string {
	return fmt.Sprintf("id: %d, u: %s, p: %s", t.Id, t.Username, t.Password)
}

func connect() factory.Factory {
	return gobatis.NewFactory(
		gobatis.SetMaxConn(100),
		gobatis.SetMaxIdleConn(50),
		gobatis.SetDataSource(&datasource.MysqlDataSource{
			Host:     "localhost",
			Port:     3306,
			DBName:   "test",
			Username: "test",
			Password: "test",
			Charset:  "utf8",
		}))
}

func initTest(t *testing.T) (err error) {
	sql_table := "CREATE TABLE IF NOT EXISTS `test_table` (" +
		"`id` int(11) NOT NULL AUTO_INCREMENT," +
		"`username` varchar(255) DEFAULT NULL," +
		"`password` varchar(255) DEFAULT NULL," +
		"`createTime` datetime DEFAULT NULL," +
		"PRIMARY KEY (`id`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;"

	db, err := sql.Open("mysql", "test:test@tcp(localhost:3306)/test?charset=utf8")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(sql_table)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec("DELETE FROM test_table")
	if err != nil {
		t.Fatal(err)
	}
	return nil
}

func TestMysql(t *testing.T) {
	initTest(t)
	var testV TestTable

	mgr := gobatis.NewSessionManager(connect())
	sess := mgr.NewSession()
	count := 0
	errI := sess.Insert("INSERT INTO test_table(username, password) VALUES('user1', 'pw1')").Param().Result(&count)
	if errI != nil {
		t.Fatal(errI)
	}
	testV = TestTable{}
	sess.Select("SELECT * FROM test_table").Param().Result(&testV)
	t.Logf("after insert: %d %v\n", count, testV)
	if count != 1 || testV.Username != "user1" || testV.Password != "pw1" {
		t.Fail()
	}

	errU := sess.Update("Update test_table SET username = 'user2', password = 'pw2'").Param().Result(&count)
	if errU != nil {
		t.Fatal(errU)
	}
	testV = TestTable{}
	sess.Select("SELECT * FROM test_table").Param().Result(&testV)
	t.Logf("after update: %d %v\n", count, testV)
	if count != 1 || testV.Username != "user2" || testV.Password != "pw2" {
		t.Fail()
	}

	errD := sess.Delete("DELETE FROM test_table").Param().Result(&count)
	if errD != nil {
		t.Fatal(errD)
	}
	testV = TestTable{}
	sess.Select("SELECT * FROM test_table").Param().Result(&testV)
	t.Logf("after delete: %d %v\n", count, testV)
	if count != 1 || testV.Username != "" || testV.Password != "" {
		t.Fail()
	}

	t.Logf("%v %v", testV.Username, testV.Password)
}

func TestMysqlTx(t *testing.T) {
	initTest(t)
	var testV TestTable

	mgr := gobatis.NewSessionManager(connect())
	session := mgr.NewSession()
	session.Tx(func(sess *gobatis.Session) error {
		count := 0
		errI := sess.Insert("INSERT INTO test_table(username, password) VALUES('user1', 'pw1')").Param().Result(&count)
		if errI != nil {
			t.Fatal(errI)
		}
		testV = TestTable{}
		sess.Select("SELECT * FROM test_table").Param().Result(&testV)
		t.Logf("after insert: %d %v\n", count, testV)
		if count != 1 || testV.Username != "user1" || testV.Password != "pw1" {
			t.Fail()
		}

		errU := sess.Update("Update test_table SET username = 'user2', password = 'pw2'").Param().Result(&count)
		if errU != nil {
			t.Fatal(errU)
		}
		testV = TestTable{}
		sess.Select("SELECT * FROM test_table").Param().Result(&testV)
		t.Logf("after update: %d %v\n", count, testV)
		if count != 1 || testV.Username != "user2" || testV.Password != "pw2" {
			t.Fail()
		}

		errD := sess.Delete("DELETE FROM test_table").Param().Result(&count)
		if errD != nil {
			t.Fatal(errD)
		}
		testV = TestTable{}
		sess.Select("SELECT * FROM test_table").Param().Result(&testV)
		t.Logf("after delete: %d %v\n", count, testV)
		if count != 1 || testV.Username != "" || testV.Password != "" {
			t.Fail()
		}
		return nil
	})
}

func TestTx1(t *testing.T) {
	initTest(t)
	mgr := gobatis.NewSessionManager(connect())
	testV := TestTable{}
	gobatis.RegisterModel(&testV)
	gobatis.RegisterSql("insert_id", "insert into test_table (id, username, password) "+
		"values (#{TestTable.id}, #{TestTable.username}, #{TestTable.password})")
	gobatis.RegisterSql("select_id", "select * from test_table where id = #{TestTable.id}")

	var retVar TestTable

	mgr.NewSession().Tx(func(session *gobatis.Session) error {
		ret := 0
		err := session.Insert("insert_id").Param(TestTable{Id: 40, Username: "user", Password: "pw"}).Result(&ret)
		if err != nil {
			return err
		}
		t.Logf("ret %d\n", ret)
		session.Select("select_id").Param(TestTable{Id: 40}).Result(&retVar)
		t.Logf("data: %v", retVar)
		//commit
		return nil
	})

	retVar = TestTable{}
	mgr.NewSession().Select("select_id").Param(TestTable{Id: 40}).Result(&retVar)
	t.Logf("out tx: %v", retVar)
	if retVar.Username != "user" || retVar.Password != "pw" {
		t.Fail()
	}
}

func TestTx2(t *testing.T) {
	initTest(t)
	mgr := gobatis.NewSessionManager(connect())
	testV := TestTable{}
	gobatis.RegisterModel(&testV)
	gobatis.RegisterSql("insert_id", "insert into test_table (id, username, password) "+
		"values (#{TestTable.id}, #{TestTable.username}, #{TestTable.password})")
	gobatis.RegisterSql("select_id", "select * from test_table where id = #{TestTable.id}")

	var retVar TestTable

	mgr.NewSession().Tx(func(session *gobatis.Session) error {
		ret := 0
		err := session.Insert("insert_id").Param(TestTable{Id: 40, Username: "user", Password: "pw"}).Result(&ret)
		if err != nil {
			return err
		}
		t.Logf("ret %d\n", ret)
		session.Select("select_id").Param(TestTable{Id: 40}).Result(&retVar)
		t.Logf("data: %v", retVar)
		//commit
		return errors.New("rollback! ")
	})

	retVar = TestTable{}
	mgr.NewSession().Select("select_id").Param(TestTable{Id: 40}).Result(&retVar)
	t.Logf("out tx: %v", retVar)
	if retVar.Username == "user" || retVar.Password == "pw" {
		t.Fail()
	}
}

func TestTx3(t *testing.T) {
	initTest(t)
	mgr := gobatis.NewSessionManager(connect())
	testV := TestTable{}
	gobatis.RegisterModel(&testV)
	gobatis.RegisterSql("insert_id", "insert into test_table (id, username, password) "+
		"values (#{TestTable.id}, #{TestTable.username}, #{TestTable.password})")
	gobatis.RegisterSql("select_id", "select * from test_table where id = #{TestTable.id}")

	var retVar TestTable

	defer func() {
		if r := recover(); r != nil {
			retVar = TestTable{}
			mgr.NewSession().Select("select_id").Param(TestTable{Id: 40}).Result(&retVar)
			t.Logf("out tx: %v", retVar)
			if retVar.Username == "user" || retVar.Password == "pw" {
				t.Fail()
			}
		}
	}()
	mgr.NewSession().Tx(func(session *gobatis.Session) error {
		ret := 0
		err := session.Insert("insert_id").Param(TestTable{Id: 40, Username: "user", Password: "pw"}).Result(&ret)
		if err != nil {
			return err
		}
		t.Logf("ret %d\n", ret)
		session.Select("select_id").Param(TestTable{Id: 40}).Result(&retVar)
		t.Logf("data: %v", retVar)
		panic("rollback! ")
		//commit
		return nil
	})
}
