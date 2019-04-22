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
    "strings"
)

type Mapper struct {
    Namespace string `xml:"namespace,attr"`
    Sql       []Sql  `xml:"sql"`

    Insert []Insert `xml:"insert"`
    Update []Update `xml:"update"`
    Select []Select `xml:"select"`
    Delete []Delete `xml:"delete"`
}

func (m *Mapper) Format() map[string]string {
    ret := map[string]string{}
    for _, v := range m.Insert {
        if d, ok := ret[v.Id]; ok {
            logging.Warn("Insert Sql id is duplicates, id: %s, before: %s, after %s\n", v.Id, d, v.Data)
        }
        ret[v.Id] = strings.TrimSpace(v.Data)
    }
    for _, v := range m.Update {
        if d, ok := ret[v.Id]; ok {
            logging.Warn("Update Sql id is duplicates, id: %s, before: %s, after %s\n", v.Id, d, v.Data)
        }
        ret[v.Id] = strings.TrimSpace(v.Data)
    }
    for _, v := range m.Select {
        if d, ok := ret[v.Id]; ok {
            logging.Warn("Select Sql id is duplicates, id: %s, before: %s, after %s\n", v.Id, d, v.Data)
        }
        ret[v.Id] = strings.TrimSpace(v.Data)
    }
    for _, v := range m.Delete {
        if d, ok := ret[v.Id]; ok {
            logging.Warn("Delete Sql id is duplicates, id: %s, before: %s, after %s\n", v.Id, d, v.Data)
        }
        ret[v.Id] = strings.TrimSpace(v.Data)
    }
    return ret
}
