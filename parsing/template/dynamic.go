// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package template

import (
	"fmt"
	"github.com/xfali/gobatis/parsing/sqlparser"
	"strings"
	"text/template"
	"time"
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
	format(string) (string, []interface{})
}

var ArgPlaceHolderFormat = argPlaceHolderFormat

func dummyUpdateSet(b interface{}, column string, value interface{}, origin string) string {
	return origin
}

func dummyWhere(b interface{}, cond, column string, value interface{}, origin string) string {
	return origin
}

//return as fast as possible
func dummyParam(p interface{}) string {
	return ""
}

func dummyNil(p interface{}) bool {
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

func (d *DummyDynamic) getFuncMap() template.FuncMap {
	return dummyFuncMap
}

func (d *DummyDynamic) getParam() []interface{} {
	return nil
}

func (d *DummyDynamic) format(s string) (string, []interface{}) {
	return s, nil
}

type CommonDynamic struct {
	index    int
	keys     []string
	paramMap map[string]interface{}
	holder   sqlparser.Holder
}

func CreateDynamicHandler(h sqlparser.Holder) Dynamic {
	return &CommonDynamic{
		index:    0,
		keys:     nil,
		paramMap: map[string]interface{}{},
		holder:   h,
	}
}

func (d *CommonDynamic) getFuncMap() template.FuncMap {
	return template.FuncMap{
		FuncNameSet:   d.UpdateSet,
		FuncNameWhere: d.Where,
		FuncNameArg:   d.Param,

		FuncNameAdd: commonAdd,
	}
}

func (d *CommonDynamic) UpdateSet(b interface{}, columnDesc string, value interface{}, origin string) string {
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
		if _, ok := d.paramMap[s]; ok {
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

func (d *CommonDynamic) Where(b interface{}, cond, columnDesc string, value interface{}, origin string) string {
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
		if _, ok := d.paramMap[s]; ok {
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

func (d *CommonDynamic) getParam() []interface{} {
	return nil
}

func (d *CommonDynamic) Param(p interface{}) string {
	d.index++
	key := getPlaceHolderKey(d.index)
	d.paramMap[key] = p
	d.keys = append(d.keys, key)
	return key
}

func (d *CommonDynamic) format(s string) (string, []interface{}) {
	i, index := 0, 1
	var params []interface{}
	for _, k := range d.keys {
		s, i = replace(s, k, d.holder(index), -1)
		if i > 0 {
			params = append(params, d.paramMap[k])
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

func IsTrue(i interface{}) bool {
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
