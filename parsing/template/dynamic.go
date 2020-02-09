// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description: 

package template

import (
	"fmt"
	"strings"
	"text/template"
)

func dummyUpdateSet(b bool, column string, value interface{}, origin string) string {
	return origin
}

func dummyWhere(b bool, cond, column string, value interface{}, origin string) string {
	return origin
}

func commonAdd(a, b int) int {
	return a + b
}

func mysqlUpdateSet(b bool, column string, value interface{}, origin string) string {
	if !b {
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
	buf.WriteString("`")
	buf.WriteString(column)
	buf.WriteString("` = ")
	if s, ok := value.(string); ok {
		buf.WriteString(`'`)
		buf.WriteString(s)
		buf.WriteString(`'`)
	} else {
		buf.WriteString(fmt.Sprint(value))
	}
	return buf.String()
}

func postgresUpdateSet(b bool, column string, value interface{}, origin string) string {
	if !b {
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
	buf.WriteString(`"`)
	buf.WriteString(column)
	buf.WriteString(`"`)
	buf.WriteString(" = ")
	if s, ok := value.(string); ok {
		buf.WriteString(`'`)
		buf.WriteString(s)
		buf.WriteString(`'`)
	} else {
		buf.WriteString(fmt.Sprint(value))
	}
	return buf.String()
}

func mysqlWhere(b bool, cond, column string, value interface{}, origin string) string {
	if !b {
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

	buf.WriteString("`")
	buf.WriteString(column)
	buf.WriteString("` = ")
	if s, ok := value.(string); ok {
		buf.WriteString(`'`)
		buf.WriteString(s)
		buf.WriteString(`'`)
	} else {
		buf.WriteString(fmt.Sprint(value))
	}
	return buf.String()
}

func postgresWhere(b bool, cond, column string, value interface{}, origin string) string {
	if !b {
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

	buf.WriteString(`"`)
	buf.WriteString(column)
	buf.WriteString(`"`)
	buf.WriteString(" = ")
	if s, ok := value.(string); ok {
		buf.WriteString(`'`)
		buf.WriteString(s)
		buf.WriteString(`'`)
	} else {
		buf.WriteString(fmt.Sprint(value))
	}
	return buf.String()
}

var mysqlFuncMap = template.FuncMap{
	"set":   mysqlUpdateSet,
	"where": mysqlWhere,
	"add":   commonAdd,
}

var postgresFuncMap = template.FuncMap{
	"set":   postgresUpdateSet,
	"where": postgresWhere,
	"add":   commonAdd,
}

var dummyFuncMap = template.FuncMap{
	"set":   dummyUpdateSet,
	"where": dummyWhere,
	"add":   commonAdd,
}

var funcMap = map[string]template.FuncMap{
	"mysql":    mysqlFuncMap,
	"postgres": postgresFuncMap,
}

func selectFuncMap(driverName string) template.FuncMap {
	if v, ok := funcMap[driverName]; ok {
		return v
	}
	return dummyFuncMap
}
