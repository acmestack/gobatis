/*
 * Copyright (c) 2022, OpeningO
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

package xml

import (
	"encoding/xml"
	"fmt"
	"github.com/xfali/gobatis/logging"
	"github.com/xfali/gobatis/parsing"
	"strconv"
	"strings"
	"unicode"
)

type Foreach struct {
	Item       string `xml:"item,attr"`
	Collection string `xml:"collection,attr"`
	Separator  string `xml:"separator,attr"`
	Index      string `xml:"index,attr"`
	Open       string `xml:"open,attr"`
	Close      string `xml:"close,attr"`
	Data       string `xml:",chardata"`
}

type Sql struct {
	Id  string `xml:"id,attr"`
	Sql string `xml:",chardata"`
}

type Property struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type Include struct {
	Refid      string     `xml:"refid,attr"`
	Properties []Property `xml:"property"`
	Sql        Sql        `xml:-`
}

type If struct {
	Foreach Foreach `xml:"foreach"`
	Test    string  `xml:"test,attr"`
	Data    string  `xml:",chardata"`
}

type Where struct {
	Choose Choose `xml:"choose"`
	If     []If   `xml:"if"`
}

type Set struct {
	If []If `xml:"if"`
}

type When struct {
	If
}

type Otherwise struct {
	Data string `xml:",chardata"`
}

type Choose struct {
	When      []When    `xml:"when"`
	Otherwise Otherwise `xml:"otherwise"`
}

//传入方法必须是通过参数名获得参数值
func (de *Sql) Format(getFunc func(key string) string) string {
	//TODO: 增加变量支持
	return de.Sql
}

//传入方法必须是通过参数名获得参数值
func (de *If) Format(getFunc func(key string) string) string {
	andStrs := strings.Split(de.Test, " and ")
	orStrs := strings.Split(de.Test, " or ")

	if len(andStrs) > 1 && len(orStrs) > 1 {
		logging.Warn(`error format. <if> element cannot both include " and " and " or "`)
		return ""
	}

	data := ""
	if de.Foreach.Data != "" {
		data = strings.TrimSpace(de.Foreach.Format(getFunc))
	}
	data += strings.TrimSpace(de.Data)
	if len(andStrs) != 0 && len(orStrs) < 2 {
		for _, v := range andStrs {
			ret := Compare(v, getFunc)
			if ret != true {
				return ""
			}
		}
		return data
	}

	ret := false
	if len(orStrs) != 0 {
		for _, v := range orStrs {
			ret = Compare(v, getFunc)
			if ret == true {
				return data
			}
		}
		if ret == false {
			return ""
		}
	}
	return ""
}

//test的参数必须是使用{}包裹起来，并且比较符号需要空格分隔，如<if test="{1} != nil"> 或者 <if test="{x.name} != nil">
func Compare(src string, getFunc func(key string) string) bool {
	params := strings.Split(src, " ")
	if len(params) > 2 {
		v1 := getValueFromFunc(params[0], getFunc)
		v2 := getValueFromFunc(params[2], getFunc)
		if v1 == "" && v2 == "" {
			return false
		}
		switch params[1] {
		case "==":
			if v1 == v2 {
				return true
			}
			break
		case "!=":
			if v1 != v2 {
				return true
			}
			break
		}
	}
	return false
}

func getKey(src string) string {
	if src == "" {
		return ""
	}
	if src[:1] == "{" {
		index := strings.Index(src, "}")
		if index == -1 {
			return src
		}
		return src[1:index]
	}
	return src
}

func getValueFromFunc(src string, getFunc func(key string) string) string {
	if src == "" {
		return ""
	}
	if src[:1] == "{" {
		index := strings.Index(src, "}")
		if index == -1 {
			return src
		}
		ret := getFunc(src[1:index])
		if ret == "" {
			return "nil"
		} else {
			return ret
		}
	}
	return src
}

//传入方法必须是通过参数名获得参数值
func (de *Include) Format(getFunc func(key string) string) string {
	//TODO: sql的参数替换特性未实现
	return de.Sql.Sql
}

//传入方法必须是通过参数名获得参数值
func (de *Set) Format(getFunc func(key string) string) string {
	ret := strings.Builder{}
	if len(de.If) > 0 {
		//ret.WriteString(" set ")
		add := false
		for i := range de.If {
			ifStr := de.If[i].Format(getFunc)
			if ifStr != "" {
				if add {
					ret.WriteString(",")
				}

				if ifStr[len(ifStr)-1:] == "," {
					ret.WriteString(ifStr[:len(ifStr)-1])
				} else {
					ret.WriteString(ifStr)
				}
				add = true
			}
		}
	}
	retStr := ret.String()
	if retStr != "" {
		retStr = " set " + retStr
	}
	return retStr
}

//传入方法必须是通过参数名获得参数值
func (de *Where) Format(getFunc func(key string) string) string {
	ret := strings.Builder{}
	if len(de.If) > 0 {
		set := false
		//ret.WriteString(" where ")
		for i := range de.If {
			ifStr := de.If[i].Format(getFunc)
			if ifStr != "" {
				if !set {
					ifStr = removeAndOr(ifStr)
					set = true
				}
				ret.WriteString(ifStr)
				ret.WriteString(" ")
			}
		}
	} else {
		chooseStr := de.Choose.Format(getFunc)
		chooseStr = removeAndOr(chooseStr)
		if chooseStr != "" {
			ret.WriteString(chooseStr)
		}
	}
	retStr := ret.String()
	if retStr != "" {
		retStr = " where " + retStr
	}
	return retStr
}

func removeAndOr(src string) string {
	if len(src) > 3 && strings.ToLower(src[:3]) == "or " {
		src = strings.TrimSpace(src[3:])
	} else if len(src) > 4 && strings.ToLower(src[:4]) == "and " {
		src = strings.TrimSpace(src[4:])
	}
	return src
}

func (de *Choose) Format(getFunc func(key string) string) string {
	ret := strings.Builder{}
	if len(de.When) > 0 {
		for i := range de.When {
			ifStr := de.When[i].Format(getFunc)
			if ifStr != "" {
				ret.WriteString(ifStr)
				ret.WriteString(" ")
				break
			}
		}
	}
	retStr := ret.String()
	if retStr == "" {
		retStr = strings.TrimSpace(de.Otherwise.Data)
	}
	return retStr
}

func (de *Foreach) Format(getFunc func(key string) string) string {
	if de.Collection == "" {
		return ""
	}

	collectionKey := getKey(de.Collection)
	value := getValueFromFunc(de.Collection, getFunc)
	//if value == "" || value == "nil" {
	//	return ""
	//}
	//
	//elems := reflection.ParseSliceParamString(value)
	//length := len(elems)
	length, _ := strconv.Atoi(value)
	ret := strings.Builder{}
	ret.WriteString(de.Open)
	indexStr := fmt.Sprintf("#{%s}", de.Index)

	originData := strings.TrimSpace(de.Data)
	realStr := ""
	i := 0
	for i < length {
		realStr = strings.Replace(originData, de.Item, fmt.Sprintf("%s[%d]", collectionKey, i), -1)
		realStr = strings.Replace(realStr, indexStr, strconv.Itoa(i), -1)
		ret.WriteString(realStr)
		if i < length-1 {
			ret.WriteString(de.Separator)
		}
		i++
	}
	ret.WriteString(de.Close)

	return ret.String()
}

func escape(src string) string {
	src = strings.Replace(src, "&lt;", "<", -1)
	src = strings.Replace(src, "&gt;", ">", -1)
	src = strings.Replace(src, "&amp;", "&", -1)
	src = strings.Replace(src, "&quot;", `"`, -1)
	src = strings.Replace(src, "&apos;", `'`, -1)
	return src
}

type typeProcessor interface {
	EndStr() string
	Parse(src string) parsing.DynamicElement
}

type IfProcessor string
type WhereProcessor string
type SetProcessor string
type IncludeProcessor string
type ChooseProcessor string
type ForeachProcessor string

var gProcessorMap = map[string]typeProcessor{
	"if":      IfProcessor("if"),
	"where":   WhereProcessor("where"),
	"set":     SetProcessor("set"),
	"include": IncludeProcessor("include"),
	"choose":  ChooseProcessor("choose"),
	"foreach": ForeachProcessor("foreach"),
}

func (d IfProcessor) EndStr() string {
	return "</" + string(d) + ">"
}

func (d IfProcessor) Parse(src string) parsing.DynamicElement {
	v := If{}
	if xml.Unmarshal([]byte(src), &v) != nil {
		logging.Warn("parse if element failed")
	}
	return &v
}

func (d WhereProcessor) EndStr() string {
	return "</" + string(d) + ">"
}

func (d WhereProcessor) Parse(src string) parsing.DynamicElement {
	v := Where{}
	if xml.Unmarshal([]byte(src), &v) != nil {
		logging.Warn("parse if element failed")
	}
	return &v
}

func (d SetProcessor) EndStr() string {
	return "</" + string(d) + ">"
}

func (d SetProcessor) Parse(src string) parsing.DynamicElement {
	v := Set{}
	if xml.Unmarshal([]byte(src), &v) != nil {
		logging.Warn("parse if element failed")
	}
	return &v
}

func (d IncludeProcessor) EndStr() string {
	return "</" + string(d) + ">"
}

func (d IncludeProcessor) Parse(src string) parsing.DynamicElement {
	v := Include{}
	if xml.Unmarshal([]byte(src), &v) != nil {
		logging.Warn("parse if element failed")
	}
	return &v
}

func (d ChooseProcessor) EndStr() string {
	return "</" + string(d) + ">"
}

func (d ChooseProcessor) Parse(src string) parsing.DynamicElement {
	v := Choose{}
	if xml.Unmarshal([]byte(src), &v) != nil {
		logging.Warn("parse if element failed")
	}
	return &v
}

func (d ForeachProcessor) EndStr() string {
	return "</" + string(d) + ">"
}

func (d ForeachProcessor) Parse(src string) parsing.DynamicElement {
	v := Foreach{}
	if xml.Unmarshal([]byte(src), &v) != nil {
		logging.Warn("parse if element failed")
	}
	return &v
}

func ParseDynamic(src string, sqls []Sql) (*parsing.DynamicData, error) {
	src = escape(src)

	start, end := -1, -1
	ret := &parsing.DynamicData{OriginData: src, DynamicElemMap: map[string]parsing.DynamicElement{}}
	strData := []rune(src)
	for i := 0; i < len(strData); {
		r := strData[i]
		if r == '<' {
			start = i
		}

		if r == '>' {
			end = i
		}

		if start < end && start != -1 && end != -1 {
			subStr := src[start+1 : end]
			subStr = strings.TrimLeftFunc(subStr, unicode.IsSpace)
			index := strings.Index(subStr, " ")
			if index != -1 {
				subStr = subStr[:index]
			}
			logging.Debug("Found element : %s\n", subStr)
			if typeProcessor, ok := gProcessorMap[subStr]; ok {
				subStr = src[start:]
				endStr := typeProcessor.EndStr()
				index = strings.Index(subStr, endStr)
				if index == -1 {
					start, end = -1, -1
					i++
					start, end = -1, -1
					continue
				} else {
					saveStr := subStr[:index+len(endStr)]
					//logging.Debug("save element : %s\n", saveStr)
					de := typeProcessor.Parse(saveStr)
					if include, ok := de.(*Include); ok {
						findSql(include, sqls)
					}
					ret.DynamicElemMap[saveStr] = de
					i = start + index + 1 + len(endStr)
					start, end = -1, -1
					continue
				}
			}
			start, end = -1, -1
		}
		i++
	}
	return ret, nil
}

func findSql(include *Include, sqls []Sql) {
	if sqls != nil {
		for i := range sqls {
			if include.Refid == sqls[i].Id {
				include.Sql = sqls[i]
				return
			}
		}
	}
}
