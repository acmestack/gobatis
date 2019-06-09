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
    "github.com/xfali/gobatis"
    "github.com/xfali/gobatis/reflection"
    "reflect"
    "testing"
    "time"
)

type TestStruct1 struct {
    TestTable gobatis.ModelName "test_table"
    Username  string            `xfield:"username"`
    Password  string            `xfield:"password"`
}

func TestReflection1(t *testing.T) {
    test := TestStruct1{"table", "abc", "123"}
    table, _ := reflection.GetStructInfo(&test)
    printTableInfo(table)
}

type TestStruct2 struct {
    TestTable gobatis.ModelName
    Username  string `xfield:"-"`
    Password  string `-`
}

func TestReflection2(t *testing.T) {
    test := TestStruct2{"table", "abc", "123"}
    table, _ := reflection.GetStructInfo(&test)
    printTableInfo(table)
}

func TestReflection3(t *testing.T) {
    var test int
    o, err := reflection.GetObjectInfo(&test)
    if err == nil {
        t.Log(o)
    } else {
        t.Fail()
    }
}

func TestReflection4(t *testing.T) {
    i := 10
    v := reflect.New(reflect.TypeOf(&i).Elem()).Elem()
    reflection.SetValue(v, i)
    t.Log(v.Interface())
}

func TestReflection4_2(t *testing.T) {
    var i complex128 = complex(1, 2)
    v := reflect.New(reflect.TypeOf(&i).Elem()).Elem()
    sv := fmt.Sprintf("%v", i)
    t.Logf("sv %v\n", sv)
    reflection.SetValue(v, []byte(sv))
    t.Log(v.Interface())
}

func TestReflection6(t *testing.T) {
    v := []TestStruct1{}
    t.Log(reflection.GetBeanClassName(v))

    v2 := 10
    t.Log(reflection.GetBeanClassName(v2))

    v3 := TestStruct2{}
    t.Log(reflection.GetBeanClassName(v3))
}

func TestReflection7(t *testing.T) {
    o := []TestStruct1{}

    rv := reflect.ValueOf(&o)
    rv = rv.Elem()
    rvv := rv
    i := TestStruct1{Username: "1", Password: "1"}
    rvv = reflect.Append(rvv, reflect.ValueOf(i))
    i = TestStruct1{Username: "2", Password: "2"}
    rvv = reflect.Append(rvv, reflect.ValueOf(i))
    rv.Set(rvv)

    for _, e := range o {
        t.Log(e)
    }
}

func TestReflectionParseEmpty(t *testing.T) {
    ret := reflection.ParseParams()
    for k, v := range ret {
        t.Logf("empty key : %s value : %v", k, v)
    }
}

func TestReflectionParseSimple(t *testing.T) {
    ret := reflection.ParseParams(1, "2", 1.3, time.Now())
    for k, v := range ret {
        t.Logf("simple key : %s value : %v", k, v)
    }
}

func TestReflectionParseMap(t *testing.T) {
    ret := reflection.ParseParams(map[string]interface{}{
        "mapkey1_int":    123,
        "mapkey2_string": "test",
        "mapkey3_float":  1.1,
        "mapkey4_time":   time.Now(),
    })
    if len(ret) == 0 {
        t.Fail()
    }
    for k, v := range ret {
        t.Logf("map key : %s value : %v", k, v)
    }
}

type testParseStruct struct {
    Name     gobatis.ModelName `parse_struct`
    Username string
    Password string
}

func TestReflectionParseStruct(t *testing.T) {
    ret := reflection.ParseParams(testParseStruct{
        "x",
        "user",
        "pw",
    })
    if len(ret) == 0 {
        t.Fail()
    }
    for k, v := range ret {
        t.Logf("struct key : %s value : %v", k, v)
    }
}

func TestReflectionParseComplex(t *testing.T) {
    ret := reflection.ParseParams(1, map[string]interface{}{
        "mapkey1_int":    123,
        "mapkey2_string": "test",
        "mapkey3_float":  1.1,
        "mapkey4_time":   time.Now(),
    }, "2", testParseStruct{
        Username: "user",
        Password: "pw",
    }, 1.3, time.Now())
    if len(ret) == 0 {
        t.Fail()
    }
    for k, v := range ret {
        t.Logf("complex key : %s value : %v", k, v)
    }
}

func TestReflectionParseSlice(t *testing.T) {
    ret := reflection.ParseParams([]int{1,2,3,4})
    if len(ret) == 0 {
        t.Fail()
    }
    for k, v := range ret {
        t.Logf("complex key : %s value : %v", k, v)
        elems := reflection.ParseSliceParamString(v.(string))
        for _, e := range elems {
            t.Logf("slice elem %v\n", e)
        }
    }
}

func TestSimpleTypeTime(t *testing.T) {
    ret := reflection.IsSimpleObject(time.Time{})
    if !ret {
        t.Fail()
    }
}

func TestSimpleTypeSliceByte(t *testing.T) {
    ret := reflection.IsSimpleObject([]byte{})
    if !ret {
        t.Fail()
    }

    ret = reflection.IsSimpleObject([]int{})
    if ret {
        t.Fail()
    }
}

func TestBeanClass(t *testing.T) {
    t.Log(reflection.GetBeanClassName(TestStruct2{}))
}

func TestTypeName(t *testing.T) {
    v := []byte{}
    t.Log(reflect.TypeOf(v).Elem().Name())

    var i interface{}
    i = v
    t.Logf("interface %s", reflect.TypeOf(i).String())

    m := map[string]interface{}{}
    t.Logf("map %s", reflect.TypeOf(m).String())

    st := TestTable{}
    t.Logf("struct %s", reflect.TypeOf(st).String())

    ptr := &v
    t.Logf("ptr %s", reflect.TypeOf(ptr).String())
}

func TestReflectSlice(t *testing.T) {
    o := []TestTable{}
    rt := reflect.TypeOf(o)
    //rv := reflect.ValueOf(o)

    rt = rt.Elem()
    t.Log(rt)
}

func TestReflectMap(t *testing.T) {
    o := map[string]int{}
    rt := reflect.TypeOf(o)
    //rv := reflect.ValueOf(o)
    t.Log(rt.Key())
    rt = rt.Elem()
    t.Log(rt)
    //t.Logf("rv valid: %v\n", rv.Elem())
}

func returnInterface(b interface{}) interface{} {
    return b
}

func TestInterfaceNil(t *testing.T) {
    var o reflection.Object
    i := returnInterface(o)
    if i == nil {
        t.Log("nil")
    } else {
        t.Fatal("not nil")
    }
}

func TestInterfaceNil2(t *testing.T) {
    var o reflection.Object
    var st *reflection.StructInfo
    st = nil
    o = st
    //i := returnInterface(o)
    if o == nil {
        t.Fatal("nil")
    } else {
        t.Log("not nil")
    }

    if reflection.IsNil(o) {
        t.Log("nil")
    } else {
        t.Fatal("not nil")
    }
}

func TestStructNil(t *testing.T) {
    o := (*TestTable)(nil)
    i := returnInterface(o)
    //var obj reflection.Object
    //i = obj
    if i == nil {
        t.Fatal("nil")
    } else {
        t.Log("not nil")
    }
    if reflection.IsNil(i) {
        t.Log("nil")
    } else {
        t.Fatal("not nil")
    }
}

func printTableInfo(table *reflection.StructInfo) {
    fmt.Printf("table name is %s\n", table.Name)
    //for k, v := range table.FieldMap {
    //    fmt.Printf("field : %s, value %s\n", k, v)
    //}

    for k, v := range table.FieldNameMap {
        fmt.Printf("origin : %s, map %s\n", k, v)
    }
}
