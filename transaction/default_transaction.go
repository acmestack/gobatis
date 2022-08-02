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

package transaction

import (
	"context"
	"database/sql"

	"github.com/acmestack/gobatis/common"
	"github.com/acmestack/gobatis/connection"
	"github.com/acmestack/gobatis/datasource"
	"github.com/acmestack/gobatis/errors"
	"github.com/acmestack/gobatis/reflection"
	"github.com/acmestack/gobatis/statement"
	"github.com/acmestack/gobatis/util"
)

type DefaultTransaction struct {
	ds datasource.DataSource
	db *sql.DB
	tx *sql.Tx
}

func NewDefaultTransaction(ds datasource.DataSource, db *sql.DB) *DefaultTransaction {
	ret := &DefaultTransaction{ds: ds, db: db}
	return ret
}

func (trans *DefaultTransaction) GetConnection() connection.Connection {
	if trans.tx == nil {
		return (*connection.DefaultConnection)(trans.db)
	} else {
		return &TransactionConnection{tx: trans.tx}
	}
}

func (trans *DefaultTransaction) Close() {

}

func (trans *DefaultTransaction) Begin() error {
	tx, err := trans.db.Begin()
	if err != nil {
		return err
	}
	trans.tx = tx
	return nil
}

func (trans *DefaultTransaction) Commit() error {
	if trans.tx == nil {
		return errors.TransactionWithoutBegin
	}

	err := trans.tx.Commit()
	if err != nil {
		return errors.TransactionCommitError
	}
	return nil
}

func (trans *DefaultTransaction) Rollback() error {
	if trans.tx == nil {
		return errors.TransactionWithoutBegin
	}

	err := trans.tx.Rollback()
	if err != nil {
		return errors.TransactionCommitError
	}
	return nil
}

type TransactionConnection struct {
	tx *sql.Tx
}

type TransactionStatement struct {
	tx  *sql.Tx
	sql string
}

func (transConnection *TransactionConnection) Prepare(sqlStr string) (statement.Statement, error) {
	ret := &TransactionStatement{
		tx:  transConnection.tx,
		sql: sqlStr,
	}
	return ret, nil
}

func (transConnection *TransactionConnection) Query(ctx context.Context, result reflection.Object, sqlStr string, params ...interface{}) error {
	db := transConnection.tx
	rows, err := db.QueryContext(ctx, sqlStr, params...)
	if err != nil {
		return errors.StatementQueryError
	}
	defer rows.Close()

	util.ScanRows(rows, result)
	return nil
}

func (transConnection *TransactionConnection) Exec(ctx context.Context, sqlStr string, params ...interface{}) (common.Result, error) {
	db := transConnection.tx
	return db.ExecContext(ctx, sqlStr, params...)
}

func (transStatement *TransactionStatement) Query(ctx context.Context, result reflection.Object, params ...interface{}) error {
	rows, err := transStatement.tx.QueryContext(ctx, transStatement.sql, params...)
	if err != nil {
		return errors.StatementQueryError
	}
	defer rows.Close()

	util.ScanRows(rows, result)
	return nil
}

func (transStatement *TransactionStatement) Exec(ctx context.Context, params ...interface{}) (common.Result, error) {
	return transStatement.tx.ExecContext(ctx, transStatement.sql, params...)
}

func (transStatement *TransactionStatement) Close() {
	//Will be closed when commit or rollback
}
