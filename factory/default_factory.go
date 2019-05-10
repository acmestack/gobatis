/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package factory

import (
    "database/sql"
    "github.com/xfali/gobatis/datasource"
    "github.com/xfali/gobatis/executor"
    "github.com/xfali/gobatis/logging"
    "github.com/xfali/gobatis/session"
    "github.com/xfali/gobatis/transaction"
)

type DefaultFactory struct {
    Host     string
    Port     int
    DBName   string
    Username string
    Password string
    Charset  string

    MaxConn     int
    MaxIdleConn int
    Log         logging.LogFunc

    ds          datasource.DataSource
    db          *sql.DB
}

func (f *DefaultFactory) InitDB() error {
    f.ds = &datasource.MysqlDataSource{
        Host:     f.Host,
        Port:     f.Port,
        DBName:   f.DBName,
        Username: f.Username,
        Password: f.Password,
        Charset:  f.Charset,
    }

    db, err := sql.Open(f.ds.DriverName(), f.ds.Url())
    db.SetMaxOpenConns(f.MaxConn)
    db.SetMaxIdleConns(f.MaxIdleConn)
    if err != nil {
        return err
    }

    f.db = db
    return nil
}

func (f *DefaultFactory) CreateSession() session.SqlSession {
    tx := transaction.NewMysqlTransaction(f.ds, f.db)
    return session.NewDefaultSqlSession(f.Log, tx, executor.NewSimpleExecutor(tx), false)
}

func (f *DefaultFactory) LogFunc() logging.LogFunc {
    return f.Log
}
