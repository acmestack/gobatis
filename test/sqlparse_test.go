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
    "github.com/xfali/gobatis/parsing/sqlparser"
    "github.com/xfali/gobatis/reflection"
    "testing"
    "time"
)

func TestSqlParser1(t *testing.T) {
    sqlStr := "SELECT * from xxx WHERE id = #{id}, name = #{name}"
    ret, _ := sqlparser.SimpleParse(sqlStr)
    t.Log(ret.String())
    if ret.Action != sqlparser.SELECT {
        t.Fail()
    }

    if ret.Vars[0] != "id" {
        t.Fail()
    }

    if ret.Vars[1] != "name" {
        t.Fail()
    }
}

func TestSqlParser2(t *testing.T) {
    sqlStr := "SELECT * from xxx WHERE id = #{id"
    _, err := sqlparser.SimpleParse(sqlStr)
    if err == nil {
        t.Fail()
    } else {
        t.Log(err)
    }
}

func TestSqlParser3(t *testing.T) {
    sqlStr := "SELECT * from xxx WHERE id = #{id, name = #{name}"
    _, err := sqlparser.SimpleParse(sqlStr)
    if err == nil {
        t.Fail()
    } else {
        t.Log(err)
    }
}

func TestSqlParser4(t *testing.T) {
    sqlStr := "DELETE from xxx WHERE id = #{id}, name = #{name}"
    ret, _ := sqlparser.SimpleParse(sqlStr)
    t.Log(ret.String())
    if ret.Action != sqlparser.DELETE {
        t.Fail()
    }

    if ret.Vars[0] != "id" {
        t.Fail()
    }

    if ret.Vars[1] != "name" {
        t.Fail()
    }
}

func TestSqlParserWithParams1(t *testing.T) {
    sqlStr := "DELETE from xxx WHERE id = #{0}, name = #{1}, id = #{0}"
    ret, _ := sqlparser.ParseWithParams(sqlStr, 123, "abc")
    t.Log(ret.String())
    if ret.Action != sqlparser.DELETE {
        t.Fail()
    }

    if ret.Vars[0] != "0" {
        t.Fail()
    }

    if ret.Vars[1] != "1" {
        t.Fail()
    }
}

func TestSqlParserWithParams2(t *testing.T) {
    sqlStr := "DELETE from ${2} WHERE id = ${0}, name = #{1}, id = #{0}"
    ret, _ := sqlparser.ParseWithParams(sqlStr, 123, "abc", "test_table")
    t.Log(ret.String())
    if ret.Action != sqlparser.DELETE {
        t.Fail()
    }

    if ret.Vars[0] != "2" {
        t.Fail()
    }

    if ret.Vars[1] != "0" {
        t.Fail()
    }
}

func TestSqlParserWithParams3(t *testing.T) {
    sqlStr := "SELECT from ${2} WHERE id = ${0, name = #{1}, id = #{0}"
    _, err := sqlparser.ParseWithParams(sqlStr, 123, "abc", "test_table")
    if err == nil {
        t.Fail()
    } else {
        t.Log(err)
    }
}

func TestSqlParserWithParams4(t *testing.T) {
    sqlStr := "SELECT from ${2} WHERE id = ${0}, name = #{1}, id = #{0}"
    _, err := sqlparser.ParseWithParams(sqlStr, 123, "abc")
    if err == nil {
        t.Fail()
    } else {
        t.Log(err)
    }
}

func TestSqlParserWithParamMap1(t *testing.T) {
    sqlStr := "SELECT * from ${tablename} WHERE id = ${id}, name = #{name}, id = #{id}"
    params := map[string]interface{}{
        "tablename": "test_table",
        "id":        123,
        "name":      "test_name",
    }
    ret, _ := sqlparser.ParseWithParamMap(sqlStr, params)
    t.Log(ret.String())
    if ret.Action != sqlparser.SELECT {
        t.Fail()
    }
}

func TestSqlParserWithParamMap2(t *testing.T) {
    sqlStr := "SELECT * from ${tablename} WHERE id = ${id, name = #{name}, id = #{id}"
    params := map[string]interface{}{
        "tablename": "test_table",
        "id":        123,
        "name":      "test_name",
    }
    _, err := sqlparser.ParseWithParamMap(sqlStr, params)
    if err == nil {
        t.Fail()
    } else {
        t.Log(err)
    }
}

func TestSqlParserWithParamMap3(t *testing.T) {
    sqlStr := "SELECT * from ${tablename} WHERE id = ${id}, name = #{name}, id = #{id}"
    params := map[string]interface{}{
        "tablename": "test_table",
        "id":        123,
        //"name" : "test_name",
    }
    _, err := sqlparser.ParseWithParamMap(sqlStr, params)
    if err == nil {
        t.Fail()
    } else {
        t.Log(err)
    }
}

type TestSqlParserStruct struct {
    TestTable gobatis.ModelName "test_table"
    Id        int               `xfield:"id"`
    Name      string            `xfield:"name"`
}

func TestSqlParserWithParamMap4(t *testing.T) {
    sqlStr := "SELECT * from ${tablename} WHERE id = ${id}, name = #{name}, id = #{id}"
    paramVar := TestSqlParserStruct{
        Id:   123,
        Name: "test_name",
    }
    ti, _ := reflection.GetTableInfo(&paramVar)
    params := ti.MapValue()
    params["tablename"] = ti.Name

    ret, _ := sqlparser.ParseWithParamMap(sqlStr, params)
    t.Log(ret.String())
    if ret.Action != sqlparser.SELECT {
        t.Fail()
    }
}

func TestSqlParserWithTime1(t *testing.T) {
    sqlStr := "SELECT * from test_table WHERE time = ${0}"
    ret, _ := sqlparser.ParseWithParams(sqlStr, time.Time{})
    t.Log(ret.String())
    if ret.Action != sqlparser.SELECT {
        t.Fail()
    }
}

func TestSqlParserWithTime2(t *testing.T) {
    sqlStr := "SELECT * from test_table WHERE time > #{0}"
    ret, _ := sqlparser.ParseWithParams(sqlStr, time.Time{})
    t.Log(ret.String())
    if ret.Action != sqlparser.SELECT {
        t.Fail()
    }
}
