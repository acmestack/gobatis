/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description:
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
