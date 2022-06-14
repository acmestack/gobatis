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

package parsing

import (
	"github.com/acmestack/gobatis/logging"
	"github.com/acmestack/gobatis/parsing/sqlparser"
	"github.com/acmestack/gobatis/reflection"
	"reflect"
	"strings"
	"time"
)

type GetFunc func(key string) string

type DynamicElement interface {
	Format(func(key string) string) string
}

type DynamicData struct {
	OriginData     string
	DynamicElemMap map[string]DynamicElement
}

func (m *DynamicData) Replace(params ...interface{}) string {
	objMap := reflection.ParseParams(params...)
	return m.ReplaceWithMap(objMap)
}

//需要外部确保param是一个struct
func (m *DynamicData) ReplaceWithMap(objParams map[string]interface{}) string {
	if len(m.DynamicElemMap) == 0 || len(objParams) == 0 {
		logging.Info("map is empty")
		//return m.OriginData
	}

	getFunc := func(s string) string {
		if o, ok := objParams[s]; ok {
			if str, ok := o.(string); ok {
				return str
			}

			//zero time convert to empty string (for <if> </if> element)
			if ti, ok := o.(time.Time); ok {
				if ti.IsZero() {
					return ""
				} else {
					return ti.String()
				}
			}

			var str string
			reflection.SafeSetValue(reflect.ValueOf(&str), o)
			return str
		}
		return ""
	}

	ret := m.OriginData
	for k, v := range m.DynamicElemMap {
		ret = strings.Replace(ret, k, v.Format(getFunc), -1)
	}
	return ret
}

func (m *DynamicData) ParseMetadata(driverName string, params ...interface{}) (*sqlparser.Metadata, error) {
	paramMap := reflection.ParseParams(params...)
	sqlStr := m.ReplaceWithMap(paramMap)
	return sqlparser.ParseWithParamMap(driverName, sqlStr, paramMap)
}
