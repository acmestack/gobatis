/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package reflection

import (
    "github.com/xfali/GoBatis"
    "reflect"
)

var typeTableName GoBatis.TableName
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
    Fields []FieldInfo
    //表字段和实体字段映射关系
    TypeMap map[string]string
}

func newTableInfo() *TableInfo {
    return &TableInfo{
        TypeMap: map[string]string{},
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
            f := FieldInfo{Name: rtf.Name, Value: rvf}
            tableInfo.TypeMap[rtf.Name] = rtf.Name
            tableInfo.Fields = append(tableInfo.Fields, f)
            continue
        }

        if rtf.Tag == "-" {
            continue
        }

        fieldName := rtf.Name
        tagName := rtf.Tag.Get(GoBatis.FIELD_NAME)
        if tagName == "-" {
            continue
        } else if tagName != "" {
            fieldName = tagName
        }
        f := FieldInfo{Name: fieldName, Value: rvf}
        tableInfo.TypeMap[rtf.Name] = fieldName
        tableInfo.Fields = append(tableInfo.Fields, f)
        continue
    }
    return tableInfo, nil
}

func ReflectValue(bean interface{}) reflect.Value {
    return reflect.Indirect(reflect.ValueOf(bean))
}
