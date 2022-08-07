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

type errCode struct {
	code    string
	message string
}

var (
	FactoryInitialized          = gobatisError("10002", "Factory have been initialized")
	ParseModelTableInfoFailed   = gobatisError("11001", "Parse Model's table info failed")
	ModelNotRegister            = gobatisError("11002", "Register model not found")
	ObjectNotSupport            = gobatisError("11101", "Object not support")
	ParseObjectNotStruct        = gobatisError("11102", "Parse interface's info but not a struct")
	ParseObjectNotSlice         = gobatisError("11103", "Parse interface's info but not a slice")
	ParseObjectNotMap           = gobatisError("11104", "Parse interface's info but not a map")
	ParseObjectNotSimpletype    = gobatisError("11105", "Parse interface's info but not a simple type")
	SliceSliceNotSupport        = gobatisError("11106", "Parse interface's info: [][]slice not support")
	GetObjectInfoFailed         = gobatisError("11121", "Parse interface's info failed")
	SqlIdDuplicates             = gobatisError("11205", "Sql id is duplicates")
	DeserializeFailed           = gobatisError("11206", "Deserialize value failed")
	ParseSqlVarError            = gobatisError("12001", "SQL PARSE ERROR")
	ParseSqlParamError          = gobatisError("12002", "SQL PARSE parameter error")
	ParseSqlParamVarNumberError = gobatisError("12003", "SQL PARSE parameter var number error")
	ParseParserNilError         = gobatisError("12004", "Dynamic sql parser is nil error")
	ParseDynamicSqlError        = gobatisError("12010", "Parse dynamic sql error")
	ParseTemplateNilError       = gobatisError("12101", "Parse template is nil")
	ExecutorCommitError         = gobatisError("21001", "executor was closed when transaction commit")
	ExecutorBeginError          = gobatisError("21002", "executor was closed when transaction begin")
	ExecutorQueryError          = gobatisError("21003", "executor was closed when exec sql")
	ExecutorGetConnectionError  = gobatisError("21003", "executor get connection error")
	TransactionWithoutBegin     = gobatisError("22001", "Transaction without begin")
	TransactionCommitError      = gobatisError("22002", "Transaction commit error")
	TransactionBusinessError    = gobatisError("22003", "Business error in transaction")
	ConnectionPrepareError      = gobatisError("23001", "Connection prepare error")
	StatementQueryError         = gobatisError("24001", "statement query error")
	StatementExecError          = gobatisError("24002", "statement exec error")
	QueryTypeError              = gobatisError("25001", "select data convert error")
	ResultPointerIsNil          = gobatisError("31000", "result type is a nil pointer")
	ResultIsnotPointer          = gobatisError("31001", "result type is not pointer")
	ResultPtrValueIsPointer     = gobatisError("31002", "result type is pointer of pointer")
	RunnerNotReady              = gobatisError("31003", "Runner not ready, may sql or param have some error")
	ResultNameNotFound          = gobatisError("31004", "result name not found")
	ResultSelectEmptyValue      = gobatisError("31005", "select return empty value")
	ResultSetValueFailed        = gobatisError("31006", "result set value failed")
)

func gobatisError(code, message string) errCode {
	return errCode{
		code:    code,
		message: message,
	}
}

func (e errCode) Error() string {
	return fmt.Sprintf("%s - %s", e.code, e.message)
}
