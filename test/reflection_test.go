/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package test

import (
    "fmt"
    "github.com/xfali/GoBatis"
    "github.com/xfali/GoBatis/reflection"
    "testing"
)

type TestStruct1 struct {
    TestTable GoBatis.TableName "test_table"
    Username  string            `xfield:"username"`
    Password  string            `xfield:"password"`
}

func TestReflection1(t *testing.T) {
    test := TestStruct1{ "table", "abc", "123"}
    table, _ := reflection.GetTableInfo(&test)
    printTableInfo(table)
}

type TestStruct2 struct {
    TestTable GoBatis.TableName
    Username  string            `xfield:"-"`
    Password  string            `-`
}

func TestReflection2(t *testing.T) {
    test := TestStruct2{ "table", "abc", "123"}
    table, _ := reflection.GetTableInfo(&test)
    printTableInfo(table)
}

func printTableInfo(table *reflection.TableInfo) {
    fmt.Printf("table name is %s\n", table.Name)
    for _, v := range  table.Fields {
        fmt.Printf("field : %s, value %s\n", v.Name, v.Value)
    }

    for k, v := range  table.TypeMap {
        fmt.Printf("origin : %s, map %s\n", k, v)
    }
}
