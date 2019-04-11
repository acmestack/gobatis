/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package executor

import (
    "github.com/xfali/gobatis/errors"
    "github.com/xfali/gobatis/statement"
    "github.com/xfali/gobatis/transaction"
)

type SimpleExecutor struct {
    transaction transaction.Transaction
    closed      bool
}

func NewSimpleExecutor(transaction transaction.Transaction) *SimpleExecutor {
    return &SimpleExecutor{transaction: transaction}
}

func (exec *SimpleExecutor) Close(rollback bool) {
    defer func() {
        if exec.transaction != nil {
            exec.transaction.Close()
        }
        exec.transaction = nil
        exec.closed = true
    }()

    if rollback {
        exec.Rollback(true)
    }
}

func (exec *SimpleExecutor) Query(statement *statement.MappedStatement, params ...interface{}) error {
    if exec.closed {
        return  errors.EXECUTOR_QUERY_ERROR
    }

    conn := exec.transaction.GetConnection()
    if conn == nil {
        return errors.EXECUTOR_GET_CONNECTION_ERROR
    }

    stmt, err := conn.Prepare(statement.Sql)
    if err != nil {
        return err
    }
    return stmt.Query(statement.ResultHandler, statement.IterFunc, params)
}

func (exec *SimpleExecutor) Exec(statement *statement.MappedStatement, params ...interface{}) (int64, error) {
    if exec.closed {
        return 0, errors.EXECUTOR_QUERY_ERROR
    }

    conn := exec.transaction.GetConnection()
    if conn == nil {
        return 0, errors.EXECUTOR_GET_CONNECTION_ERROR
    }

    stmt, err := conn.Prepare(statement.Sql)
    if err != nil {
        return 0, err
    }
    return stmt.Exec(params)
}

func (exec *SimpleExecutor) Begin() error {
    if exec.closed {
        return errors.EXECUTOR_BEGIN_ERROR
    }

    return exec.transaction.Begin()
}

func (exec *SimpleExecutor) Commit(require bool) error {
    if exec.closed {
        return errors.EXECUTOR_COMMIT_ERROR
    }

    if require {
        return exec.transaction.Commit()
    }

    return nil
}

func (exec *SimpleExecutor) Rollback(require bool) error {
    if !exec.closed {
        if require {
            return exec.transaction.Rollback()
        }
    }
    return nil
}
