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
	"github.com/xfali/gobatis/parsing/sqlparser"
	"github.com/xfali/gobatis/parsing/template"
	"github.com/xfali/gobatis/parsing/xml"
	"sync"
)

type dynamicSqlManager struct {
	sqlMap map[string]*parsing.DynamicData
	lock   sync.Mutex
}

type sqlManager struct {
	dynamicSqlMgr  *dynamicSqlManager
	templateSqlMgr *template.Manager
}

var g_sql_mgr = sqlManager{
	dynamicSqlMgr:  &dynamicSqlManager{sqlMap: map[string]*parsing.DynamicData{}},
	templateSqlMgr: template.NewManager(),
}

func RegisterSql(sqlId string, sql string) error {
	g_sql_mgr.dynamicSqlMgr.lock.Lock()
	defer g_sql_mgr.dynamicSqlMgr.lock.Unlock()

	if _, ok := g_sql_mgr.dynamicSqlMgr.sqlMap[sqlId]; ok {
		return errors.SQL_ID_DUPLICATES
	} else {
		dd := &parsing.DynamicData{OriginData: sql}
		g_sql_mgr.dynamicSqlMgr.sqlMap[sqlId] = dd
	}
	return nil
}

func UnregisterSql(sqlId string) {
	g_sql_mgr.dynamicSqlMgr.lock.Lock()
	defer g_sql_mgr.dynamicSqlMgr.lock.Unlock()

	delete(g_sql_mgr.dynamicSqlMgr.sqlMap, sqlId)
}

func RegisterMapperData(data []byte) error {
	g_sql_mgr.dynamicSqlMgr.lock.Lock()
	defer g_sql_mgr.dynamicSqlMgr.lock.Unlock()

	mapper, err := xml.Parse(data)
	if err != nil {
		logging.Warn("register mapper data failed: %s err: %v\n", string(data), err)
		return err
	}

	return formatMapper(mapper)
}

func RegisterMapperFile(file string) error {
	g_sql_mgr.dynamicSqlMgr.lock.Lock()
	defer g_sql_mgr.dynamicSqlMgr.lock.Unlock()

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
		if _, ok := g_sql_mgr.dynamicSqlMgr.sqlMap[k]; ok {
			return errors.SQL_ID_DUPLICATES
		} else {
			g_sql_mgr.dynamicSqlMgr.sqlMap[k] = v
		}
	}
	return nil
}

func FindDynamicSql(sqlId string) (sqlparser.SqlParser, bool) {
	g_sql_mgr.dynamicSqlMgr.lock.Lock()
	defer g_sql_mgr.dynamicSqlMgr.lock.Unlock()

	v, ok := g_sql_mgr.dynamicSqlMgr.sqlMap[sqlId]
	return v, ok
}

func RegisterTemplateData(data []byte) error {
	return g_sql_mgr.templateSqlMgr.RegisterData(data)
}

func RegisterTemplateFile(file string) error {
	return g_sql_mgr.templateSqlMgr.RegisterFile(file)
}

func FindTemplateSql(sqlId string) (sqlparser.SqlParser, bool) {
	return g_sql_mgr.templateSqlMgr.FindSql(sqlId)
}
