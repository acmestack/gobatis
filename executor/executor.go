/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package executor

import (
    "context"
    "github.com/xfali/gobatis/common"
    "github.com/xfali/gobatis/reflection"
)

type Executor interface {
    Close(rollback bool)

    Query(ctx context.Context, result reflection.Object, sql string, params ... interface{}) error

    Exec(ctx context.Context, sql string, params ... interface{}) (common.Result, error)

    Begin() error

    Commit(require bool) error

    Rollback(require bool) error
}
