/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package handler

import (
    "github.com/xfali/gobatis/errors"
    "github.com/xfali/gobatis/reflection"
    "reflect"
)

type DefaultResultHandle struct {
    tableName string
}

func (r *DefaultResultHandle)Deserialize(columns []string, value []interface{}) (interface{}, error) {
    mi := reflection.FindModelInfo(r.tableName)
    if mi == nil {
        return nil, errors.MODEL_NOT_REGISTER
    }
    ti := mi.TableInfo
    v := reflect.New(reflect.TypeOf(mi.Model).Elem()).Elem()
    for i:=0; i<len(columns); i++{
        if fieldName, ok := ti.TypeMap[columns[i]]; ok {
            f := v.FieldByName(fieldName)
            ret := reflection.SetField(f, value)
            if !ret {
                continue
            }
        }
    }
    return v.Interface(), nil
}
