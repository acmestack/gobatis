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

var PARSE_MODEL_TABLEINFO_FAILED *ErrCode = New("11001", "Parse Model's table info failed")
var MODEL_NOT_REGISTER *ErrCode = New("11002", "Register model not found")
var PARSE_TABLEINFO_NOT_STRUCT *ErrCode = New("11003", "Parse interface's table info but not a struct")
var DESERIALIZE_FAILED *ErrCode = New("11004", "Deserialize value failed")
var PARSE_SQL_VAR_ERROR *ErrCode = New("12001", "SQL PARSE ERROR")
var PARSE_SQL_PARAM_ERROR *ErrCode = New("12002", "SQL PARSE parameter error")
var PARSE_SQL_PARAM_VAR_NUMBER_ERROR *ErrCode = New("12003", "SQL PARSE parameter var number error")
var EXECUTOR_COMMIT_ERROR *ErrCode = New("21001", "executor was closed when transaction commit")
var EXECUTOR_BEGIN_ERROR *ErrCode = New("21002", "executor was closed when transaction begin")
var EXECUTOR_QUERY_ERROR *ErrCode = New("21003", "executor was closed when exec sql")
var EXECUTOR_GET_CONNECTION_ERROR *ErrCode = New("21003", "executor get connection error")
var TRANSACTION_WITHOUT_BEGIN *ErrCode = New("22001", "Transaction without begin")
var TRANSACTION_COMMIT_ERROR *ErrCode = New("22002", "Transaction commit error")
var CONNECTION_PREPARE_ERROR *ErrCode = New("23001", "Connection prepare error")
var STATEMENT_QUERY_ERROR *ErrCode = New("24001", "statement query error")
var STATEMENT_EXEC_ERROR *ErrCode = New("24002", "statement exec error")
var QUERY_TYPE_ERROR *ErrCode = New("25001", "select data convert error")
var RESULT_ISNOT_POINTER *ErrCode = New("31001", "result type is not pointer")
var RESULT_PTR_VALUE_IS_POINTER *ErrCode = New("31002", "result type is pointer of pointer")
var RUNNER_NOT_READY *ErrCode = New("31003", "Runner not ready, may sql or param have some error")
var RESULT_NAME_NOT_FOUND *ErrCode = New("31004", "result name not found")
var RESULT_SELECT_EMPTY_VALUE *ErrCode = New("31005", "select return empty value")

func New(code, message string) *ErrCode {
    ret := &ErrCode{
        Code: code,
        Message: message,
        fmtErr: fmt.Sprintf("{ \"code\" : \"%s\", \"msg\" : \"%s\" }", code, message),
    }
    return ret
}

func (e *ErrCode) Error() string {
    return e.fmtErr
}
