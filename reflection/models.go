/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package reflection

import "github.com/xfali/gobatis/errors"

type ModelInfo struct {
    TableInfo *TableInfo
    Model     interface{}
}

var g_model_map map[string]*ModelInfo

func init() {
    g_model_map = map[string]*ModelInfo{}
}

func RegisterModel(model interface{}) *errors.ErrCode {
    tableInfo, err := GetTableInfo(model)
    if err != nil {
        return errors.Parse_MODEL_TABLEINFO_FAILED
    }
    g_model_map[tableInfo.Name] = &ModelInfo{TableInfo: tableInfo, Model: model}
    return nil
}

func FindModelInfo(tableName string) *ModelInfo {
    return g_model_map[tableName]
}
