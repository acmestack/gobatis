/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description:
 */

package test

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
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

func (t *TestTable) String() string {
	return fmt.Sprintf("id: %d, u: %s, p: %s", t.Id, t.Username, t.Password)
}

func connect() factory.Factory {
	return gobatis.NewFactory(
		gobatis.SetMaxConn(100),
		gobatis.SetMaxIdleConn(50),
		gobatis.SetDataSource(&datasource.PostgreDataSource{
			Host:     "localhost",
			Port:     5432,
			DBName:   "testdb",
			Username: "test",
			Password: "test",
		}))
}

func initTest(t *testing.T) (err error) {
	sql_table := "CREATE TABLE IF NOT EXISTS test_table (" +
		"id serial NOT NULL,"+
		"username varchar(255) DEFAULT NULL," +
		"password varchar(255) DEFAULT NULL," +
		"createTime timestamp DEFAULT NULL,"+
		"PRIMARY KEY (id)" +
	")"

	db, err := sql.Open("postgres", "host=localhost port=5432 user=test password=test dbname=testdb sslmode=disable")
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

func TestPostgres(t *testing.T) {
	initTest(t)
	var testV TestTable

	mgr := gobatis.NewSessionManager(connect())
	sess := mgr.NewSession()
	count := 0
	errI := sess.Insert("INSERT INTO test_table(username, password) VALUES('user1', 'pw1')").Param().Result(&count)
	if errI != nil {
		t.Log(errI)
	}
	testV = TestTable{}
	sess.Select("SELECT * FROM test_table").Param().Result(&testV)
	t.Logf("after insert: %d %v\n", count, testV)
	if testV.Username != "user1" || testV.Password != "pw1" {
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

func TestPostgresTx(t *testing.T) {
	initTest(t)
	var testV TestTable

	mgr := gobatis.NewSessionManager(connect())
	session := mgr.NewSession()
	session.Tx(func(sess *gobatis.Session) error {
		count := 0
		errI := sess.Insert("INSERT INTO test_table(username, password) VALUES('user1', 'pw1')").Param().Result(&count)
		if errI != nil {
			t.Log(errI)
		}
		testV = TestTable{}
		sess.Select("SELECT * FROM test_table").Param().Result(&testV)
		t.Logf("after insert: %d %v\n", count, testV)
		if testV.Username != "user1" || testV.Password != "pw1" {
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
			t.Log(err)
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
		session.Insert("insert_id").Param(TestTable{Id: 40, Username: "user", Password: "pw"}).Result(&ret)
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
		session.Insert("insert_id").Param(TestTable{Id: 40, Username: "user", Password: "pw"}).Result(&ret)
		t.Logf("ret %d\n", ret)
		session.Select("select_id").Param(TestTable{Id: 40}).Result(&retVar)
		t.Logf("data: %v", retVar)
		panic("rollback! ")
		//commit
		return nil
	})
}

func TestSession(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		mgr := gobatis.NewSessionManager(connect())
		var ret []TestTable
		mgr.NewSession().Select("SELECT * FROM test_table WHERE id = #{0}").Param(2).Result(&ret)
		t.Log(ret)
	})

	t.Run("template", func(t *testing.T) {
		mgr := gobatis.NewSessionManager(connect())
		mgr.SetParserFactory(gobatis.TemplateParserFactory)
		var ret []TestTable
		mgr.NewSession().Select("SELECT * FROM test_table WHERE id = {{.}}").Param(2).Result(&ret)
		t.Log(ret)
	})

	t.Run("template arg", func(t *testing.T) {
		mgr := gobatis.NewSessionManager(connect())
		mgr.SetParserFactory(gobatis.TemplateParserFactory)
		var ret []TestTable
		mgr.NewSession().Select("SELECT * FROM test_table WHERE id = {{arg .}}").Param(2).Result(&ret)
		t.Log(ret)
	})

	t.Run("template 2param", func(t *testing.T) {
		mgr := gobatis.NewSessionManager(connect())
		mgr.SetParserFactory(gobatis.TemplateParserFactory)
		var ret []TestTable
		mgr.NewSession().Select("SELECT * FROM test_table WHERE id = {{arg (index . 0)}} AND username = {{arg (index . 1)}}").Param(2, "user").Result(&ret)
		t.Log(ret)
	})
}
