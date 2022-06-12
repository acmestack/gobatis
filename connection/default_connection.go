/*
 * Copyright (c) 2022, OpeningO
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
	"github.com/xfali/gobatis/common"
	"github.com/xfali/gobatis/errors"
	"github.com/xfali/gobatis/reflection"
	"github.com/xfali/gobatis/statement"
	"github.com/xfali/gobatis/util"
)

type DefaultConnection sql.DB
type DefaultStatement sql.Stmt

func (c *DefaultConnection) Prepare(sqlStr string) (statement.Statement, error) {
	db := (*sql.DB)(c)
	s, err := db.Prepare(sqlStr)
	if err != nil {
		return nil, errors.CONNECTION_PREPARE_ERROR
	}
	return (*DefaultStatement)(s), nil
}

func (c *DefaultConnection) Query(ctx context.Context, result reflection.Object, sqlStr string, params ...interface{}) error {
	db := (*sql.DB)(c)
	rows, err := db.QueryContext(ctx, sqlStr, params...)
	if err != nil {
		return errors.STATEMENT_QUERY_ERROR
	}
	defer rows.Close()

	util.ScanRows(rows, result)
	return nil
}

func (c *DefaultConnection) Exec(ctx context.Context, sqlStr string, params ...interface{}) (common.Result, error) {
	db := (*sql.DB)(c)
	return db.ExecContext(ctx, sqlStr, params...)
}

func (s *DefaultStatement) Query(ctx context.Context, result reflection.Object, params ...interface{}) error {
	stmt := (*sql.Stmt)(s)
	rows, err := stmt.QueryContext(ctx, params...)
	if err != nil {
		return errors.STATEMENT_QUERY_ERROR
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
