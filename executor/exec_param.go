/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package executor

import (
    "github.com/xfali/gobatis"
    "github.com/xfali/gobatis/handler"
)

type ExecParam struct {
    Type          int
    Sql           string
    ResultHandler handler.ResultHandler
    IterFunc      gobatis.IterFunc

}
