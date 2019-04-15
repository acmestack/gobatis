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
    "github.com/xfali/gobatis/errors"
    "reflect"
    "strconv"
    "strings"
    "time"
)

var typeModelName gobatis.ModelName
var modelNameType = reflect.TypeOf(typeModelName)

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
    FieldNameMap map[string]string
}

func newTableInfo() *TableInfo {
    return &TableInfo{
        FieldMap:     map[string]reflect.Value{},
        FieldNameMap: map[string]string{},
    }
}

func GetBeanName(model interface{}) string {
    rt := reflect.TypeOf(model)
    rv := reflect.ValueOf(model)

    if rt.Kind() == reflect.Ptr {
        rt = rt.Elem()
        rv = rv.Elem()
    }

    if rt.Kind() == reflect.Slice {
        rt = rt.Elem()
    }

    //Default name is struct name
    name := rt.Name()

    if rt.Kind() != reflect.Struct {
        return name
    }

    //字段解析
    for i, j := 0, rt.NumField(); i < j; i++ {
        rtf := rt.Field(i)
        if rtf.Type == modelNameType {
            if rtf.Tag != "" {
                name = string(rtf.Tag)
            }
        }
    }
    return name
}

func GetTableInfo(model interface{}) (*TableInfo, error) {
    tableInfo := newTableInfo()

    rt := reflect.TypeOf(model)
    rv := reflect.ValueOf(model)

    if rt.Kind() == reflect.Ptr {
        rt = rt.Elem()
        rv = rv.Elem()
    }

    if rt.Kind() != reflect.Struct {
       return nil, errors.PARSE_TABLEINFO_NOT_STRUCT
    }
    //Default name is struct name
    tableInfo.Name = rt.Name()

    //字段解析
    for i, j := 0, rt.NumField(); i < j; i++ {
        rtf := rt.Field(i)
        rvf := rv.Field(i)
        if rtf.Type == modelNameType {
            if rtf.Tag != "" {
                tableInfo.Name = string(rtf.Tag)
            }
            continue
        }

        //没有tag,表字段名与实体字段名一致
        if rtf.Tag == "" {
            tableInfo.FieldNameMap[rtf.Name] = rtf.Name
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
        tableInfo.FieldNameMap[fieldName] = rtf.Name
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

func SetValue(f reflect.Value, v interface{}) bool {
    hasAssigned := false
    rawValue := reflect.Indirect(reflect.ValueOf(v))
    rawValueType := reflect.TypeOf(rawValue.Interface())
    vv := reflect.ValueOf(rawValue.Interface())

    switch f.Type().Kind() {
    case reflect.Bool:
        switch rawValueType.Kind() {
        case reflect.Bool:
            hasAssigned = true
            f.SetBool(vv.Bool())
            break
        case reflect.Slice:
            if d, ok := vv.Interface().([]uint8); ok {
                hasAssigned = true
                f.SetBool(d[0] != 0)
            }
            break
        }
        break
    case reflect.String:
        switch rawValueType.Kind() {
        case reflect.String:
            hasAssigned = true
            f.SetString(vv.String())
            break
        case reflect.Slice:
            if d, ok := vv.Interface().([]uint8); ok {
                hasAssigned = true
                f.SetString(string(d))
            }
            break
        }
        break
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        switch rawValueType.Kind() {
        case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
            hasAssigned = true
            f.SetInt(vv.Int())
            break
        case reflect.Slice:
            if d, ok := vv.Interface().([]uint8); ok {
                intV, err := strconv.ParseInt(string(d), 10, 64)
                if err == nil {
                    hasAssigned = true
                    f.SetInt(intV)
                }
            }
            break
        }
        break
    case reflect.Float32, reflect.Float64:
        switch rawValueType.Kind() {
        case reflect.Float32, reflect.Float64:
            hasAssigned = true
            f.SetFloat(vv.Float())
            break
        case reflect.Slice:
            if d, ok := vv.Interface().([]uint8); ok {
                floatV, err := strconv.ParseFloat(string(d), 64)
                if err == nil {
                    hasAssigned = true
                    f.SetFloat(floatV)
                }
            }
            break
        }
        break
    case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
        switch rawValueType.Kind() {
        case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
            hasAssigned = true
            f.SetUint(vv.Uint())
            break
        case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
            hasAssigned = true
            f.SetUint(uint64(vv.Int()))
            break
        case reflect.Slice:
            if d, ok := vv.Interface().([]uint8); ok {
                uintV, err := strconv.ParseUint(string(d), 10, 64)
                if err == nil {
                    hasAssigned = true
                    f.SetUint(uintV)
                }
            }
            break
        }
        break
    case reflect.Struct:
        fieldType := f.Type()
        if fieldType.ConvertibleTo(TimeType) {
            if rawValueType == TimeType {
                hasAssigned = true
                t := vv.Convert(TimeType).Interface().(time.Time)
                f.Set(reflect.ValueOf(t).Convert(fieldType))
            } else if rawValueType == IntType || rawValueType == Int64Type ||
                rawValueType == Int32Type {
                hasAssigned = true

                t := time.Unix(vv.Int(), 0)
                f.Set(reflect.ValueOf(t).Convert(fieldType))
            } else {
                if d, ok := vv.Interface().([]byte); ok {
                    t, err := convert2Time(d, time.Local)
                    if err == nil {
                        hasAssigned = true
                        f.Set(reflect.ValueOf(t).Convert(fieldType))
                    }
                }
            }
        } else {
            f.Set(reflect.ValueOf(v))
        }
    }

    return hasAssigned
}

const (
    zeroTime0 = "0000-00-00 00:00:00"
    zeroTime1 = "0001-01-01 00:00:00"
)

func convert2Time(data []byte, location *time.Location) (time.Time, error) {
    timeStr := strings.TrimSpace(string(data))
    var timeRet time.Time
    var err error
    if timeStr == zeroTime0 || timeStr == zeroTime1 {
    } else if !strings.ContainsAny(timeStr, "- :") {
        // time stamp
        sd, err := strconv.ParseInt(timeStr, 10, 64)
        if err == nil {
            timeRet = time.Unix(sd, 0)
        }
    } else if len(timeStr) > 19 && strings.Contains(timeStr, "-") {
        timeRet, err = time.ParseInLocation(time.RFC3339Nano, timeStr, location)
        if err != nil {
            timeRet, err = time.ParseInLocation("2006-01-02 15:04:05.999999999", timeStr, location)
        }
        if err != nil {
            timeRet, err = time.ParseInLocation("2006-01-02 15:04:05.9999999 Z07:00", timeStr, location)
        }
    } else if len(timeStr) == 19 && strings.Contains(timeStr, "-") {
        timeRet, err = time.ParseInLocation("2006-01-02 15:04:05", timeStr, location)
    } else if len(timeStr) == 10 && timeStr[4] == '-' && timeStr[7] == '-' {
        timeRet, err = time.ParseInLocation("2006-01-02", timeStr, location)
    }
    return timeRet, nil
}
