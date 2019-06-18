/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description:
 */

package builder

import (
	"fmt"
	"strconv"
	"strings"
)

type SQLFragment struct {
	builder strings.Builder
	child   *SQLFragment
	parent  *SQLFragment
}

func Select(columns ...string) *SQLFragment {
	fragment := &SQLFragment{}

	fragment.builder.WriteString("SELECT ")
	if len(columns) == 0 {
		panic("param error")
	}
	fragment.builder.WriteString(columns[0])
	for i := 1; i < len(columns); i++ {
		fragment.builder.WriteString(", ")
		fragment.builder.WriteString(columns[i])
	}
	fragment.builder.WriteString(" ")
	return fragment
}

func DeleteFrom(table string) *SQLFragment {
	fragment := &SQLFragment{}

	fragment.builder.WriteString(fmt.Sprintf("DELETE FROM %s ", table))

	return fragment
}

func InsertInto(table string) *SQLFragment {
	fragment := &SQLFragment{}

	fragment.builder.WriteString(fmt.Sprintf("INSERT INTO %s ", table))

	return fragment
}

func Update(table string) *SQLFragment {
	fragment := &SQLFragment{}

	fragment.builder.WriteString(fmt.Sprintf("UPDATE %s ", table))

	return fragment
}

func (f *SQLFragment) Select(columns ...string) *SQLFragment {
	str := f.builder.String()
	str = str[7 : len(str)-1]

	if len(columns) == 0 {
		panic("param error")
	}
	f.builder.Reset()
	f.builder.WriteString("SELECT ")
	f.builder.WriteString(str)
	f.builder.WriteString(", ")
	f.builder.WriteString(columns[0])
	for i := 1; i < len(columns); i++ {
		f.builder.WriteString(", ")
		f.builder.WriteString(columns[i])
	}
	f.builder.WriteString(" ")
	return f
}

func (f *SQLFragment) From(table string) *SQLFragment {
	fragment := &SQLFragment{}
	fragment.initParent(f)

	if !check(f.builder.String(), "FROM", 0) {
		fragment.builder.WriteString("FROM ")
	} else {
		//当前元素为values，则清空当前元素值，移到子元素
		curStr := f.builder.String()
		curStr = curStr[5 : len(curStr)-1]
		f.builder.Reset()

		fragment.builder.WriteString("FROM ")
		fragment.builder.WriteString(curStr)
		fragment.builder.WriteString(", ")
	}
	fragment.builder.WriteString(table)
	fragment.builder.WriteString(" ")
	return fragment
}

func (f *SQLFragment) Join(join string) *SQLFragment {
	fragment := &SQLFragment{}
	fragment.initParent(f)

	fragment.builder.WriteString("JOIN ")
	fragment.builder.WriteString(join)
	fragment.builder.WriteString(" ")
	return fragment
}

func (f *SQLFragment) InnerJoin(join string) *SQLFragment {
	fragment := &SQLFragment{}
	fragment.initParent(f)

	fragment.builder.WriteString("INNER JOIN ")
	fragment.builder.WriteString(join)
	fragment.builder.WriteString(" ")
	return fragment
}

func (f *SQLFragment) LeftJoin(join string) *SQLFragment {
	fragment := &SQLFragment{}
	fragment.initParent(f)

	fragment.builder.WriteString("LEFT JOIN ")
	fragment.builder.WriteString(join)
	fragment.builder.WriteString(" ")
	return fragment
}

func (f *SQLFragment) RightJoin(join string) *SQLFragment {
	fragment := &SQLFragment{}
	fragment.initParent(f)

	fragment.builder.WriteString("RIGHT JOIN ")
	fragment.builder.WriteString(join)
	fragment.builder.WriteString(" ")
	return fragment
}

func (f *SQLFragment) Where(condition string) *SQLFragment {
	fragment := &SQLFragment{}
	fragment.initParent(f)
	str := f.builder.String()
	if str == "AND" || str == "OR" {
		fragment.builder.WriteString(" ")
		fragment.builder.WriteString(condition)
		fragment.builder.WriteString(" ")
	} else {
		fragment.builder.WriteString("WHERE ")
		fragment.builder.WriteString(condition)
		fragment.builder.WriteString(" ")
	}
	return fragment
}

func (f *SQLFragment) GroupBy(columns ...string) *SQLFragment {
	fragment := &SQLFragment{}
	fragment.initParent(f)

	fragment.builder.WriteString("GROUP BY ")
	if len(columns) == 0 {
		panic("param error")
	}
	fragment.builder.WriteString(columns[0])
	for i := 1; i < len(columns); i++ {
		fragment.builder.WriteString(", ")
		fragment.builder.WriteString(columns[i])
	}
	fragment.builder.WriteString(" ")
	return fragment
}

func (f *SQLFragment) Having(condition string) *SQLFragment {
	fragment := &SQLFragment{}
	fragment.initParent(f)

	str := f.builder.String()
	if str == "AND" || str == "OR" {
		fragment.builder.WriteString(" ")
		fragment.builder.WriteString(condition)
		fragment.builder.WriteString(" ")
	} else {
		fragment.builder.WriteString("HAVING ")
		fragment.builder.WriteString(condition)
		fragment.builder.WriteString(" ")
	}
	return fragment
}

func (f *SQLFragment) And() *SQLFragment {
	fragment := &SQLFragment{}
	fragment.initParent(f)

	str := f.builder.String()
	if check(str, "HAVING", 0) || check(str, "WHERE", 0) {
		fragment.builder.WriteString("AND")
	}

	return fragment
}

func (f *SQLFragment) Or() *SQLFragment {
	fragment := &SQLFragment{}
	fragment.initParent(f)

	str := f.builder.String()
	if check(str, "HAVING", 0) || check(str, "WHERE", 0) {
		fragment.builder.WriteString("OR")
	}

	return fragment
}

func (f *SQLFragment) OrderBy(columns ...string) *SQLFragment {
	fragment := &SQLFragment{}
	fragment.initParent(f)

	if len(columns) == 0 {
		panic("param error")
	}
	fragment.builder.WriteString("ORDER BY ")
	fragment.builder.WriteString(columns[0])
	for i := 1; i < len(columns); i++ {
		fragment.builder.WriteString(", ")
		fragment.builder.WriteString(columns[i])
	}
	fragment.builder.WriteString(" ")

	return fragment
}

func (f *SQLFragment) Desc() *SQLFragment {
	fragment := &SQLFragment{}
	fragment.initParent(f)

	fragment.builder.WriteString("DESC ")
	return fragment
}

func (f *SQLFragment) Asc() *SQLFragment {
	fragment := &SQLFragment{}
	fragment.initParent(f)

	fragment.builder.WriteString("ASC ")
	return fragment
}

func (f *SQLFragment) Offset(offset int64) *SQLFragment {
	fragment := &SQLFragment{}
	fragment.initParent(f)

	fragment.builder.WriteString("OFFSET ")
	fragment.builder.WriteString(strconv.FormatInt(offset, 10))
	fragment.builder.WriteString(" ")

	return fragment
}

func (f *SQLFragment) Limit(limit int64) *SQLFragment {
	fragment := &SQLFragment{}
	fragment.initParent(f)

	fragment.builder.WriteString("LIMIT ")
	fragment.builder.WriteString(strconv.FormatInt(limit, 10))
	fragment.builder.WriteString(" ")

	return fragment
}

func (f *SQLFragment) Set(column string, value string) *SQLFragment {
	fragment := &SQLFragment{}
	fragment.initParent(f)

	if !check(f.builder.String(), "SET", 0) {
		fragment.builder.WriteString("SET ")
	} else {
		fragment.builder.WriteString(", ")
	}

	fragment.builder.WriteString(column)
	fragment.builder.WriteString(" = ")
	fragment.builder.WriteString(value)
	fragment.builder.WriteString(" ")

	return fragment
}

func (f *SQLFragment) IntoColumns(columns ...string) *SQLFragment {
	fragment := &SQLFragment{}
	fragment.initParent(f)

	if len(columns) == 0 {
		panic("param error")
	}
	if !check(f.builder.String(), "(", 0) {
		fragment.builder.WriteString("(")
	} else {
		//当前元素为values，则清空当前元素值，移到子元素
		curStr := f.builder.String()
		curStr = curStr[1 : len(curStr)-2]
		f.builder.Reset()

		fragment.builder.WriteString("(")
		fragment.builder.WriteString(curStr)
		fragment.builder.WriteString(", ")
	}
	fragment.builder.WriteString(columns[0])
	for i := 1; i < len(columns); i++ {
		fragment.builder.WriteString(", ")
		fragment.builder.WriteString(columns[i])
	}
	fragment.builder.WriteString(") ")
	return fragment
}

func (f *SQLFragment) IntoValues(values ...string) *SQLFragment {
	fragment := &SQLFragment{}
	fragment.initParent(f)

	if len(values) == 0 {
		panic("param error")
	}
	if !check(f.builder.String(), "VALUES", 0) {
		fragment.builder.WriteString("VALUES(")
	} else {
		//当前元素为values，则清空当前元素值，移到子元素
		curStr := f.builder.String()
		curStr = curStr[7 : len(curStr)-2]
		f.builder.Reset()

		fragment.builder.WriteString("VALUES(")
		fragment.builder.WriteString(curStr)
		fragment.builder.WriteString(", ")
	}
	fragment.builder.WriteString(values[0])
	for i := 1; i < len(values); i++ {
		fragment.builder.WriteString(", ")
		fragment.builder.WriteString(values[i])
	}
	fragment.builder.WriteString(") ")
	return fragment
}

func (f *SQLFragment) String() string {
	buf := strings.Builder{}
	root := f
	for root.parent != nil {
		root = root.parent
	}
	buf.WriteString(root.builder.String())
	child := root.child
	for child != nil {
		buf.WriteString(child.builder.String())
		child = child.child
	}
	return buf.String()
}

func (f *SQLFragment) Hook(hookFunc func(*SQLFragment) *SQLFragment) *SQLFragment {
	return hookFunc(f)
}

func (f *SQLFragment) initParent(parent *SQLFragment) *SQLFragment {
	f.parent = parent
	parent.child = f
	return f
}

func check(src, dest string, srcOffset int) bool {
	length := len(dest)
	if srcOffset+length >= len(src) {
		//LOG:
		return false
	}
	return src[srcOffset:srcOffset+length] == dest
}
