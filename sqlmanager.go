/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description:
 */

package gobatis

import (
	"github.com/xfali/gobatis/errors"
	"github.com/xfali/gobatis/logging"
	"github.com/xfali/gobatis/parsing"
	"github.com/xfali/gobatis/parsing/xml"
	"sync"
)

type SqlManager struct {
	sqlMap map[string]*parsing.DynamicData
	lock   sync.Mutex
}

var g_sql_mgr = SqlManager{sqlMap: map[string]*parsing.DynamicData{}}

func RegisterSql(sqlId string, sql string) error {
	g_sql_mgr.lock.Lock()
	defer g_sql_mgr.lock.Unlock()

	if _, ok := g_sql_mgr.sqlMap[sqlId]; ok {
		return errors.SQL_ID_DUPLICATES
	} else {
		dd := &parsing.DynamicData{OriginData: sql}
		g_sql_mgr.sqlMap[sqlId] = dd
	}
	return nil
}

func UnregisterSql(sqlId string) {
	g_sql_mgr.lock.Lock()
	defer g_sql_mgr.lock.Unlock()

	delete(g_sql_mgr.sqlMap, sqlId)
}

func RegisterMapperData(data []byte) error {
	g_sql_mgr.lock.Lock()
	defer g_sql_mgr.lock.Unlock()

	mapper, err := xml.Parse(data)
	if err != nil {
		logging.Warn("register mapper data failed: %s err: %v\n", string(data), err)
		return err
	}

	return formatMapper(mapper)
}

func RegisterMapperFile(file string) error {
	g_sql_mgr.lock.Lock()
	defer g_sql_mgr.lock.Unlock()

	mapper, err := xml.ParseFile(file)
	if err != nil {
		logging.Warn("register mapper file failed: %s err: %v\n", file, err)
		return err
	}

	return formatMapper(mapper)
}

func formatMapper(mapper *xml.Mapper) error {
	ret := mapper.Format()
	for k, v := range ret {
		if _, ok := g_sql_mgr.sqlMap[k]; ok {
			return errors.SQL_ID_DUPLICATES
		} else {
			g_sql_mgr.sqlMap[k] = v
		}
	}
	return nil
}

func FindSql(sqlId string) *parsing.DynamicData {
	g_sql_mgr.lock.Lock()
	defer g_sql_mgr.lock.Unlock()

	return g_sql_mgr.sqlMap[sqlId]
}
