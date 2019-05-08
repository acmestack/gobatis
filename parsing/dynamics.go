/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package parsing

import (
    "fmt"
    "github.com/xfali/gobatis/logging"
    "github.com/xfali/gobatis/reflection"
    "reflect"
    "strings"
    "time"
)

type GetFunc func(key string) string

type DynamicElement interface {
    Format(func(key string) string) string
}

type DynamicData struct {
    OriginData     string
    DynamicElemMap map[string]DynamicElement
}

func (m *DynamicData) Replace(params ...interface{}) string {
    if len(m.DynamicElemMap) == 0 || len(params) == 0 {
        return m.OriginData
    }

    if len(params) == 1 {
        t := reflect.TypeOf(params[0])
        if t.Kind() == reflect.Ptr {
            t = t.Elem()
        }
        if reflection.IsSimpleType(params[0]) {
            return m.ReplaceWithParams(params...)
        }
        if t.Kind() == reflect.Struct {
            return m.ReplaceWithBean(params[0])
        }
    } else {
        objParams := map[string]interface{}{}
        for i, v := range params {
            if !reflection.IsSimpleType(v) {
                logging.Warn("Param error: expect simple type, but get other type")
                return m.OriginData
            }
            key := fmt.Sprintf("{%d}", i)
            objParams[key] = v
        }
        return m.ReplaceWithMap(objParams)
    }

    return m.OriginData
}

//需要外部确保param是一个struct
func (m *DynamicData) ReplaceWithBean(param interface{}) string {
    if len(m.DynamicElemMap) == 0 {
        return m.OriginData
    }

    ti, err := reflection.GetTableInfo(param)
    if err != nil {
        logging.Info("%s", err.Error())
        return m.OriginData
    }
    objParams := ti.MapValue()

    return m.ReplaceWithMap(objParams)
}

//需要外部确认params是一个简单类型（simple type）的切片slice
func (m *DynamicData) ReplaceWithParams(params ...interface{}) string {
    if len(m.DynamicElemMap) == 0 || len(params) == 0 {
        return m.OriginData
    }
    objParams := map[string]interface{}{}
    for i, v := range params {
        //if !reflection.IsSimpleType(v) {
        //    logging.Warn("Param error: expect simple type, but get other type")
        //    return m.OriginData
        //}
        key := fmt.Sprintf("{%d}", i)
        objParams[key] = v
    }
    return m.ReplaceWithMap(objParams)
}

//需要外部确保param是一个struct
func (m *DynamicData) ReplaceWithMap(objParams map[string]interface{}) string {
    if len(m.DynamicElemMap) == 0 || len(objParams) == 0 {
        logging.Info("map is empty")
        return m.OriginData
    }

    getFunc := func(s string) string {
        if o, ok := objParams[s]; ok {
            //zero time convert to empty string (for <if> </if> element)
            if ti, ok := o.(time.Time); ok {
                if ti.IsZero() {
                    return ""
                } else {
                    return ti.String()
                }
            }
            var str string
            reflection.SafeSetValue(reflect.ValueOf(&str), o)
            return str
        }
        return ""
    }

    ret := m.OriginData
    for k, v := range m.DynamicElemMap {
        ret = strings.Replace(ret, k, v.Format(getFunc), -1)
    }
    return ret
}
