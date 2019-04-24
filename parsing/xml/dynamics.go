/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package xml

import (
    "encoding/xml"
    "fmt"
    "github.com/xfali/gobatis/logging"
    "github.com/xfali/gobatis/reflection"
    "reflect"
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
    Test string `xml:"test,attr"`
    Date string `xml:",chardata"`
}

type Where struct {
    If []If `xml:"if"`
}

type Set struct {
    If []If `xml:"if"`
}

type DynamicData struct {
    OriginData string
    //If       []If    `xml:"if"`
    //Include  Include `xml:"include"`
    //Set      Set     `xml:"set"`
    //Where    Where   `xml:"where"`

    dynamicElemMap map[string]DynamicElement
}

type GetFunc func(key string) string

type DynamicElement interface {
    Format(func(key string) string) string
}

//传入方法必须是通过参数名获得参数值
func (de *Sql) Format(getFunc func(key string) string) string {
    //TODO: 增加变量支持
    return de.Sql
}

//传入方法必须是通过参数名获得参数值
func (de *If) Format(getFunc func(key string) string) string {
    andStrs := strings.Split(de.Test, " and ")
    if len(andStrs) != 0 {
        for _, v := range andStrs {
            ret := Compare(v, getFunc)
            if ret != true {
                return ""
            }
        }
        return strings.TrimSpace(de.Date)
    }

    orStrs := strings.Split(de.Test, " or ")
    ret := false
    if len(orStrs) != 0 {
        for _, v := range orStrs {
            ret = Compare(v, getFunc)
            if ret == true {
                return strings.TrimSpace(de.Date)
            }
        }
        if ret == false {
            return ""
        }
    }
    return ""
}

func Compare(src string, getFcun func(key string) string) bool {
    params := strings.Split(src, " ")
    if len(params) > 2 {
        value := getFcun(params[0])
        if value == "" {
            value = "nil"
        }
        switch params[1] {
        case "=":
            if value == params[2] {
                return true
            }
            break
        case "!=":
            if value != params[2] {
                return true
            }
            break
        }
    }
    return false
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
        ret.WriteString(" set ")
        for i := range de.If {
            ret.WriteString(de.If[i].Format(getFunc))
            ret.WriteString(" ")
        }
    }
    return ret.String()
}

//传入方法必须是通过参数名获得参数值
func (de *Where) Format(getFunc func(key string) string) string {
    ret := strings.Builder{}
    if len(de.If) > 0 {
        ret.WriteString(" where ")
        //FIXME: where 不会自动移除第一个元素的AND和OR，该特性后续添加
        for i := range de.If {
            ret.WriteString(de.If[i].Format(getFunc))
            ret.WriteString(" ")
        }
    }
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
    Parse(src string) DynamicElement
}

type IfProcessor string
type WhereProcessor string
type SetProcessor string
type IncludeProcessor string

var gProcessorMap = map[string]typeProcessor{
    "if":      IfProcessor("if"),
    "where":   WhereProcessor("where"),
    "set":     SetProcessor("set"),
    "include": IncludeProcessor("include"),
}

func (d IfProcessor) EndStr() string {
    return "</" + string(d) + ">"
}

func (d IfProcessor) Parse(src string) DynamicElement {
    v := If{}
    if xml.Unmarshal([]byte(src), &v) != nil {
        logging.Warn("parse if element failed")
    }
    return &v
}

func (d WhereProcessor) EndStr() string {
    return "</" + string(d) + ">"
}

func (d WhereProcessor) Parse(src string) DynamicElement {
    v := Where{}
    if xml.Unmarshal([]byte(src), &v) != nil {
        logging.Warn("parse if element failed")
    }
    return &v
}

func (d SetProcessor) EndStr() string {
    return "</" + string(d) + ">"
}

func (d SetProcessor) Parse(src string) DynamicElement {
    v := Set{}
    if xml.Unmarshal([]byte(src), &v) != nil {
        logging.Warn("parse if element failed")
    }
    return &v
}

func (d IncludeProcessor) EndStr() string {
    return "</" + string(d) + ">"
}

func (d IncludeProcessor) Parse(src string) DynamicElement {
    v := Include{}
    if xml.Unmarshal([]byte(src), &v) != nil {
        logging.Warn("parse if element failed")
    }
    return &v
}

func ParseDynamic(src string, sqls []Sql) (*DynamicData, error) {
    src = escape(src)

    start, end := -1, -1
    ret := &DynamicData{OriginData: src, dynamicElemMap: map[string]DynamicElement{}}
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
                    logging.Debug("save element : %s\n", saveStr)
                    de := typeProcessor.Parse(saveStr)
                    if include, ok := de.(*Include); ok {
                        findSql(include, sqls)
                    }
                    ret.dynamicElemMap[saveStr] = de
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

func (m *DynamicData) Replace(params ...interface{}) string {
    if len(params) == 0 {
        return m.OriginData
    }

    var getFunc func(string) string
    paramMap := map[string]string{}
    var paramStr string
    if len(params) == 1 {
        t := reflect.TypeOf(params[0])
        if t.Kind() == reflect.Ptr {
            t = t.Elem()
        }
        if t.Kind() == reflect.Struct {
            ti, err := reflection.GetTableInfo(params[0])
            if err != nil {
                logging.Info("%s", err.Error())
                return m.OriginData
            }
            objParams := ti.MapValue()
            getFunc = func(s string) string {
                if o, ok := objParams[s]; ok {
                    var str string
                    reflection.SetValue(reflect.ValueOf(str), o)
                    return str
                }
                return ""
            }
        } else {
            if reflection.IsSimpleType(params[0]) {
                reflection.SetValue(reflect.ValueOf(paramStr), params[0])
                paramMap["${0}"] = paramStr
                getFunc = func(s string) string {
                    return paramMap[s]
                }
            }
        }
    } else {
        objParams := map[string]interface{}{}
        for i, v := range params {
            if !reflection.IsSimpleType(v) {
                logging.Warn("Param error: expect simple type, but get other type")
                return m.OriginData
            }
            key := fmt.Sprintf("${%d}", i)
            objParams[key] = v
        }
        getFunc = func(s string) string {
            if o, ok := objParams[s]; ok {
                var str string
                reflection.SetValue(reflect.ValueOf(str), o)
                return str
            }
            return ""
        }
    }

    ret := m.OriginData
    for k, v := range m.dynamicElemMap {
        ret = strings.Replace(ret, k, v.Format(getFunc), -1)
    }
    return ret
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
