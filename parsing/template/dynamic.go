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

package template

import (
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/acmestack/gobatis/parsing/sqlparser"
)

const (
	argPlaceHolder       = "_xfali_Arg_Holder"
	argPlaceHolderLen    = 17
	argPlaceHolderFormat = "%s%08d"

	FuncNameSet   = "set"
	FuncNameWhere = "where"
	FuncNameArg   = "arg"
	FuncNameAdd   = "add"
)

type Dynamic interface {
	getFuncMap() template.FuncMap
	format(string) (string, []any)
}

var ArgPlaceHolderFormat = argPlaceHolderFormat

func dummyUpdateSet(b any, column string, value any, origin string) string {
	return origin
}

func dummyWhere(b any, cond, column string, value any, origin string) string {
	return origin
}

//return as fast as possible
func dummyParam(p any) string {
	return ""
}

func dummyNil(p any) bool {
	return true
}

func commonAdd(a, b int) int {
	return a + b
}

type DummyDynamic struct{}

var dummyFuncMap = template.FuncMap{
	FuncNameSet:   dummyUpdateSet,
	FuncNameWhere: dummyWhere,
	FuncNameArg:   dummyParam,

	FuncNameAdd: commonAdd,
}

var gDummyDynamic = &DummyDynamic{}

func (dummyDynamic *DummyDynamic) getFuncMap() template.FuncMap {
	return dummyFuncMap
}

func (dummyDynamic *DummyDynamic) getParam() []any {
	return nil
}

func (dummyDynamic *DummyDynamic) format(s string) (string, []any) {
	return s, nil
}

type CommonDynamic struct {
	index    int
	keys     []string
	paramMap map[string]any
	holder   sqlparser.Holder
}

func CreateDynamicHandler(holder sqlparser.Holder) Dynamic {
	return &CommonDynamic{
		index:    0,
		keys:     nil,
		paramMap: map[string]any{},
		holder:   holder,
	}
}

func (dynamic *CommonDynamic) getFuncMap() template.FuncMap {
	return template.FuncMap{
		FuncNameSet:   dynamic.UpdateSet,
		FuncNameWhere: dynamic.Where,
		FuncNameArg:   dynamic.Param,

		FuncNameAdd: commonAdd,
	}
}

func (dynamic *CommonDynamic) UpdateSet(b any, columnDesc string, value any, origin string) string {
	if !IsTrue(b) {
		return origin
	}

	buf := strings.Builder{}
	if origin == "" {
		buf.WriteString(" SET ")
	} else {
		origin = strings.TrimSpace(origin)
		buf.WriteString(origin)
		if origin[:len(origin)-1] != "," {
			buf.WriteString(",")
		}
	}
	buf.WriteString(columnDesc)
	if s, ok := value.(string); ok {
		if _, ok := dynamic.paramMap[s]; ok {
			buf.WriteString(s)
		} else {
			buf.WriteString(`'`)
			buf.WriteString(s)
			buf.WriteString(`'`)
		}
	} else {
		buf.WriteString(fmt.Sprint(value))
	}
	return buf.String()
}

func (dynamic *CommonDynamic) Where(b any, cond, columnDesc string, value any, origin string) string {
	if !IsTrue(b) {
		return origin
	}

	buf := strings.Builder{}
	if origin == "" {
		buf.WriteString(" WHERE ")
		cond = ""
	} else {
		buf.WriteString(strings.TrimSpace(origin))
		buf.WriteString(" ")
		buf.WriteString(cond)
		buf.WriteString(" ")
	}

	buf.WriteString(columnDesc)
	if s, ok := value.(string); ok {
		if _, ok := dynamic.paramMap[s]; ok {
			buf.WriteString(s)
		} else {
			buf.WriteString(`'`)
			buf.WriteString(s)
			buf.WriteString(`'`)
		}
	} else {
		buf.WriteString(fmt.Sprint(value))
	}
	return buf.String()
}

func (dynamic *CommonDynamic) getParam() []any {
	return nil
}

func (dynamic *CommonDynamic) Param(p any) string {
	dynamic.index++
	key := getPlaceHolderKey(dynamic.index)
	dynamic.paramMap[key] = p
	dynamic.keys = append(dynamic.keys, key)
	return key
}

func (dynamic *CommonDynamic) format(s string) (string, []any) {
	i, index := 0, 1
	var params []any
	for _, k := range dynamic.keys {
		s, i = replace(s, k, dynamic.holder(index), -1)
		if i > 0 {
			params = append(params, dynamic.paramMap[k])
			index++
		}
	}
	return s, params
}

func selectDynamic(driverName string) Dynamic {
	if h, ok := sqlparser.GetMarker(driverName); ok {
		return dynamicFac(h)
	}
	return gDummyDynamic
}

func replace(s, old, new string, n int) (string, int) {
	if old == new || n == 0 {
		return s, 0 // avoid allocation
	}

	if old == "" {
		return s, 0
	}

	if n < 0 {
		if m := strings.Count(s, old); m == 0 {
			return s, 0 // avoid allocation
		} else if n < 0 || m < n {
			n = m
		}
	}
	makeSize := len(s) + n*(len(new)-len(old))
	// Apply replacements to buffer.
	t := make([]byte, makeSize)
	w, count := 0, 0
	start := 0
	for {
		if n == 0 {
			break
		}
		j := start
		index := strings.Index(s[start:], old)
		if index == -1 {
			return string(t[0:w]), count
		} else {
			j += index
			count++
		}
		w += copy(t[w:], s[start:j])
		w += copy(t[w:], new)
		start = j + len(old)
		n--
	}
	w += copy(t[w:], s[start:])
	return string(t[0:w]), count
}

func IsTrue(i any) bool {
	t, _ := template.IsTrue(i)
	if !t {
		return t
	}

	if ti, ok := i.(time.Time); ok {
		if ti.IsZero() {
			return false
		}
	}
	return t
}

func getPlaceHolderKey(index int) string {
	return fmt.Sprintf(ArgPlaceHolderFormat, argPlaceHolder, index)
}

var dynamicFac = CreateDynamicHandler

func SetDynamicFactory(f func(h sqlparser.Holder) Dynamic) {
	dynamicFac = f
}
