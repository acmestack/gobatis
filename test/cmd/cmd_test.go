/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package test_package

import (
    "testing"
    "github.com/xfali/gobatis"
    "github.com/xfali/gobatis/factory"
    "github.com/xfali/gobatis/logging"
)

var sessionMgr *gobatis.SessionManager

func init() {
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
    fac.InitDB()
    sessionMgr = gobatis.NewSessionManager(&fac)
}

func TestSelect(t *testing.T) {
    proxy := New(sessionMgr)

    ret := proxy.SelectTestTable(TestTable{Username: "test_user"})

    for _, v := range ret {
        t.Logf("value: %v", v)
    }
}

func TestSelectCount(t *testing.T) {
    proxy := New(sessionMgr)

    //all count
    ret := proxy.SelectTestTable(TestTable{})

    t.Logf("count is %d", ret)
}

func TestInsert(t *testing.T) {
    proxy := New(sessionMgr)

    ret, id := proxy.InsertTestTable(TestTable{Username: "test_insert_user", Password: "test_pw"})

    t.Logf("insert ret is %d, id is %d", ret, id)
}

func TestUpdate(t *testing.T) {
    proxy := New(sessionMgr)

    ret := proxy.UpdateTestTable(TestTable{Id: 1, Username: "test_insert_user", Password: "update_pw"})

    t.Logf("update ret is %d", ret)
}

func TestDelete(t *testing.T) {
    proxy := New(sessionMgr)

    //Same as: DELETE FROM test_table WHERE username = 'test_insert_user'
    ret := proxy.DeleteTestTable(TestTable{Username: "test_insert_user"})

    t.Logf("delete ret is %d", ret)
}

func TestTxSuccess(t *testing.T) {
    proxy := New(sessionMgr)
    proxy.Tx(func(s *TestTableCallProxy) bool {
        s.InsertTestTable(TestTable{Username: "tx_user", Password: "tx_pw"})
        //select all
        ret := s.SelectTestTable(TestTable{})

        for _, v := range ret {
            t.Logf("value: %v", v)
        }
        //will commit
        return true
    })

    //select all
    ret := New(sessionMgr).SelectTestTable(TestTable{Username: "tx_user"})

    for _, v := range ret {
        t.Logf("value: %v", v)
    }
}


func TestTxFail(t *testing.T) {
    proxy := New(sessionMgr)
    proxy.Tx(func(s *TestTableCallProxy) bool {
        s.InsertTestTable(TestTable{Username: "tx_fail_user", Password: "tx_fail_pw"})
        //select all
        ret := s.SelectTestTable(TestTable{})

        for _, v := range ret {
            t.Logf("value: %v", v)
        }

        //will rollback
        //panic("test rollback")

        //will rollback too
        return false
    })

    //select all
    ret := New(sessionMgr).SelectTestTable(TestTable{Username: "tx_fail_user"})

    for _, v := range ret {
        t.Logf("value: %v", v)
    }
}
