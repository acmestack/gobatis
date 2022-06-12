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

package template

import (
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"github.com/xfali/gobatis/parsing/sqlparser"
	"strings"
	"text/template"
)

type CommonV2Dynamic struct {
	index    int
	keys     []string
	paramMap map[string]interface{}
	holder   sqlparser.Holder
}

func (d *CommonV2Dynamic) getFuncMap() template.FuncMap {
	ret := sprig.TxtFuncMap()
	ret[FuncNameSet] = d.UpdateSet
	ret[FuncNameWhere] = d.Where
	ret[FuncNameArg] = d.Param
	ret[FuncNameAdd] = commonAdd
	return ret
}

func (d *CommonV2Dynamic) UpdateSet(segments ...interface{}) string {
	buf := strings.Builder{}
	if len(segments) > 0 {
		buf.WriteString(" SET ")
	} else {
		return ""
	}
	for _, value := range segments {
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
	}

	return buf.String()
}

func (d *CommonV2Dynamic) Where(segments ...interface{}) string {
	buf := strings.Builder{}
	if len(segments) > 0 {
		buf.WriteString(" WHERE ")
	} else {
		return ""
	}
	for _, value := range segments {
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
	}

	return buf.String()
}

func (d *CommonV2Dynamic) getParam() []interface{} {
	return nil
}

func (d *CommonV2Dynamic) Param(p interface{}) string {
	d.index++
	key := getPlaceHolderKey(d.index)
	d.paramMap[key] = p
	d.keys = append(d.keys, key)
	return key
}

func (d *CommonV2Dynamic) format(s string) (string, []interface{}) {
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

func CreateV2DynamicHandler(h sqlparser.Holder) Dynamic {
	return &CommonV2Dynamic{
		index:    0,
		keys:     nil,
		paramMap: map[string]interface{}{},
		holder:   h,
	}
}
