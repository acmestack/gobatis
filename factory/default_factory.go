/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package factory

import (
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
    transaction transaction.Transaction
}

func (f *DefaultFactory) Init() {
    f.ds = &datasource.MysqlDataSource{
        Host:     f.Host,
        Port:     f.Port,
        DBName:   f.DBName,
        Username: f.Username,
        Password: f.Password,
        Charset:  f.Charset,
    }
    f.transaction = transaction.NewMysqlTransaction(f.ds, f.MaxConn, f.MaxIdleConn)
}

func (f *DefaultFactory) CreateSession() session.Session {
    return session.NewDefaultSqlSession(f.Log, f.transaction, executor.NewSimpleExecutor(f.transaction), false)
}

func (f *DefaultFactory) LogFunc() logging.LogFunc {
    return f.Log
}
