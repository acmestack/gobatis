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

package xml

import (
	"strings"

	"github.com/acmestack/gobatis/logging"
	"github.com/acmestack/gobatis/parsing"
)

type Mapper struct {
	Namespace  string      `xml:"namespace,attr"`
	ResultMaps []ResultMap `xml:"resultMap"`
	Sql        []Sql       `xml:"sql"`

	Insert []Insert `xml:"insert"`
	Update []Update `xml:"update"`
	Select []Select `xml:"select"`
	Delete []Delete `xml:"delete"`
}

func (mapper *Mapper) Format() map[string]*parsing.DynamicData {
	ret := map[string]*parsing.DynamicData{}
	keyPre := strings.TrimSpace(mapper.Namespace)
	if keyPre != "" {
		keyPre = keyPre + "."
	}
	for _, v := range mapper.Insert {
		key := keyPre + v.Id
		if d, ok := ret[key]; ok {
			logging.Warn("Insert Sql id is duplicates, id: %s, before: %s, after %s\n", v.Id, d, v.Data)
		}
		d, err := ParseDynamic(strings.TrimSpace(v.Data), mapper.Sql)
		if err == nil {
			ret[key] = d
		}
	}
	for _, v := range mapper.Update {
		key := keyPre + v.Id
		if d, ok := ret[key]; ok {
			logging.Warn("Update Sql id is duplicates, id: %s, before: %s, after %s\n", v.Id, d, v.Data)
		}
		d, err := ParseDynamic(strings.TrimSpace(v.Data), mapper.Sql)
		if err == nil {
			ret[key] = d
		}
	}
	for _, v := range mapper.Select {
		key := keyPre + v.Id
		if d, ok := ret[key]; ok {
			logging.Warn("Select Sql id is duplicates, id: %s, before: %s, after %s\n", v.Id, d, v.Data)
		}
		d, err := ParseDynamic(strings.TrimSpace(v.Data), mapper.Sql)
		if err == nil {
			ret[key] = d
		}
	}
	for _, v := range mapper.Delete {
		key := keyPre + v.Id
		if d, ok := ret[v.Id]; ok {
			logging.Warn("Delete Sql id is duplicates, id: %s, before: %s, after %s\n", v.Id, d, v.Data)
		}
		d, err := ParseDynamic(strings.TrimSpace(v.Data), mapper.Sql)
		if err == nil {
			ret[key] = d
		}
	}
	return ret
}
