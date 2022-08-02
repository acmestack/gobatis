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

package sqlparser

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/acmestack/gobatis/errors"
)

const (
	SELECT = "select"
	INSERT = "insert"
	UPDATE = "update"
	DELETE = "delete"
)

type Metadata struct {
	Action     string
	PrepareSql string
	Vars       []string
	Params     []interface{}
}

type SqlParser interface {
	ParseMetadata(driverName string, params ...interface{}) (*Metadata, error)
}

func SimpleParse(sql string) (*Metadata, error) {
	ret := Metadata{}
	sql = strings.Trim(sql, " ")
	action := sql[:6]
	action = strings.ToLower(action)
	ret.Action = action

	subStr := sql
	firstIndex, lastIndex := -1, -1
	for {
		firstIndex = strings.Index(subStr, "#{")
		if firstIndex == -1 {
			break
		} else {
			subStr = subStr[firstIndex+2:]
			lastIndex = findFirst(subStr, '}')
			//lastIndex = strings.Index(subStr, "}")
			if lastIndex == -1 {
				return nil, errors.ParseSqlVarError
			} else {
				varName := subStr[:lastIndex]
				if varName != "" {
					ret.Vars = append(ret.Vars, varName)
				}
			}
		}
		subStr = subStr[lastIndex+1:]
	}

	ret.PrepareSql = sql
	for _, varName := range ret.Vars {
		ret.PrepareSql = strings.Replace(ret.PrepareSql, "#{"+varName+"}", "?", -1)
	}

	return &ret, nil
}

func ParseWithParams(sql string, params ...interface{}) (*Metadata, error) {
	ret := Metadata{}
	sql = strings.Trim(sql, " ")
	action := sql[:6]
	action = strings.ToLower(action)
	ret.Action = action

	ret.PrepareSql = sql
	subStr := sql
	firstIndex, lastIndex := -1, -1
	var c string
	for {
		firstIndex = strings.Index(subStr, "{")
		if firstIndex == -1 || firstIndex == 0 {
			break
		} else {
			c = subStr[firstIndex-1 : firstIndex]
			subStr = subStr[firstIndex+1:]
			lastIndex = findFirst(subStr, '}')
			//lastIndex = strings.Index(subStr, "}")
			if lastIndex == -1 {
				return nil, errors.ParseSqlVarError
			} else {
				varName := subStr[:lastIndex]
				if varName != "" {
					ret.Vars = append(ret.Vars, varName)
					indexV, err := strconv.Atoi(varName)
					if err != nil {
						return nil, errors.ParseSqlParamVarNumberError
					}
					if c == "$" {
						if len(params) <= indexV {
							return nil, errors.ParseSqlParamError
						}
						oldStr := "${" + varName + "}"
						newStr := interface2String(params[indexV])
						ret.PrepareSql = strings.Replace(ret.PrepareSql, oldStr, newStr, -1)
						subStr = strings.Replace(subStr, oldStr, newStr, -1)
					} else if c == "#" {
						if len(params) < indexV {
							return nil, errors.ParseSqlParamError
						}
						oldStr := "#{" + varName + "}"
						ret.PrepareSql = strings.Replace(ret.PrepareSql, oldStr, "?", -1)
						ret.Params = append(ret.Params, params[indexV])
					}
				}
			}
		}
		subStr = subStr[lastIndex+1:]
	}

	return &ret, nil
}

func ParseWithParamMap(driverName, sql string, params map[string]interface{}) (*Metadata, error) {
	ret := Metadata{}
	sql = strings.Trim(sql, " ")
	action := sql[:6]
	action = strings.ToLower(action)
	ret.Action = action

	ret.PrepareSql = sql
	subStr := sql
	firstIndex, lastIndex := -1, -1
	var c string
	var index int = 0
	holder := SelectMarker(driverName)

	for {
		firstIndex = strings.Index(subStr, "{")
		if firstIndex == -1 || firstIndex == 0 {
			break
		} else {
			c = subStr[firstIndex-1 : firstIndex]
			subStr = subStr[firstIndex+1:]
			lastIndex = findFirst(subStr, '}')
			//lastIndex = strings.Index(subStr, "}")
			if lastIndex == -1 {
				return nil, errors.ParseSqlVarError
			} else {
				varName := subStr[:lastIndex]
				if varName != "" {
					ret.Vars = append(ret.Vars, varName)
					if value, ok := params[varName]; ok {
						if c == "$" {
							oldStr := "${" + varName + "}"
							newStr := interface2String(value)
							ret.PrepareSql = strings.Replace(ret.PrepareSql, oldStr, newStr, -1)
							subStr = strings.Replace(subStr, oldStr, newStr, -1)
						} else if c == "#" {
							oldStr := "#{" + varName + "}"
							index++
							h := holder(index)
							ret.PrepareSql = strings.Replace(ret.PrepareSql, oldStr, h, -1)
							ret.Params = append(ret.Params, value)
						}
					} else {
						return nil, errors.ParseSqlParamError
					}
				}
			}
		}
		subStr = subStr[lastIndex+1:]
	}

	return &ret, nil
}

type Holder func(int) string

var gHolderMap = map[string]Holder{
	"mysql":    MysqlMarker,    //mysql
	"postgres": PostgresMarker, //postgresql
	"oci8":     Oci8Marker,     //oracle
	"adodb":    MysqlMarker,    //sqlserver
}

func RegisterParamMarker(driverName string, h Holder) bool {
	_, ok := GetMarker(driverName)
	gHolderMap[driverName] = h
	return ok
}

func SelectMarker(driverName string) Holder {
	if v, ok := GetMarker(driverName); ok {
		return v
	}
	return MysqlMarker
}

func GetMarker(driverName string) (Holder, bool) {
	v, ok := gHolderMap[driverName]
	return v, ok
}

func MysqlMarker(int) string {
	return "?"
}

func PostgresMarker(i int) string {
	return "$" + strconv.Itoa(i)
}

func Oci8Marker(i int) string {
	return ":" + strconv.Itoa(i)
}

func interface2String(i interface{}) string {
	return fmt.Sprintf("%v", i)
}

func findFirst(subStr string, char rune) int {
	for i, r := range subStr {
		//switch r {
		//case ',', ' ', '\t', '\n', '\r':
		//    return -1
		//case char:
		//    return i
		//}
		if unicode.IsSpace(r) || r == ',' {
			return -1
		} else if r == char {
			return i
		}
	}
	return -1
}

func (md *Metadata) String() string {
	return fmt.Sprintf("action: %s, prepareSql: %s, varmap: %v, params: %v", md.Action, md.PrepareSql, md.Vars, md.Params)
}
