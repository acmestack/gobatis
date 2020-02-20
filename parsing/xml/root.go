/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description:
 */

package xml

import (
	"github.com/xfali/gobatis/logging"
	"github.com/xfali/gobatis/parsing"
	"strings"
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

func (m *Mapper) Format() map[string]*parsing.DynamicData {
	ret := map[string]*parsing.DynamicData{}
	keyPre := strings.TrimSpace(m.Namespace)
	if keyPre != "" {
		keyPre = keyPre + "."
	}
	for _, v := range m.Insert {
		key := keyPre + v.Id
		if d, ok := ret[key]; ok {
			logging.Warn("Insert Sql id is duplicates, id: %s, before: %s, after %s\n", v.Id, d, v.Data)
		}
		d, err := ParseDynamic(strings.TrimSpace(v.Data), m.Sql)
		if err == nil {
			ret[key] = d
		}
	}
	for _, v := range m.Update {
		key := keyPre + v.Id
		if d, ok := ret[key]; ok {
			logging.Warn("Update Sql id is duplicates, id: %s, before: %s, after %s\n", v.Id, d, v.Data)
		}
		d, err := ParseDynamic(strings.TrimSpace(v.Data), m.Sql)
		if err == nil {
			ret[key] = d
		}
	}
	for _, v := range m.Select {
		key := keyPre + v.Id
		if d, ok := ret[key]; ok {
			logging.Warn("Select Sql id is duplicates, id: %s, before: %s, after %s\n", v.Id, d, v.Data)
		}
		d, err := ParseDynamic(strings.TrimSpace(v.Data), m.Sql)
		if err == nil {
			ret[key] = d
		}
	}
	for _, v := range m.Delete {
		key := keyPre + v.Id
		if d, ok := ret[v.Id]; ok {
			logging.Warn("Delete Sql id is duplicates, id: %s, before: %s, after %s\n", v.Id, d, v.Data)
		}
		d, err := ParseDynamic(strings.TrimSpace(v.Data), m.Sql)
		if err == nil {
			ret[key] = d
		}
	}
	return ret
}
