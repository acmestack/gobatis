/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package session

import (
    "github.com/xfali/gobatis"
    "github.com/xfali/gobatis/handler"
)

type Session interface {
    Close(rollback bool)

    SelectOne(handler handler.ResultHandler, sql string, params ...interface{}) (interface{}, error)

    Select(handler handler.ResultHandler, sql string, params ...interface{}) ([]interface{}, error)

    Query(handler handler.ResultHandler, iterFunc gobatis.IterFunc, sql string, params ...interface{}) error

    Insert(sql string, params ...interface{}) int64

    Update(sql string, params ...interface{}) int64

    Delete(sql string, params ...interface{}) int64

    Begin()

    Commit()

    Rollback()
}
