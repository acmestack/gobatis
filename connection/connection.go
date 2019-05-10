/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package connection

import (
    "context"
    "github.com/xfali/gobatis/common"
    "github.com/xfali/gobatis/handler"
    "github.com/xfali/gobatis/statement"
)

type Connection interface {
    Prepare(sql string) (statement.Statement, error)
    Query(ctx context.Context, handler handler.ResultHandler, iterFunc common.IterFunc, sql string, params ...interface{}) error
    Exec(ctx context.Context, sql string, params ...interface{}) (common.Result, error)
}
