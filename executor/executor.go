/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package executor

type Executor interface {
    Close(rollback bool)

    Query(statement *ExecParam, params ... interface{}) error

    Exec(statement *ExecParam, params ... interface{}) (int64, error)

    Begin() error

    Commit(require bool) error

    Rollback(require bool) error
}
