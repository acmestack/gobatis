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
	"github.com/xfali/gobatis/errors"
	"github.com/xfali/gobatis/executor"
	"github.com/xfali/gobatis/logging"
	"github.com/xfali/gobatis/session"
	"github.com/xfali/gobatis/transaction"
	"sync"
	"time"
)

type DefaultFactory struct {
	// Deprecated: Use gobatis.NewFactory instead.
	Host string
	// Deprecated: Use gobatis.NewFactory instead.
	Port int
	// Deprecated: Use gobatis.NewFactory instead.
	DBName string
	// Deprecated: Use gobatis.NewFactory instead.
	Username string
	// Deprecated: Use gobatis.NewFactory instead.
	Password string
	// Deprecated: Use gobatis.NewFactory instead.
	Charset string

	MaxConn         int
	MaxIdleConn     int
	ConnMaxLifetime time.Duration
	Log             logging.LogFunc

	DataSource datasource.DataSource

	db    *sql.DB
	mutex sync.Mutex
}

func (f *DefaultFactory) Open(ds datasource.DataSource) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if f.db != nil {
		return errors.FACTORY_INITED
	}

	if ds != nil {
		f.DataSource = ds
	}

	db, err := sql.Open(f.DataSource.DriverName(), f.DataSource.DriverInfo())
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(f.MaxConn)
	db.SetMaxIdleConns(f.MaxIdleConn)
	db.SetConnMaxLifetime(f.ConnMaxLifetime)

	f.db = db
	return nil
}

func (f *DefaultFactory) Close() error {
	if f.db != nil {
		return f.db.Close()
	}
	return nil
}

func (f *DefaultFactory) GetDataSource() datasource.DataSource {
	return f.DataSource
}

func (f *DefaultFactory) CreateTransaction() transaction.Transaction {
	return transaction.NewDefaultTransaction(f.DataSource, f.db)
}

func (f *DefaultFactory) CreateExecutor(transaction transaction.Transaction) executor.Executor {
	return executor.NewSimpleExecutor(transaction)
}

func (f *DefaultFactory) CreateSession() session.SqlSession {
	tx := f.CreateTransaction()
	return session.NewDefaultSqlSession(f.Log, tx, f.CreateExecutor(tx), false)
}

func (f *DefaultFactory) LogFunc() logging.LogFunc {
	return f.Log
}

func (f *DefaultFactory) WithLock(lockFunc func(fac *DefaultFactory)) {
	f.mutex.Lock()
	lockFunc(f)
	f.mutex.Unlock()
}
