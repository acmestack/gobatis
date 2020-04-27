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
	"os"
	"path/filepath"
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

type ParserFactory func(sql string) (sqlparser.SqlParser, error)

func DynamicParserFactory(sql string) (sqlparser.SqlParser, error) {
	return &parsing.DynamicData{OriginData: sql}, nil
}

func TemplateParserFactory(sql string) (sqlparser.SqlParser, error) {
	return template.CreateParser([]byte(sql))
}

func ScanMapperFile(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			filename := filepath.Base(path)
			length := len(filename)
			if length > 4 {
				if filename[length-4:] == ".xml" {
					err := RegisterMapperFile(path)
					if err != nil {
						return err
					}
				}
				if filename[length-4:] == ".tpl" {
					err := RegisterTemplateFile(path)
					if err != nil {
						return err
					}
				}
			}
		}
		return nil
	})
}