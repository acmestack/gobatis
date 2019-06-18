/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description:
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

type MysqlConnection sql.DB
type MysqlStatement sql.Stmt

func (c *MysqlConnection) Prepare(sqlStr string) (statement.Statement, error) {
	db := (*sql.DB)(c)
	s, err := db.Prepare(sqlStr)
	if err != nil {
		return nil, errors.CONNECTION_PREPARE_ERROR
	}
	return (*MysqlStatement)(s), nil
}

func (c *MysqlConnection) Query(ctx context.Context, result reflection.Object, sqlStr string, params ...interface{}) error {
	db := (*sql.DB)(c)
	rows, err := db.QueryContext(ctx, sqlStr, params...)
	if err != nil {
		return errors.STATEMENT_QUERY_ERROR
	}
	defer rows.Close()

	util.ScanRows(rows, result)
	return nil
}

func (c *MysqlConnection) Exec(ctx context.Context, sqlStr string, params ...interface{}) (common.Result, error) {
	db := (*sql.DB)(c)
	return db.ExecContext(ctx, sqlStr, params...)
}

func (s *MysqlStatement) Query(ctx context.Context, result reflection.Object, params ...interface{}) error {
	stmt := (*sql.Stmt)(s)
	rows, err := stmt.QueryContext(ctx, params...)
	if err != nil {
		return errors.STATEMENT_QUERY_ERROR
	}
	defer rows.Close()

	util.ScanRows(rows, result)
	return nil
}

func (s *MysqlStatement) Exec(ctx context.Context, params ...interface{}) (common.Result, error) {
	stmt := (*sql.Stmt)(s)
	return stmt.ExecContext(ctx, params...)
}

func (s *MysqlStatement) Close() {
	stmt := (*sql.Stmt)(s)
	stmt.Close()
}
