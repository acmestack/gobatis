/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package test

import (
    "github.com/xfali/gobatis/reflection"
    "reflect"
    "testing"
    "time"
)

func TestReflectObjectStruct(t *testing.T) {
    v := TestTable{}
    info, err := reflection.GetObjectInfo(&v)
    if err != nil {
        t.Fatal()
    }
    t.Logf("classname :%v", info.GetClassName())
    t.Log(v)
    newOne := TestTable{
        Username: "1",
    }
    info.SetValue(reflect.ValueOf(newOne))

    t.Logf("after set :%v\n", v)

    info.SetField("username", reflect.ValueOf("123"))
    t.Logf("after setField :%v\n", v)
}

func TestReflectObjectSimpleTime(t *testing.T) {
    v := time.Time{}
    info, err := reflection.GetObjectInfo(&v)
    if err != nil {
        t.Fatal()
    }
    t.Logf("classname :%v", info.GetClassName())
    t.Log(v)
    newOne := TestTable{
        Username: "1",
    }
    info.SetValue(reflect.ValueOf(newOne))

    t.Logf("after set error type :%v\n", v)

    info.SetValue(reflect.ValueOf(time.Now()))

    t.Logf("after set now type :%v\n", v)

    info.SetField("username", reflect.ValueOf("123"))
    t.Logf("after setField :%v\n", v)
}

func TestReflectObjectSimpleFloat(t *testing.T) {
    v := 0.0
    info, err := reflection.GetObjectInfo(&v)
    if err != nil {
        t.Fatal()
    }
    t.Logf("classname :%v", info.GetClassName())
    t.Log(v)
    info.SetValue(reflect.ValueOf(1))

    t.Logf("after set int type :%v\n", v)

    info.SetValue(reflect.ValueOf(1.5))

    t.Logf("after set float type :%v\n", v)
}

func TestReflectObjectMap(t *testing.T) {
    v := map[string]interface{}{}
    info, err := reflection.GetObjectInfo(&v)
    if err != nil {
        t.Fatal()
    }
    t.Logf("classname :%v", info.GetClassName())
    t.Log(v)
    info.SetValue(reflect.ValueOf(1))

    t.Logf("after set int type :%v\n", v)

    info.SetValue(reflect.ValueOf(map[string]int{"1": 1, "2": 2}))

    t.Logf("after set map[string]int type :%v\n", v)

    info.SetValue(reflect.ValueOf(map[string]interface{}{"1": 1, "2": 2}))

    t.Logf("after set map[string]interface{} type :%v\n", v)

    info.SetField("username", reflect.ValueOf("123"))
    t.Logf("after setField :%v\n", v)
}

func TestReflectObjectSlice2(t *testing.T) {
    v := []int{}
    info, err := reflection.GetObjectInfo(&v)
    if err != nil {
        t.Fatal()
    }
    t.Logf("classname :%v", info.GetClassName())
    t.Log(v)
    info.SetValue(reflect.ValueOf(1))

    t.Logf("after set int type :%v\n", v)

    info.SetValue(reflect.ValueOf([]float32{1.0,2,3}))

    t.Logf("after set []float32{1.0,2,3} :%v\n", v)

    info.SetValue(reflect.ValueOf([]int{1,2,3}))

    t.Logf("after set []int{1,2,3} :%v\n", v)

    info.SetField("username", reflect.ValueOf(123))
    t.Logf("after setField :%v\n", v)

    info.AddValue(reflect.ValueOf(123))
    t.Logf("after AddValue :%v\n", v)
}

func TestReflectObjectSlice(t *testing.T) {
    v := []TestTable{}
    info, err := reflection.GetObjectInfo(&v)
    if err != nil {
        t.Fatal()
    }
    t.Logf("classname :%v", info.GetClassName())
    t.Log(v)
    info.SetValue(reflect.ValueOf(1))

    t.Logf("after set int type :%v\n", v)

    info.SetValue(reflect.ValueOf([]float32{1.0,2,3}))

    t.Logf("after set []float32{1.0,2,3} :%v\n", v)

    info.SetValue(reflect.ValueOf([]TestTable{{Username:"1", Password:"1"}}))

    t.Logf(`after set []TestTable{{Username:"1", Password:"1"}} :%v\n`, v)

    info.SetField("username", reflect.ValueOf(123))
    t.Logf("after setField :%v\n", v)

    info.AddValue(reflect.ValueOf(1))
    t.Logf("after AddValue 1 :%v\n", v)

    info.AddValue(reflect.ValueOf(TestTable{Username: "2", Password:"2"}))
    t.Logf(`after AddValue TestTable{Username: "2", Password:"2"} :%v\n`, v)

    ev := info.NewElemValue()
    einfo, _ := reflection.GetReflectObjectInfo(ev.Type(), ev)
    einfo.SetField("username", "x")
    info.AddValue(ev)

    t.Logf(`after AddValue new elem {Username: "x"} :%v\n`, v)
}
