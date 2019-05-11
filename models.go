/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package gobatis

import (
    "github.com/xfali/gobatis/errors"
    "github.com/xfali/gobatis/reflection"
    "reflect"
    "sync"
)

type ModelName string

type ModelInfo struct {
    ObjectInfo *reflection.ObjectInfo
    Model      interface{}
}

type ModelManager struct {
    modelMap map[string]*ModelInfo
    lock     sync.Mutex
}

var g_model_mgr = &ModelManager{modelMap: map[string]*ModelInfo{}}

func init() {
    registerBuildin()
}

func registerBuildin() {
    g_model_mgr.lock.Lock()
    defer g_model_mgr.lock.Unlock()

    g_model_mgr.modelMap[reflection.StringType.Name()] = &ModelInfo{ObjectInfo: nil, Model: &reflection.EMPTY_STRING}
    g_model_mgr.modelMap[reflection.BoolType.Name()] = &ModelInfo{ObjectInfo: nil, Model: &reflection.BOOL_DEFAULT}
    g_model_mgr.modelMap[reflection.ByteType.Name()] = &ModelInfo{ObjectInfo: nil, Model: &reflection.BYTE_DEFAULT}
    g_model_mgr.modelMap[reflection.Complex64Type.Name()] = &ModelInfo{ObjectInfo: nil, Model: &reflection.COMPLEX64_DEFAULT}
    g_model_mgr.modelMap[reflection.Complex128Type.Name()] = &ModelInfo{ObjectInfo: nil, Model: &reflection.COMPLEX128_DEFAULT}
    g_model_mgr.modelMap[reflection.Float32Type.Name()] = &ModelInfo{ObjectInfo: nil, Model: &reflection.FLOAT32_DEFAULT}
    g_model_mgr.modelMap[reflection.Float64Type.Name()] = &ModelInfo{ObjectInfo: nil, Model: &reflection.FLOAT64_DEFAULT}
    g_model_mgr.modelMap[reflection.Int64Type.Name()] = &ModelInfo{ObjectInfo: nil, Model: &reflection.INT64_DEFAULT}
    g_model_mgr.modelMap[reflection.Uint64Type.Name()] = &ModelInfo{ObjectInfo: nil, Model: &reflection.UINT64_DEFAULT}
    g_model_mgr.modelMap[reflection.Int32Type.Name()] = &ModelInfo{ObjectInfo: nil, Model: &reflection.INT32_DEFAULT}
    g_model_mgr.modelMap[reflection.Uint32Type.Name()] = &ModelInfo{ObjectInfo: nil, Model: &reflection.UINT32_DEFAULT}
    g_model_mgr.modelMap[reflection.Int16Type.Name()] = &ModelInfo{ObjectInfo: nil, Model: &reflection.INT16_DEFAULT}
    g_model_mgr.modelMap[reflection.Uint16Type.Name()] = &ModelInfo{ObjectInfo: nil, Model: &reflection.UINT16_DEFAULT}
    g_model_mgr.modelMap[reflection.Int8Type.Name()] = &ModelInfo{ObjectInfo: nil, Model: &reflection.INT8_DEFAULT}
    g_model_mgr.modelMap[reflection.Uint8Type.Name()] = &ModelInfo{ObjectInfo: nil, Model: &reflection.UINT8_DEFAULT}
    g_model_mgr.modelMap[reflection.IntType.Name()] = &ModelInfo{ObjectInfo: nil, Model: &reflection.INT_DEFAULT}
    g_model_mgr.modelMap[reflection.UintType.Name()] = &ModelInfo{ObjectInfo: nil, Model: &reflection.UINT_DEFAULT}
    g_model_mgr.modelMap[reflection.TimeType.Name()] = &ModelInfo{ObjectInfo: nil, Model: &reflection.TIME_DEFAULT}
}

// 注册模型，模型描述了column和field之间的关联关系；
// 用于获得数据库数据反序列化。未注册的模型将无法正确反序列化。
func RegisterModel(model interface{}) *errors.ErrCode {
    tableInfo, err := reflection.GetObjectInfo(model)
    if err != nil {
        return errors.PARSE_MODEL_TABLEINFO_FAILED
    }
    g_model_mgr.lock.Lock()
    defer g_model_mgr.lock.Unlock()
    g_model_mgr.modelMap[tableInfo.GetClassName()] = &ModelInfo{ObjectInfo: tableInfo, Model: model}
    return nil
}

func RegisterModelWithName(name string, model interface{}) *errors.ErrCode {
    tableInfo, err := reflection.GetObjectInfo(model)
    if err != nil {
        return errors.PARSE_MODEL_TABLEINFO_FAILED
    }
    g_model_mgr.lock.Lock()
    defer g_model_mgr.lock.Unlock()
    g_model_mgr.modelMap[name] = &ModelInfo{ObjectInfo: tableInfo, Model: model}
    return nil
}

func FindModelInfo(name string) *ModelInfo {
    g_model_mgr.lock.Lock()
    defer g_model_mgr.lock.Unlock()

    return g_model_mgr.modelMap[name]
}

func FindModelInfoOfBean(bean interface{}) *ModelInfo {
    name := reflection.GetBeanClassName(bean)
    return FindModelInfo(name)
}

func (mi *ModelInfo) Deserialize(columns []string, values []interface{}) (interface{}, error) {
    ti := mi.ObjectInfo
    v := reflect.New(reflect.TypeOf(mi.Model).Elem()).Elem()
    //struct
    if ti != nil {
        for i := 0; i < len(columns); i++ {
            if fieldName, ok := ti.FieldNameMap[columns[i]]; ok {
                f := v.FieldByName(fieldName)
                ret := reflection.SetValue(f, values[i])
                if !ret {
                    continue
                }
            }
        }
    } else {
        if len(values) > 1 {
            return nil, errors.DESERIALIZE_FAILED
        }
        if !reflection.SetValue(v, values[0]) {
            return nil, errors.DESERIALIZE_FAILED
        }
    }

    return v.Interface(), nil
}
