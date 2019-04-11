/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package executor

import "github.com/xfali/gobatis/statement"

type Executor interface {
    Close(rollback bool)

    Query(statement *statement.MappedStatement, params ... interface{}) error

    Exec(statement *statement.MappedStatement, params ... interface{}) (int64, error)

    Begin() error

    Commit(require bool) error

    Rollback(require bool) error
}
