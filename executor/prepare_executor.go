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
    "github.com/xfali/gobatis/transaction"
)

type PrepareExecutor struct {
    transaction transaction.Transaction
    closed      bool
}

func NewPrepareExecutor(transaction transaction.Transaction) *PrepareExecutor {
    return &PrepareExecutor{transaction: transaction}
}

func (exec *PrepareExecutor) Close(rollback bool) {
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

func (exec *PrepareExecutor) Query(execParam *ExecParam, params ...interface{}) error {
    if exec.closed {
        return  errors.EXECUTOR_QUERY_ERROR
    }

    conn := exec.transaction.GetConnection()
    if conn == nil {
        return errors.EXECUTOR_GET_CONNECTION_ERROR
    }

    //FIXME: stmt must be close, use stmtCache instead
    stmt, err := conn.Prepare(execParam.Sql)
    defer stmt.Close()
    if err != nil {
        return err
    }
    return stmt.Query(execParam.ResultHandler, execParam.IterFunc, params...)
}

func (exec *PrepareExecutor) Exec(execParam *ExecParam, params ...interface{}) (int64, error) {
    if exec.closed {
        return 0, errors.EXECUTOR_QUERY_ERROR
    }

    conn := exec.transaction.GetConnection()
    if conn == nil {
        return 0, errors.EXECUTOR_GET_CONNECTION_ERROR
    }

    //FIXME: stmt must be close, use stmtCache instead
    stmt, err := conn.Prepare(execParam.Sql)
    defer stmt.Close()

    if err != nil {
        return 0, err
    }
    return stmt.Exec(params...)
}

func (exec *PrepareExecutor) Begin() error {
    if exec.closed {
        return errors.EXECUTOR_BEGIN_ERROR
    }

    return exec.transaction.Begin()
}

func (exec *PrepareExecutor) Commit(require bool) error {
    if exec.closed {
        return errors.EXECUTOR_COMMIT_ERROR
    }

    if require {
        return exec.transaction.Commit()
    }

    return nil
}

func (exec *PrepareExecutor) Rollback(require bool) error {
    if !exec.closed {
        if require {
            return exec.transaction.Rollback()
        }
    }
    return nil
}
