/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description:
 */

package sqlparser

import (
	"fmt"
	"github.com/xfali/gobatis/errors"
	"strconv"
	"strings"
	"unicode"
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
				return nil, errors.PARSE_SQL_VAR_ERROR
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
				return nil, errors.PARSE_SQL_VAR_ERROR
			} else {
				varName := subStr[:lastIndex]
				if varName != "" {
					ret.Vars = append(ret.Vars, varName)
					indexV, err := strconv.Atoi(varName)
					if err != nil {
						return nil, errors.PARSE_SQL_PARAM_VAR_NUMBER_ERROR
					}
					if c == "$" {
						if len(params) <= indexV {
							return nil, errors.PARSE_SQL_PARAM_ERROR
						}
						oldStr := "${" + varName + "}"
						newStr := interface2String(params[indexV])
						ret.PrepareSql = strings.Replace(ret.PrepareSql, oldStr, newStr, -1)
						subStr = strings.Replace(subStr, oldStr, newStr, -1)
					} else if c == "#" {
						if len(params) < indexV {
							return nil, errors.PARSE_SQL_PARAM_ERROR
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

func ParseWithParamMap(sql string, params map[string]interface{}) (*Metadata, error) {
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
				return nil, errors.PARSE_SQL_VAR_ERROR
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
							ret.PrepareSql = strings.Replace(ret.PrepareSql, oldStr, "?", -1)
							ret.Params = append(ret.Params, value)
						}
					} else {
						return nil, errors.PARSE_SQL_PARAM_ERROR
					}
				}
			}
		}
		subStr = subStr[lastIndex+1:]
	}

	return &ret, nil
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
