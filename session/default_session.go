/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package session

import (
    "github.com/xfali/GoBatis/executor"
    "github.com/xfali/GoBatis/handler"
    "github.com/xfali/GoBatis/logging"
    "github.com/xfali/GoBatis/statement"
    "github.com/xfali/GoBatis/transaction"
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

func (sess *DefaultSqlSession) Select(handler handler.ResultHandler, sql string, params ...string) error {
    sess.logLastSql(sql, params...)
    return sess.executor.Query(sess.getStatement(sql, handler), params)
}

func (sess *DefaultSqlSession) Insert(sql string, params ...string) int64 {
    sess.logLastSql(sql, params...)
    return sess.executor.Exec()
}

func (sess *DefaultSqlSession) Update(sql string, params ...string) int64 {
    sess.logLastSql(sql, params...)
}

func (sess *DefaultSqlSession) Delete(sql string, params ...string) int64 {
    sess.logLastSql(sql, params...)
}

func (sess *DefaultSqlSession) Commit() {
    sess.logLastSql("Commit", "")
}

func (sess *DefaultSqlSession) Rollback() {
    sess.logLastSql("Rollback", "")
}

func (sess *DefaultSqlSession) logLastSql(sql string, params ... string) {
    sess.Log(logging.INFO, "sql: %s, param %v", params)
}

func (sess *DefaultSqlSession) getStatement(sql string, handler handler.ResultHandler) statement.MappedStatement {
    return statement.NewSimpleStatment(sql, handler)
}
