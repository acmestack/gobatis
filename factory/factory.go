/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description:
 */

package factory

import (
	"github.com/xfali/gobatis/executor"
	"github.com/xfali/gobatis/logging"
	"github.com/xfali/gobatis/session"
	"github.com/xfali/gobatis/transaction"
)

type Factory interface {
	InitDB() error
	Close() error

	CreateTransaction() transaction.Transaction
	CreateExecutor(transaction.Transaction) executor.Executor
	CreateSession() session.SqlSession
	LogFunc() logging.LogFunc
}
