/*
 * Copyright (c) 2022, AcmeStack
 * All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package errors

import "fmt"

type ErrCode struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
	fmtErr  string `json:"-"`
}

var (
	FACTORY_INITED                   = New("10002", "Factory have been initialized")
	PARSE_MODEL_TABLEINFO_FAILED     = New("11001", "Parse Model's table info failed")
	MODEL_NOT_REGISTER               = New("11002", "Register model not found")
	OBJECT_NOT_SUPPORT               = New("11101", "Object not support")
	PARSE_OBJECT_NOT_STRUCT          = New("11102", "Parse interface's info but not a struct")
	PARSE_OBJECT_NOT_SLICE           = New("11103", "Parse interface's info but not a slice")
	PARSE_OBJECT_NOT_MAP             = New("11104", "Parse interface's info but not a map")
	PARSE_OBJECT_NOT_SIMPLETYPE      = New("11105", "Parse interface's info but not a simple type")
	SLICE_SLICE_NOT_SUPPORT          = New("11106", "Parse interface's info: [][]slice not support")
	GET_OBJECTINFO_FAILED            = New("11121", "Parse interface's info failed")
	SQL_ID_DUPLICATES                = New("11205", "Sql id is duplicates")
	DESERIALIZE_FAILED               = New("11206", "Deserialize value failed")
	PARSE_SQL_VAR_ERROR              = New("12001", "SQL PARSE ERROR")
	PARSE_SQL_PARAM_ERROR            = New("12002", "SQL PARSE parameter error")
	PARSE_SQL_PARAM_VAR_NUMBER_ERROR = New("12003", "SQL PARSE parameter var number error")
	PARSE_PARSER_NIL_ERROR           = New("12004", "Dynamic sql parser is nil error")
	PARSE_DYNAMIC_SQL_ERROR          = New("12010", "Parse dynamic sql error")
	PARSE_TEMPLATE_NIL_ERROR         = New("12101", "Parse template is nil")
	EXECUTOR_COMMIT_ERROR            = New("21001", "executor was closed when transaction commit")
	EXECUTOR_BEGIN_ERROR             = New("21002", "executor was closed when transaction begin")
	EXECUTOR_QUERY_ERROR             = New("21003", "executor was closed when exec sql")
	EXECUTOR_GET_CONNECTION_ERROR    = New("21003", "executor get connection error")
	TRANSACTION_WITHOUT_BEGIN        = New("22001", "Transaction without begin")
	TRANSACTION_COMMIT_ERROR         = New("22002", "Transaction commit error")
	TRANSACTION_BUSINESS_ERROR       = New("22003", "Business error in transaction")
	CONNECTION_PREPARE_ERROR         = New("23001", "Connection prepare error")
	STATEMENT_QUERY_ERROR            = New("24001", "statement query error")
	STATEMENT_EXEC_ERROR             = New("24002", "statement exec error")
	QUERY_TYPE_ERROR                 = New("25001", "select data convert error")
	RESULT_POINTER_IS_NIL            = New("31000", "result type is a nil pointer")
	RESULT_ISNOT_POINTER             = New("31001", "result type is not pointer")
	RESULT_PTR_VALUE_IS_POINTER      = New("31002", "result type is pointer of pointer")
	RUNNER_NOT_READY                 = New("31003", "Runner not ready, may sql or param have some error")
	RESULT_NAME_NOT_FOUND            = New("31004", "result name not found")
	RESULT_SELECT_EMPTY_VALUE        = New("31005", "select return empty value")
	RESULT_SET_VALUE_FAILED          = New("31006", "result set value failed")
)

func New(code, message string) *ErrCode {
	ret := &ErrCode{
		Code:    code,
		Message: message,
		fmtErr:  fmt.Sprintf("{ \"code\" : \"%s\", \"msg\" : \"%s\" }", code, message),
	}
	return ret
}

func (e *ErrCode) Error() string {
	return e.fmtErr
}
