/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package connection

type Statement interface {
    Query(params ...interface{}) error
    Exec(params ...interface{}) (int64, error)
}

type Connection interface {
    Prepare(sql string) (Statement, error)
}
