/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package transaction

import (
    "github.com/xfali/GoBatis/connection"
)

type Transaction interface {
    Close()

    GetConnection() connection.Connection

    Begin() error

    Commit() error

    Rollback() error
}
