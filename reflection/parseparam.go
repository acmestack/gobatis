/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package reflection

import (
    "reflect"
    "strconv"
)

type paramParser struct {
    ret   map[string]interface{}
    index int
}

func ParseParams(params ...interface{}) map[string]interface{} {
    parser := paramParser{
        ret:   map[string]interface{}{},
        index: 0,
    }
    parser.innerParse(params...)
    return parser.ret
}

func (parser *paramParser)innerParse(params ...interface{}) {
    for i := range params {
        parser.parseOne(params[i])
    }
}

func (parser *paramParser)parseOne(v interface{}) {
    rt := reflect.TypeOf(v)
    rv := reflect.ValueOf(v)

    if rt.Kind() == reflect.Ptr {
        rt = rt.Elem()
        rv = rv.Elem()
    }

    if IsSimpleType(rt) {
        parser.ret[strconv.Itoa(parser.index)] = v
        parser.index++
    } else if rt.Kind() == reflect.Struct {
        oi, _ := GetStructInfo(v)
        structMap := oi.MapValue()
        for key, value := range structMap {
            parser.ret[structKey(oi, key)] = value
        }
    } else if rt.Kind() == reflect.Slice {
        //l := rv.Len()
        //for i := 0; i < l; i++ {
        //    elemV := rv.Index(i)
        //    if !elemV.CanInterface() {
        //        elemV = reflect.Indirect(elemV)
        //    }
        //    parser.parseOne(elemV.Interface())
        //}
        parser.ret[strconv.Itoa(parser.index)] = v
        parser.index++
    } else if rt.Kind() == reflect.Map {
        keys := rv.MapKeys()
        for _, key := range keys {
            if key.Kind() == reflect.String {
                value := rv.MapIndex(key)
                value = value.Elem()
                if IsSimpleType(value.Type()) {
                    if !value.CanInterface() {
                        value = reflect.Indirect(value)
                    }
                    parser.ret[key.String()] = value
                }
            }
        }
    }
}

func structKey(oi *StructInfo, field string) string {
    return oi.Name + "." + field
}
