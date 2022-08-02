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
	"github.com/acmestack/gobatis/datasource"
	"github.com/acmestack/gobatis/executor"
	"github.com/acmestack/gobatis/logging"
	"github.com/acmestack/gobatis/session"
	"github.com/acmestack/gobatis/transaction"
)

type Factory interface {
	Open(datasource.DataSource) error
	Close() error

	GetDataSource() datasource.DataSource

	CreateTransaction() transaction.Transaction
	CreateExecutor(transaction.Transaction) executor.Executor
	CreateSession() session.SqlSession
	LogFunc() logging.LogFunc
}
