/*
 * Licensed to the AcmeStack under one or more contributor license
 * agreements. See the NOTICE file distributed with this work for
 * additional information regarding copyright ownership.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package session

import (
	"context"
	"fmt"

	"github.com/acmestack/gobatis/common"
	"github.com/acmestack/gobatis/executor"
	"github.com/acmestack/gobatis/logging"
	"github.com/acmestack/gobatis/reflection"
	"github.com/acmestack/gobatis/transaction"
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

func (session *DefaultSqlSession) Close(rollback bool) {
	session.executor.Close(rollback)
}

func (session *DefaultSqlSession) Query(ctx context.Context, result reflection.Object, sql string, params ...interface{}) error {
	session.logLastSql(sql, params...)
	return session.executor.Query(ctx, result, sql, params...)
}

func (session *DefaultSqlSession) Insert(ctx context.Context, sql string, params ...interface{}) (int64, int64, error) {
	session.logLastSql(sql, params...)
	ret, err := session.exec(ctx, sql, params...)
	if err != nil {
		return 0, -1, err
	}

	count, err := ret.RowsAffected()
	if err != nil {
		return 0, -1, err
	}

	id, err := ret.LastInsertId()
	if err != nil {
		return count, id, err
	}
	return count, id, nil
}

func (session *DefaultSqlSession) Update(ctx context.Context, sql string, params ...interface{}) (int64, error) {
	session.logLastSql(sql, params...)
	ret, err := session.exec(ctx, sql, params...)
	if err != nil {
		return 0, err
	}

	count, err := ret.RowsAffected()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (session *DefaultSqlSession) Delete(ctx context.Context, sql string, params ...interface{}) (int64, error) {
	session.logLastSql(sql, params...)
	ret, err := session.exec(ctx, sql, params...)
	if err != nil {
		return 0, err
	}

	count, err := ret.RowsAffected()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (session *DefaultSqlSession) Begin() error {
	session.logLastSql("Begin", "")
	return session.tx.Begin()
}

func (session *DefaultSqlSession) Commit() error {
	session.logLastSql("Commit", "")
	return session.tx.Commit()
}

func (session *DefaultSqlSession) Rollback() error {
	session.logLastSql("Rollback", "")
	return session.tx.Rollback()
}

func (session *DefaultSqlSession) logLastSql(sql string, params ...interface{}) {
	session.Log(logging.INFO, "sql: [%s], param: %s\n", sql, fmt.Sprint(params...))
}

func (session *DefaultSqlSession) exec(ctx context.Context, sql string, params ...interface{}) (common.Result, error) {
	return session.executor.Exec(ctx, sql, params...)
}
