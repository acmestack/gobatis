/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package transaction

import (
    "context"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "github.com/xfali/gobatis"
    "github.com/xfali/gobatis/common"
    "github.com/xfali/gobatis/connection"
    "github.com/xfali/gobatis/datasource"
    "github.com/xfali/gobatis/errors"
    "github.com/xfali/gobatis/handler"
    "github.com/xfali/gobatis/statement"
    "github.com/xfali/gobatis/util"
)

type MysqlTransaction struct {
    ds datasource.DataSource
    db *sql.DB
    tx *sql.Tx
}

func NewMysqlTransaction(ds datasource.DataSource, db *sql.DB) *MysqlTransaction {
    ret := &MysqlTransaction{ds: ds, db: db}
    return ret
}

func (trans *MysqlTransaction) GetConnection() connection.Connection {
    if trans.tx == nil {
        return (*connection.MysqlConnection)(trans.db)
    } else {
        return &TansactionConnection{tx: trans.tx}
    }
}

func (trans *MysqlTransaction) Close() {

}

func (trans *MysqlTransaction) Begin() error {
    tx, err := trans.db.Begin()
    if err != nil {
        return err
    }
    trans.tx = tx
    return nil
}

func (trans *MysqlTransaction) Commit() error {
    if trans.tx == nil {
        return errors.TRANSACTION_WITHOUT_BEGIN
    }

    err := trans.tx.Commit()
    if err != nil {
        return errors.TRANSACTION_COMMIT_ERROR
    }
    return nil
}

func (trans *MysqlTransaction) Rollback() error {
    if trans.tx == nil {
        return errors.TRANSACTION_WITHOUT_BEGIN
    }

    err := trans.tx.Rollback()
    if err != nil {
        return errors.TRANSACTION_COMMIT_ERROR
    }
    return nil
}

type TansactionConnection struct {
    tx *sql.Tx
}

type TransactionStatement struct {
    tx  *sql.Tx
    sql string
}

func (c *TansactionConnection) Prepare(sqlStr string) (statement.Statement, error) {
    ret := &TransactionStatement{
        tx:  c.tx,
        sql: sqlStr,
    }
    return ret, nil
}

func (c *TansactionConnection) Query(ctx context.Context, handler handler.ResultHandler, iterFunc gobatis.IterFunc, sqlStr string, params ...interface{}) error {
    db := c.tx
    rows, err := db.QueryContext(ctx, sqlStr, params...)
    if err != nil {
        return errors.STATEMENT_QUERY_ERROR
    }
    defer rows.Close()

    util.ScanRows(rows, handler, iterFunc)
    return nil
}

func (c *TansactionConnection) Exec(ctx context.Context, sqlStr string, params ...interface{}) (common.Result, error) {
    db := c.tx
    return db.ExecContext(ctx, sqlStr, params...)
}

func (s *TransactionStatement) Query(ctx context.Context, handler handler.ResultHandler, iterFunc gobatis.IterFunc, params ...interface{}) error {
    rows, err := s.tx.QueryContext(ctx, s.sql, params...)
    if err != nil {
        return errors.STATEMENT_QUERY_ERROR
    }
    defer rows.Close()

    util.ScanRows(rows, handler, iterFunc)
    return nil
}

func (s *TransactionStatement) Exec(ctx context.Context, params ...interface{}) (common.Result, error) {
    return s.tx.ExecContext(ctx, s.sql, params...)
}

func (s *TransactionStatement) Close() {
    //Will be closed when commit or rollback
}
