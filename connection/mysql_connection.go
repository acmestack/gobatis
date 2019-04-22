/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package connection

import (
    "database/sql"
    "github.com/xfali/gobatis"
    "github.com/xfali/gobatis/errors"
    "github.com/xfali/gobatis/handler"
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

func (c *MysqlConnection)Query(handler handler.ResultHandler, iterFunc gobatis.IterFunc, sqlStr string, params ...interface{}) error {
    db := (*sql.DB)(c)
    rows, err := db.Query(sqlStr, params...)
    if err != nil {
        return errors.STATEMENT_QUERY_ERROR
    }
    defer rows.Close()

    util.ScanRows(rows, handler, iterFunc)
    return nil
}

func (c *MysqlConnection)Exec(sqlStr string, params ...interface{}) (int64, error) {
    db := (*sql.DB)(c)
    result, err := db.Exec(sqlStr, params...)
    if err != nil {
        return 0, errors.STATEMENT_QUERY_ERROR
    }
    ret, err := result.RowsAffected()
    if err != nil {
        return 0, errors.STATEMENT_QUERY_ERROR
    }
    return ret, nil
}

func (s *MysqlStatement) Query(handler handler.ResultHandler, iterFunc gobatis.IterFunc, params ...interface{}) error {
    stmt := (*sql.Stmt)(s)
    rows, err := stmt.Query(params...)
    if err != nil {
        return errors.STATEMENT_QUERY_ERROR
    }
    defer rows.Close()

    util.ScanRows(rows, handler, iterFunc)
    return nil
}

func (s *MysqlStatement) Exec(params ...interface{}) (int64, error) {
    stmt := (*sql.Stmt)(s)
    result, err := stmt.Exec(params...)
    if err != nil {
        return 0, errors.STATEMENT_QUERY_ERROR
    }
    ret, err := result.RowsAffected()
    if err != nil {
        return 0, errors.STATEMENT_QUERY_ERROR
    }
    return ret, nil
}

func (s *MysqlStatement) Close() {
    stmt := (*sql.Stmt)(s)
    stmt.Close()
}
