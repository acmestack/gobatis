/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package reflection

import (
    "github.com/xfali/gobatis"
    "reflect"
)

var typeTableName gobatis.TableName
var tableNameType = reflect.TypeOf(typeTableName)

type FieldInfo struct {
    //字段名
    Name string
    //值
    Value reflect.Value
}

type TableInfo struct {
    //表名
    Name string
    //字段信息
    //Fields []FieldInfo
    FieldMap map[string]reflect.Value
    //表字段和实体字段映射关系
    TypeMap map[string]string
}

func newTableInfo() *TableInfo {
    return &TableInfo{
        FieldMap: map[string]reflect.Value{},
        TypeMap:  map[string]string{},
    }
}

func GetTableInfo(model interface{}) (*TableInfo, error) {
    tableInfo := newTableInfo()

    rt := reflect.TypeOf(model)
    rv := reflect.ValueOf(model)

    if rt.Kind() == reflect.Ptr {
        rt = rt.Elem()
        rv = rv.Elem()
    }
    //Default name is struct name
    tableInfo.Name = rt.Name()

    //字段解析
    for i, j := 0, rt.NumField(); i < j; i++ {
        rtf := rt.Field(i)
        rvf := rv.Field(i)
        if rtf.Type == tableNameType {
            if rtf.Tag != "" {
                tableInfo.Name = string(rtf.Tag)
            }
            continue
        }

        //没有tag,表字段名与实体字段名一致
        if rtf.Tag == "" {
            tableInfo.TypeMap[rtf.Name] = rtf.Name
            //f := FieldInfo{Name: rtf.Name, Value: rvf}
            //tableInfo.Fields = append(tableInfo.Fields, f)
            tableInfo.FieldMap[rtf.Name] = rvf
            continue
        }

        if rtf.Tag == "-" {
            continue
        }

        fieldName := rtf.Name
        tagName := rtf.Tag.Get(gobatis.FIELD_NAME)
        if tagName == "-" {
            continue
        } else if tagName != "" {
            fieldName = tagName
        }
        tableInfo.TypeMap[rtf.Name] = fieldName
        //f := FieldInfo{Name: fieldName, Value: rvf}
        //tableInfo.Fields = append(tableInfo.Fields, f)
        tableInfo.FieldMap[fieldName] = rvf
        continue
    }
    return tableInfo, nil
}

func (ti *TableInfo) MapValue() map[string]interface{} {
    params := map[string]interface{}{}
    for k, v := range ti.FieldMap {
        if !v.CanInterface() {
            v = reflect.Indirect(v)
        }
        params[k] = v.Interface()
    }
    params["tablename"] = ti.Name
    return params
}

func ReflectValue(bean interface{}) reflect.Value {
    return reflect.Indirect(reflect.ValueOf(bean))
}

func SetField(f reflect.Value, v interface{}) bool {
    hasAssigned := false
    rawValue := reflect.Indirect(reflect.ValueOf(v))
    rawValueType := reflect.TypeOf(rawValue.Interface())
    vv := reflect.ValueOf(rawValue.Interface())

    switch f.Type().Kind() {
    case reflect.Bool:
        if rawValueType.Kind() == reflect.Bool {
            hasAssigned = true
            f.SetBool(vv.Bool())
        }
        break
    case reflect.String:
        if rawValueType.Kind() == reflect.String {
            hasAssigned = true
            f.SetString(vv.String())
        }
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        switch rawValueType.Kind() {
        case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
            hasAssigned = true
            f.SetInt(vv.Int())
        }
    case reflect.Float32, reflect.Float64:
        switch rawValueType.Kind() {
        case reflect.Float32, reflect.Float64:
            hasAssigned = true
            f.SetFloat(vv.Float())
        }
    case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
        switch rawValueType.Kind() {
        case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
            hasAssigned = true
            f.SetUint(vv.Uint())
        case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
            hasAssigned = true
            f.SetUint(uint64(vv.Int()))
        }
    }

    return hasAssigned
}
