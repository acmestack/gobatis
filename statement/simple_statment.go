/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package statement

import "github.com/xfali/gobatis/handler"

type SimpleStatment struct {
    sql           string
    resultHandler handler.ResultHandler
}

func NewSimpleStatment(sql string, resultHandler handler.ResultHandler) *SimpleStatment {
    return &SimpleStatment{sql: sql, resultHandler: resultHandler}
}

func (this *SimpleStatment) Type() int {
    return 1
}

func (this *SimpleStatment) Sql() string {
    return this.sql
}
