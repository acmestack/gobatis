/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package session

import (
    "github.com/xfali/gobatis/handler"
)

//参数
//  idx:迭代数
//  bean:序列化后的值
//返回值:
//  打断迭代返回true
type IterFunc func(idx int64, bean interface{}) bool

type Session interface {
    Close(rollback bool)

    SelectOne(handler handler.ResultHandler, sql string, params ...interface{}) (interface{}, error)

    Select(handler handler.ResultHandler, sql string, params ...interface{}) ([]interface{}, error)

    Query(handler handler.ResultHandler, iterFunc IterFunc, sql string, params ...interface{}) error

    Insert(sql string, params ...interface{}) int64

    Update(sql string, params ...interface{}) int64

    Delete(sql string, params ...interface{}) int64

    Begin()

    Commit()

    Rollback()
}
