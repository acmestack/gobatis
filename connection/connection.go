/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package connection

import "github.com/xfali/gobatis/handler"

type IterFunc func(idx int64, bean interface{}) bool

type Statement interface {
    Query(handler handler.ResultHandler, iterFunc IterFunc, params ...interface{}) error
    Exec(params ...interface{}) (int64, error)
    Close()
}

type Connection interface {
    Prepare(sql string) (Statement, error)
}
