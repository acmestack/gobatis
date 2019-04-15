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
    TestTable gobatis.TableName "test_table"
    Username  string            `xfield:"username"`
    Password  string            `xfield:"password"`
}

func TestRunner(t *testing.T) {
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
    mgr.Select("queryTest").Params(1, 2).Result(&testV)
}
