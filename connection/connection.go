/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package connection

import (
    "github.com/xfali/gobatis"
    "github.com/xfali/gobatis/handler"
    "github.com/xfali/gobatis/statement"
)

type Connection interface {
    Prepare(sql string) (statement.Statement, error)
    Query(handler handler.ResultHandler, iterFunc gobatis.IterFunc, sql string, params ...interface{}) error
    Exec(sql string, params ...interface{}) (int64, error)
}
