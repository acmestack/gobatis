/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package test

import (
    "github.com/xfali/gobatis"
    "github.com/xfali/gobatis/config"
    "github.com/xfali/gobatis/factory"
    "github.com/xfali/gobatis/logging"
    "github.com/xfali/gobatis/session/runner"
    "testing"
)

type TestTable struct {
    TestTable gobatis.ModelName "test_table"
    Username  string            `xfield:"username"`
    Password  string            `xfield:"password"`
}

func TestSelectRunner(t *testing.T) {
    var testV TestTable
    fac := factory.DefaultFactory{
        Host:     "localhost",
        Port:     3306,
        DBName:   "test",
        Username: "root",
        Password: "123",
        Charset:  "utf8",

        MaxConn:     1000,
        MaxIdleConn: 500,

        Log: logging.DefaultLogf,
    }
    fac.Init()
    mgr := runner.NewSqlManager(&fac)
    mgr.RegisterSql("queryTest", "select * from test_table where id = #{0}")
    config.RegisterModel(&testV)
    mgr.Select("queryTest").Params(1).Result(&testV)

    t.Logf("%v %v", testV.Username, testV.Password)
}

func TestSelectRunner2(t *testing.T) {
    var testV TestTable
    fac := factory.DefaultFactory{
        Host:     "localhost",
        Port:     3306,
        DBName:   "test",
        Username: "root",
        Password: "123",
        Charset:  "utf8",

        MaxConn:     1000,
        MaxIdleConn: 500,

        Log: logging.DefaultLogf,
    }
    fac.Init()
    mgr := runner.NewSqlManager(&fac)
    mgr.RegisterSql("queryTest", "select * from test_table limit 10")
    config.RegisterModel(&testV)
    testList := []TestTable{}
    mgr.Select("queryTest").Params().Result(&testList)

    for _, v := range testList {
        t.Logf("%v %v", v.Username, v.Password)
    }
}

func TestSelectRunner3(t *testing.T) {
    var testV TestTable
    fac := factory.DefaultFactory{
        Host:     "localhost",
        Port:     3306,
        DBName:   "test",
        Username: "root",
        Password: "123",
        Charset:  "utf8",

        MaxConn:     1000,
        MaxIdleConn: 500,

        Log: logging.DefaultLogf,
    }
    fac.Init()
    mgr := runner.NewSqlManager(&fac)
    mgr.RegisterSql("queryTest", "select count(*) from test_table")
    config.RegisterModel(&testV)
    i := 0
    mgr.Select("queryTest").Params().Result(&i)
}

func TestInsertRunner(t *testing.T) {
    var testV TestTable
    fac := factory.DefaultFactory{
        Host:     "localhost",
        Port:     3306,
        DBName:   "test",
        Username: "root",
        Password: "123",
        Charset:  "utf8",

        MaxConn:     1000,
        MaxIdleConn: 500,

        Log: logging.DefaultLogf,
    }
    fac.Init()
    mgr := runner.NewSqlManager(&fac)
    mgr.RegisterSql("insertTest", "insert into test_table (username, password) value(#{username}, #{password})")
    config.RegisterModel(&testV)
    testV.Username = "test_user"
    testV.Password = "test_pw"
    i := 0
    mgr.Insert("insertTest").ParamType(testV).Result(&i)
    t.Logf("insert %d\n", i)
}

func TestUpdateRunner(t *testing.T) {
    var testV TestTable
    fac := factory.DefaultFactory{
        Host:     "localhost",
        Port:     3306,
        DBName:   "test",
        Username: "root",
        Password: "123",
        Charset:  "utf8",

        MaxConn:     1000,
        MaxIdleConn: 500,

        Log: logging.DefaultLogf,
    }
    fac.Init()
    mgr := runner.NewSqlManager(&fac)
    mgr.RegisterSql("updateTest", "update test_table set username = #{username}, password = #{password} where id = 1")
    config.RegisterModel(&testV)
    testV.Username = "test_user"
    testV.Password = "test_pw"
    i := 0
    mgr.Update("updateTest").ParamType(testV).Result(&i)
    t.Logf("update %d\n", i)
}

func TestDeleteRunner(t *testing.T) {
    var testV TestTable
    fac := factory.DefaultFactory{
        Host:     "localhost",
        Port:     3306,
        DBName:   "test",
        Username: "root",
        Password: "123",
        Charset:  "utf8",

        MaxConn:     1000,
        MaxIdleConn: 500,

        Log: logging.DefaultLogf,
    }
    fac.Init()
    mgr := runner.NewSqlManager(&fac)
    mgr.RegisterSql("deleteTest", "delete from test_table where username = #{username}, password = #{password}")
    config.RegisterModel(&testV)
    testV.Username = "test_user"
    testV.Password = "test_pw"
    i := 0
    mgr.Delete("deleteTest").ParamType(testV).Result(&i)
    t.Logf("delete %d\n", i)
}
