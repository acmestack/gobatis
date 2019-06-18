/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description:
 */

package statement

import (
	"context"
	"github.com/xfali/gobatis/common"
	"github.com/xfali/gobatis/reflection"
)

type Statement interface {
	Query(ctx context.Context, result reflection.Object, params ...interface{}) error
	Exec(ctx context.Context, params ...interface{}) (common.Result, error)
	Close()
}
