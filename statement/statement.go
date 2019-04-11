/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package statement

import (
    "github.com/xfali/gobatis/connection"
    "github.com/xfali/gobatis/handler"
)

type MappedStatement struct {
    Type          int
    Sql           string
    ResultHandler handler.ResultHandler
    IterFunc      connection.IterFunc
}
