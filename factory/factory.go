/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package factory

import (
    "github.com/xfali/gobatis/logging"
    "github.com/xfali/gobatis/session"
)

type Factory interface {
    InitDB() error
    CreateSession() session.SqlSession
    LogFunc() logging.LogFunc
}
