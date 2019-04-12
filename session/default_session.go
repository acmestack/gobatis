/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package session

import (
    "fmt"
    "github.com/xfali/gobatis/connection"
    "github.com/xfali/gobatis/executor"
    "github.com/xfali/gobatis/handler"
    "github.com/xfali/gobatis/logging"
    "github.com/xfali/gobatis/statement"
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

func (sess *DefaultSqlSession) SelectOne(handler handler.ResultHandler, sql string, params ...interface{}) (interface{}, error) {
    sess.logLastSql(sql, params...)
    var ret interface{}
    iterFunc := func(idx int64, bean interface{}) bool {
        ret = bean
        return false
    }
    err := sess.executor.Query(sess.getStatement(sql, handler, iterFunc), params...)
    if err != nil {
        return nil, err
    }
    return ret, nil
}

func (sess *DefaultSqlSession) Select(handler handler.ResultHandler, sql string, params ...interface{}) ([]interface{}, error) {
    sess.logLastSql(sql, params...)
    var ret []interface{}
    iterFunc := func(idx int64, bean interface{}) bool {
        ret = append(ret, bean)
        return true
    }
    err := sess.executor.Query(sess.getStatement(sql, handler, iterFunc), params...)
    if err != nil {
        return nil, err
    }
    return ret, nil
}

func (sess *DefaultSqlSession) Insert(sql string, params ...interface{}) int64 {
    sess.logLastSql(sql, params...)
    return sess.exec(sql, params...)
}

func (sess *DefaultSqlSession) Update(sql string, params ...interface{}) int64 {
    sess.logLastSql(sql, params...)
    return sess.exec(sql, params...)
}

func (sess *DefaultSqlSession) Delete(sql string, params ...interface{}) int64 {
    sess.logLastSql(sql, params...)
    return sess.exec(sql, params...)
}

func (sess *DefaultSqlSession) Begin() {
    sess.logLastSql("Commit", "")
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

func (sess *DefaultSqlSession) getStatement(sql string, handler handler.ResultHandler, iterFunc connection.IterFunc) *statement.MappedStatement {
    return &statement.MappedStatement{Sql: sql, ResultHandler: handler, IterFunc: iterFunc}
}

func (sess *DefaultSqlSession) exec(sql string, params ...interface{}) int64 {
    i, e := sess.executor.Exec(sess.getStatement(sql, nil, nil), params)
    if e != nil {
        return 0
    } else {
        return i
    }
}
