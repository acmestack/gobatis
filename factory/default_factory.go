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

package factory

import (
	"database/sql"
	"sync"
	"time"

	"github.com/acmestack/gobatis/datasource"
	"github.com/acmestack/gobatis/errors"
	"github.com/acmestack/gobatis/executor"
	"github.com/acmestack/gobatis/logging"
	"github.com/acmestack/gobatis/session"
	"github.com/acmestack/gobatis/transaction"
)

type DefaultFactory struct {
	MaxConn         int
	MaxIdleConn     int
	ConnMaxLifetime time.Duration
	Log             logging.LogFunc

	DataSource datasource.DataSource

	db    *sql.DB
	mutex sync.Mutex
}

func (factory *DefaultFactory) Open(ds datasource.DataSource) error {
	factory.mutex.Lock()
	defer factory.mutex.Unlock()

	if factory.db != nil {
		return errors.FactoryInitialized
	}

	if ds != nil {
		factory.DataSource = ds
	}

	db, err := sql.Open(factory.DataSource.DriverName(), factory.DataSource.DriverInfo())
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(factory.MaxConn)
	db.SetMaxIdleConns(factory.MaxIdleConn)
	db.SetConnMaxLifetime(factory.ConnMaxLifetime)

	factory.db = db
	return nil
}

func (factory *DefaultFactory) Close() error {
	if factory.db != nil {
		return factory.db.Close()
	}
	return nil
}

func (factory *DefaultFactory) GetDataSource() datasource.DataSource {
	return factory.DataSource
}

func (factory *DefaultFactory) CreateTransaction() transaction.Transaction {
	return transaction.NewDefaultTransaction(factory.DataSource, factory.db)
}

func (factory *DefaultFactory) CreateExecutor(transaction transaction.Transaction) executor.Executor {
	return executor.NewSimpleExecutor(transaction)
}

func (factory *DefaultFactory) CreateSession() session.SqlSession {
	tx := factory.CreateTransaction()
	return session.NewDefaultSqlSession(factory.Log, tx, factory.CreateExecutor(tx), false)
}

func (factory *DefaultFactory) LogFunc() logging.LogFunc {
	return factory.Log
}

func (factory *DefaultFactory) WithLock(lockFunc func(fac *DefaultFactory)) {
	factory.mutex.Lock()
	lockFunc(factory)
	factory.mutex.Unlock()
}

// Deprecated: Use Open instead
func (factory *DefaultFactory) InitDB() error {
	return factory.Open(factory.DataSource)
}
