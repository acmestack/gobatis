/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package session

import (
    "context"
    "github.com/xfali/gobatis"
    "github.com/xfali/gobatis/handler"
)

type SqlSession interface {
    Close(rollback bool)

    SelectOne(ctx context.Context, handler handler.ResultHandler, sql string, params ...interface{}) (interface{}, error)

    Select(ctx context.Context, handler handler.ResultHandler, sql string, params ...interface{}) ([]interface{}, error)

    Query(ctx context.Context, handler handler.ResultHandler, iterFunc gobatis.IterFunc, sql string, params ...interface{}) error

    Insert(ctx context.Context, sql string, params ...interface{}) (int64, int64, error)

    Update(ctx context.Context, sql string, params ...interface{}) (int64, error)

    Delete(ctx context.Context, sql string, params ...interface{}) (int64, error)

    Begin()

    Commit()

    Rollback()
}
