/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package errors

import "fmt"

type ErrCode struct {
    Code    string `json:"code"`
    Message string `json:"msg"`
    fmtErr  string `json:"-"`
}

var EXECUTOR_COMMIT_ERROR *ErrCode = New("11001", "executor was closed when transaction commit")
var EXECUTOR_BEGIN_ERROR *ErrCode = New("11002", "executor was closed when transaction begin")
var EXECUTOR_QUERY_ERROR *ErrCode = New("11003", "executor was closed when exec sql")
var EXECUTOR_GET_CONNECTION_ERROR *ErrCode = New("11003", "executor get connection error")
var TRANSACTION_WITHOUT_BEGIN *ErrCode = New("12001", "Transaction without begin")
var TRANSACTION_COMMIT_ERROR *ErrCode = New("12002", "Transaction commit error")
var CONNECTION_PREPARE_ERROR *ErrCode = New("13001", "Connection prepare error")
var STATEMENT_QUERY_ERROR *ErrCode = New("14001", "statement query error")
var STATEMENT_EXEC_ERROR *ErrCode = New("14002", "statement exec error")
var QUERY_TYPE_ERROR *ErrCode = New("15001", "select data convert error")

func New(code, message string) *ErrCode {
    ret := &ErrCode{
        Code: code,
        Message: message,
        fmtErr: fmt.Sprintf("{ \"code\" : \"%s\", \"msg\" : \"%s\"", code, message),
    }
    return ret
}

func (e *ErrCode) Error() string {
    return e.fmtErr
}
