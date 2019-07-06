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
	"github.com/xfali/gobatis/common"
	"github.com/xfali/gobatis/executor"
	"github.com/xfali/gobatis/logging"
	"github.com/xfali/gobatis/reflection"
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

func (sess *DefaultSqlSession) Query(ctx context.Context, result reflection.Object, sql string, params ...interface{}) error {
	sess.logLastSql(sql, params...)
	return sess.executor.Query(ctx, result, sql, params...)
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
	sess.Log(logging.INFO, "sql: [%s], param: %s\n", sql, fmt.Sprint(params))
}

func (sess *DefaultSqlSession) exec(ctx context.Context, sql string, params ...interface{}) (common.Result, error) {
	return sess.executor.Exec(ctx, sql, params...)
}
