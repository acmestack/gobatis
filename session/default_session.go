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
    "fmt"
    "github.com/xfali/gobatis"
    "github.com/xfali/gobatis/common"
    "github.com/xfali/gobatis/executor"
    "github.com/xfali/gobatis/handler"
    "github.com/xfali/gobatis/logging"
    "github.com/xfali/gobatis/transaction"
)

type DefaultSqlSession struct {
    Log        logging.LogFunc
    tx         transaction.Transaction
    executor   executor.Executor
    autoCommit bool
}

func NewDefaultSqlSession(log logging.LogFunc, tx transaction.Transaction, e executor.Executor, autoCommit bool) *DefaultSqlSession {
    ret := &DefaultSqlSession{
        Log:        log,
        tx:         tx,
        executor:   e,
        autoCommit: autoCommit,
    }
    return ret
}

func (sess *DefaultSqlSession) Close(rollback bool) {
    sess.executor.Close(rollback)
}

func (sess *DefaultSqlSession) SelectOne(ctx context.Context, handler handler.ResultHandler, sql string, params ...interface{}) (interface{}, error) {
    sess.logLastSql(sql, params...)
    var ret interface{} = nil
    iterFunc := func(idx int64, bean interface{}) bool {
        ret = bean
        return false
    }
    err := sess.executor.Query(ctx, handler, iterFunc, sql, params...)
    if err != nil {
        return nil, err
    }
    return ret, nil
}

func (sess *DefaultSqlSession) Select(ctx context.Context, handler handler.ResultHandler, sql string, params ...interface{}) ([]interface{}, error) {
    sess.logLastSql(sql, params...)
    var ret []interface{}
    iterFunc := func(idx int64, bean interface{}) bool {
        ret = append(ret, bean)
        return true
    }
    err := sess.executor.Query(ctx, handler, iterFunc, sql, params...)
    if err != nil {
        return nil, err
    }
    return ret, nil
}

func (sess *DefaultSqlSession) Query(ctx context.Context, handler handler.ResultHandler, iterFunc gobatis.IterFunc, sql string, params ...interface{}) error {
    sess.logLastSql(sql, params...)
    return sess.executor.Query(ctx, handler, iterFunc, sql, params...)
}

func (sess *DefaultSqlSession) Insert(ctx context.Context, sql string, params ...interface{}) (int64, int64, error) {
    sess.logLastSql(sql, params...)
    ret, err := sess.exec(ctx, sql, params...)
    if err != nil {
        return 0, -1, err
    }

    count, err := ret.RowsAffected()
    if err != nil {
        return 0, -1, err
    }

    id, err := ret.LastInsertId()
    if err != nil {
        return 0, -1, err
    }
    return count, id, nil
}

func (sess *DefaultSqlSession) Update(ctx context.Context, sql string, params ...interface{}) (int64, error) {
    sess.logLastSql(sql, params...)
    ret, err := sess.exec(ctx, sql, params...)
    if err != nil {
        return 0, err
    }

    count, err := ret.RowsAffected()
    if err != nil {
        return 0, err
    }
    return count, nil
}

func (sess *DefaultSqlSession) Delete(ctx context.Context, sql string, params ...interface{}) (int64, error) {
    sess.logLastSql(sql, params...)
    ret, err := sess.exec(ctx, sql, params...)
    if err != nil {
        return 0, err
    }

    count, err := ret.RowsAffected()
    if err != nil {
        return 0, err
    }
    return count, nil
}

func (sess *DefaultSqlSession) Begin() {
    sess.logLastSql("Begin", "")
    sess.tx.Begin()
}

func (sess *DefaultSqlSession) Commit() {
    sess.logLastSql("Commit", "")
    sess.tx.Commit()
}

func (sess *DefaultSqlSession) Rollback() {
    sess.logLastSql("Rollback", "")
    sess.tx.Rollback()
}

func (sess *DefaultSqlSession) logLastSql(sql string, params ...interface{}) {
    sess.Log(logging.INFO, "sql: [%s], param: %s", sql, fmt.Sprint(params))
}

func (sess *DefaultSqlSession) getSqlContext(ctx context.Context, sql string, handler handler.ResultHandler, iterFunc gobatis.IterFunc) *common.SqlContext {
    return &common.SqlContext{Ctx: ctx, Sql: sql, ResultHandler: handler, IterFunc: iterFunc}
}

func (sess *DefaultSqlSession) exec(ctx context.Context, sql string, params ...interface{}) (common.Result, error) {
    return sess.executor.Exec(ctx, sql, params...)
}
