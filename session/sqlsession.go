/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package session

import "github.com/xfali/GoBatis/handler"

type Session interface {
    Close(rollback bool)

    Select(handler handler.ResultHandler, sql string, params ... string) error

    Insert(sql string, params ... string) int64

    Update(sql string, params ... string) int64

    Delete(sql string, params ... string) int64

    Commit()

    Rollback()
}


