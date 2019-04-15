/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package config

import (
    "github.com/xfali/gobatis/errors"
    "github.com/xfali/gobatis/reflection"
    "reflect"
    "sync"
)

type ModelInfo struct {
    TableInfo *reflection.TableInfo
    Model     interface{}
}

type ModelManager struct {
    modelMap map[string]*ModelInfo
    lock     sync.Mutex
}

var g_model_mgr ModelManager = ModelManager{modelMap: map[string]*ModelInfo{}}

func RegisterModel(model interface{}) *errors.ErrCode {
    tableInfo, err := reflection.GetTableInfo(model)
    if err != nil {
        return errors.Parse_MODEL_TABLEINFO_FAILED
    }
    g_model_mgr.lock.Lock()
    defer g_model_mgr.lock.Unlock()
    g_model_mgr.modelMap[tableInfo.Name] = &ModelInfo{TableInfo: tableInfo, Model: model}
    return nil
}

func FindModelInfo(tableName string) *ModelInfo {
    g_model_mgr.lock.Lock()
    defer g_model_mgr.lock.Unlock()

    return g_model_mgr.modelMap[tableName]
}

func FindModelInfoOfBean(bean interface{}) *ModelInfo {
    name := reflection.GetTableName(bean)
    return FindModelInfo(name)
}

func (mi *ModelInfo) Deserialize(columns []string, values []interface{}) (interface{}, error) {
    ti := mi.TableInfo
    v := reflect.New(reflect.TypeOf(mi.Model).Elem()).Elem()
    for i := 0; i < len(columns); i++ {
        if fieldName, ok := ti.FieldNameMap[columns[i]]; ok {
            f := v.FieldByName(fieldName)
            ret := reflection.SetField(f, values[i])
            if !ret {
                continue
            }
        }
    }
    return v.Interface(), nil
}

