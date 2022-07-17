/*
 * Copyright (c) 2022, AcmeStack
 * All rights reserved.
 *
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

package connection

import (
	"context"
	"database/sql"

	"github.com/acmestack/gobatis/common"
	"github.com/acmestack/gobatis/errors"
	"github.com/acmestack/gobatis/reflection"
	"github.com/acmestack/gobatis/statement"
	"github.com/acmestack/gobatis/util"
)

type DefaultConnection sql.DB
type DefaultStatement sql.Stmt

func (conn *DefaultConnection) Prepare(sqlStr string) (statement.Statement, error) {
	db := (*sql.DB)(conn)
	s, err := db.Prepare(sqlStr)
	if err != nil {
		return nil, errors.ConnectionPrepareError
	}
	return (*DefaultStatement)(s), nil
}

func (conn *DefaultConnection) Query(ctx context.Context, result reflection.Object, sqlStr string, params ...interface{}) error {
	db := (*sql.DB)(conn)
	rows, err := db.QueryContext(ctx, sqlStr, params...)
	if err != nil {
		return errors.StatementQueryError
	}
	defer rows.Close()

	util.ScanRows(rows, result)
	return nil
}

func (conn *DefaultConnection) Exec(ctx context.Context, sqlStr string, params ...interface{}) (common.Result, error) {
	db := (*sql.DB)(conn)
	return db.ExecContext(ctx, sqlStr, params...)
}

func (s *DefaultStatement) Query(ctx context.Context, result reflection.Object, params ...interface{}) error {
	stmt := (*sql.Stmt)(s)
	rows, err := stmt.QueryContext(ctx, params...)
	if err != nil {
		return errors.StatementQueryError
	}
	defer rows.Close()

	util.ScanRows(rows, result)
	return nil
}

func (s *DefaultStatement) Exec(ctx context.Context, params ...interface{}) (common.Result, error) {
	stmt := (*sql.Stmt)(s)
	return stmt.ExecContext(ctx, params...)
}

func (s *DefaultStatement) Close() {
	stmt := (*sql.Stmt)(s)
	stmt.Close()
}
