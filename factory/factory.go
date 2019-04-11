/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package factory

import (
    "github.com/xfali/gobatis/session"
)

type Factory interface {
    CreateSession() session.Session
}
