/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package transaction

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "github.com/xfali/gobatis/connection"
    "github.com/xfali/gobatis/datasource"
    "github.com/xfali/gobatis/errors"
)

type MysqlTransaction struct {
    ds datasource.DataSource
    db *sql.DB
    tx *sql.Tx
}

func NewMysqlTransaction(ds datasource.DataSource, maxConn, maxIdleConn int) *MysqlTransaction {
    db, err := sql.Open(ds.DriverName(), ds.Url())
    db.SetMaxOpenConns(maxConn)
    db.SetMaxIdleConns(maxIdleConn)
    if err != nil {
        return nil
    }
    ret := &MysqlTransaction{ds: ds, db: db}
    return ret
}

func (trans *MysqlTransaction) GetConnection() connection.Connection {
    return (*connection.MysqlConnection)(trans.db)
}

func (trans *MysqlTransaction) Close() {
    trans.db.Close()
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
