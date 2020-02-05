/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description:
 */

package gobatis

import (
	"github.com/xfali/gobatis/parsing"
	"github.com/xfali/gobatis/parsing/sqlparser"
	"github.com/xfali/gobatis/parsing/template"
	"github.com/xfali/gobatis/parsing/xml"
)

type sqlManager struct {
	dynamicSqlMgr  *xml.Manager
	templateSqlMgr *template.Manager
}

var g_sql_mgr = sqlManager{
	dynamicSqlMgr:  xml.NewManager(),
	templateSqlMgr: template.NewManager(),
}

func RegisterSql(sqlId string, sql string) error {
	return g_sql_mgr.dynamicSqlMgr.RegisterSql(sqlId, sql)
}

func UnregisterSql(sqlId string) {
	g_sql_mgr.dynamicSqlMgr.UnregisterSql(sqlId)
}

func RegisterMapperData(data []byte) error {
	return g_sql_mgr.dynamicSqlMgr.RegisterData(data)
}

func RegisterMapperFile(file string) error {
	return g_sql_mgr.dynamicSqlMgr.RegisterFile(file)
}

func FindDynamicSqlParser(sqlId string) (sqlparser.SqlParser, bool) {
	return g_sql_mgr.dynamicSqlMgr.FindSqlParser(sqlId)
}

func RegisterTemplateData(data []byte) error {
	return g_sql_mgr.templateSqlMgr.RegisterData(data)
}

func RegisterTemplateFile(file string) error {
	return g_sql_mgr.templateSqlMgr.RegisterFile(file)
}

func FindTemplateSqlParser(sqlId string) (sqlparser.SqlParser, bool) {
	return g_sql_mgr.templateSqlMgr.FindSqlParser(sqlId)
}

func FindSqlParser(sqlId string) sqlparser.SqlParser {
	ret, ok := FindDynamicSqlParser(sqlId)
	if !ok {
		ret, ok = FindTemplateSqlParser(sqlId)
	}
	//FIXME: 当没有查找到sqlId对应的sql语句，则尝试使用sqlId直接操作数据库
	//该设计可能需要设计一个更合理的方式
	if !ok {
		return &parsing.DynamicData{OriginData: sqlId}
	}
	return ret
}
