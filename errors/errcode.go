/*
 * Licensed to the AcmeStack under one or more contributor license
 * agreements. See the NOTICE file distributed with this work for
 * additional information regarding copyright ownership.
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
	FactoryInited               = New("10002", "Factory have been initialized")
	ParseModelTableinfoFailed   = New("11001", "Parse Model's table info failed")
	ModelNotRegister            = New("11002", "Register model not found")
	ObjectNotSupport            = New("11101", "Object not support")
	ParseObjectNotStruct        = New("11102", "Parse interface's info but not a struct")
	ParseObjectNotSlice         = New("11103", "Parse interface's info but not a slice")
	ParseObjectNotMap           = New("11104", "Parse interface's info but not a map")
	ParseObjectNotSimpletype    = New("11105", "Parse interface's info but not a simple type")
	SliceSliceNotSupport        = New("11106", "Parse interface's info: [][]slice not support")
	GetObjectinfoFailed         = New("11121", "Parse interface's info failed")
	SqlIdDuplicates             = New("11205", "Sql id is duplicates")
	DeserializeFailed           = New("11206", "Deserialize value failed")
	ParseSqlVarError            = New("12001", "SQL PARSE ERROR")
	ParseSqlParamError          = New("12002", "SQL PARSE parameter error")
	ParseSqlParamVarNumberError = New("12003", "SQL PARSE parameter var number error")
	ParseParserNilError         = New("12004", "Dynamic sql parser is nil error")
	ParseDynamicSqlError        = New("12010", "Parse dynamic sql error")
	ParseTemplateNilError       = New("12101", "Parse template is nil")
	ExecutorCommitError         = New("21001", "executor was closed when transaction commit")
	ExecutorBeginError          = New("21002", "executor was closed when transaction begin")
	ExecutorQueryError          = New("21003", "executor was closed when exec sql")
	ExecutorGetConnectionError  = New("21003", "executor get connection error")
	TransactionWithoutBegin     = New("22001", "Transaction without begin")
	TransactionCommitError      = New("22002", "Transaction commit error")
	TransactionBusinessError    = New("22003", "Business error in transaction")
	ConnectionPrepareError      = New("23001", "Connection prepare error")
	StatementQueryError         = New("24001", "statement query error")
	StatementExecError          = New("24002", "statement exec error")
	QueryTypeError              = New("25001", "select data convert error")
	ResultPointerIsNil          = New("31000", "result type is a nil pointer")
	ResultIsnotPointer          = New("31001", "result type is not pointer")
	ResultPtrValueIsPointer     = New("31002", "result type is pointer of pointer")
	RunnerNotReady              = New("31003", "Runner not ready, may sql or param have some error")
	ResultNameNotFound          = New("31004", "result name not found")
	ResultSelectEmptyValue      = New("31005", "select return empty value")
	ResultSetValueFailed        = New("31006", "result set value failed")
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
