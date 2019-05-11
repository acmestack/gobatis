/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package reflection

import (
    "encoding/json"
    "fmt"
    "github.com/xfali/gobatis/common"
    "github.com/xfali/gobatis/errors"
    "github.com/xfali/gobatis/logging"
    "reflect"
    "strconv"
    "strings"
    "time"
)

var modelNameType reflect.Type

func SetModelNameType(mtype reflect.Type) {
    modelNameType = mtype
}

type FieldInfo struct {
    //字段名
    Name string
    //值
    Value reflect.Value
}

type ObjectInfo struct {
    //表名
    Name string
    //字段信息
    //Fields []FieldInfo
    FieldMap map[string]reflect.Value
    //表字段和实体字段映射关系
    FieldNameMap map[string]string
}

func newObjectInfo() *ObjectInfo {
    return &ObjectInfo{
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

//GetObjectInfo 解析结构体，使用：
//1、如果结构体中含有gobatis.ModelName类型的字段，则：
// a)、如果含有tag，则使用tag作为tablename；
// b)、如果不含有tag，则使用fieldName作为tablename。
//2、如果结构体中不含有gobatis.ModelName类型的字段，则使用结构体名称作为tablename
//3、如果结构体中含有xfield的tag，则：
// a）、如果tag为‘-’，则不进行columne与field的映射；
// b）、如果tag不为‘-’使用tag name作为column名称与field映射。
//4、如果结构体中不含有xfield的tag，则使用field name作为column名称与field映射
//5、如果字段的tag为‘-’，则不进行columne与field的映射；
func GetObjectInfo(model interface{}) (*ObjectInfo, error) {
    tableInfo := newObjectInfo()

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
            } else {
                tableInfo.Name = rtf.Name
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
        tagName := rtf.Tag.Get(common.FIELD_NAME)
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

func (ti *ObjectInfo) MapValue() map[string]interface{} {
    paramMap := map[string]interface{}{}
    ti.FillMapValue(&paramMap)
    return paramMap
}

func (ti *ObjectInfo) FillMapValue(paramMap *map[string]interface{}) {
    for k, v := range ti.FieldMap {
        if !v.CanInterface() {
            v = reflect.Indirect(v)
        }
        (*paramMap)[k] = v.Interface()
    }
    //(*paramMap)["tablename"] = ti.Name
}

func ReflectValue(bean interface{}) reflect.Value {
    return reflect.Indirect(reflect.ValueOf(bean))
}

func IsSimpleObject(bean interface{}) bool {
    rt := reflect.TypeOf(bean)
    if rt.Kind() == reflect.Ptr {
        rt = rt.Elem()
    }
    return IsSimpleType(rt)
}

//是否是数据库使用的简单类型，注意不能是PTR
func IsSimpleType(t reflect.Type) bool {
    switch t.Kind() {
    case IntKind, Int8Kind, Int16Kind, Int32Kind, Int64Kind, UintKind, Uint8Kind, Uint16Kind, Uint32Kind, Uint64Kind,
        Float32Kind, Float64Kind, Complex64Kind, Complex128Kind, StringKind, BoolKind, ByteKind, BytesKind /*, TimeKind*/ :
        return true
    }
    if t.ConvertibleTo(TimeType) {
        return true
    }
    return false
}

func checkBeanValue(beanValue reflect.Value) bool {
    if beanValue.Kind() != reflect.Ptr {
        return false
    } else if beanValue.Elem().Kind() == reflect.Ptr {
        return false
    }
    return true
}

func SafeSetValue(f reflect.Value, v interface{}) bool {
    if !checkBeanValue(f) {
        logging.Info("value cannot be set")
        return false
    }
    f = f.Elem()
    return SetValue(f, v)
}

func SetValue(f reflect.Value, v interface{}) bool {
    if v == nil {
        return false
    }

    hasAssigned := false
    rawValue := reflect.Indirect(reflect.ValueOf(v))
    rawValueType := reflect.TypeOf(rawValue.Interface())
    vv := reflect.ValueOf(rawValue.Interface())

    ft := f.Type()
    switch ft.Kind() {
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
        case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
            hasAssigned = true
            f.SetString(strconv.FormatUint(vv.Uint(), 10))
            break
        case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
            hasAssigned = true
            f.SetString(strconv.FormatInt(vv.Int(), 10))
            break
        case reflect.Float64:
            hasAssigned = true
            f.SetString(strconv.FormatFloat(vv.Float(), 'g', -1, 64))
            break
        case reflect.Float32:
            hasAssigned = true
            f.SetString(strconv.FormatFloat(vv.Float(), 'g', -1, 32))
            break
        case reflect.Bool:
            hasAssigned = true
            f.SetString(strconv.FormatBool(vv.Bool()))
            break
        //case reflect.Struct:
        //    if ti, ok := v.(time.Time); ok {
        //        hasAssigned = true
        //        if ti.IsZero() {
        //            f.SetString("")
        //        } else {
        //            f.SetString(ti.String())
        //        }
        //    } else {
        //        hasAssigned = true
        //        f.SetString(fmt.Sprintf("%v", v))
        //    }
        default:
            hasAssigned = true
            f.SetString(fmt.Sprintf("%v", v))
        }
        break
    case reflect.Complex64, reflect.Complex128:
        switch rawValueType.Kind() {
        case reflect.Complex64, reflect.Complex128:
            hasAssigned = true
            f.SetComplex(vv.Complex())
            break
        case reflect.Slice:
            if rawValueType.ConvertibleTo(BytesType) {
                d := vv.Bytes()
                if len(d) > 0 {
                    if f.CanAddr() {
                        err := json.Unmarshal(d, f.Addr().Interface())
                        if err != nil {
                            return false
                        }
                    } else {
                        x := reflect.New(ft)
                        err := json.Unmarshal(d, x.Interface())
                        if err != nil {
                            return false
                        }
                        hasAssigned = true
                        f.Set(x.Elem())
                        break
                    }
                }
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

func ToSlice(arr interface{}) []interface{} {
    v := reflect.ValueOf(arr)
    if v.Kind() != reflect.Slice {
        panic("toslice arr not slice")
    }
    l := v.Len()
    ret := make([]interface{}, l)
    for i := 0; i < l; i++ {
        ret[i] = v.Index(i).Interface()
    }
    return ret
}