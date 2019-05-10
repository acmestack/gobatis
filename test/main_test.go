/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package test

import (
    "github.com/xfali/gobatis/config"
    "github.com/xfali/gobatis/factory"
    "github.com/xfali/gobatis/logging"
    "github.com/xfali/gobatis/session/runner"
    "testing"
)

func TestSql1(t *testing.T) {
    testV := TestTable{
    }
    config.RegisterModel(&testV)
    config.RegisterMapperData([]byte(main_xml))

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
    mgr := runner.NewSessionManager(&fac)

    ret := 0
    mgr.NewSession().Select("selectCount").Param(testV).Result(&ret)
    t.Logf("ret = %d\n", ret)

    dataList := []TestTable{}
    mgr.NewSession().Select("selectUser").Param(testV).Result(&dataList)
    for _, v := range dataList {
        t.Logf("select data: %v\n", v)
    }

    mgr.NewSession().Insert("insertUser").Param(testV).Result(&ret)
    t.Logf("insert ret = %d\n", ret)

    mgr.NewSession().Update("updateUser").Param(testV).Result(&ret)
    t.Logf("update ret = %d\n", ret)

    mgr.NewSession().Delete("deleteUser").Param(testV).Result(&ret)
    t.Logf("deleteUser ret = %d\n", ret)

    mgr.NewSession().Tx(func(session *runner.Session) bool {
        session.Insert("insertUser").Param(testV).Result(&ret)
        //commit
        return true
    })
}

func TestMain2(t *testing.T) {
    testV := TestTable{
    }
    config.RegisterModel(&testV)
    config.RegisterMapperData([]byte(main_xml))

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
    mgr := runner.NewSessionManager(&fac)

    mgr.NewSession().Select("select * from test_table where id = ${0}").Param(100).Result(&testV)
    t.Logf("ret = %v\n", testV)
}
