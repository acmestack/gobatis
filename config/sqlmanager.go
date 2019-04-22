/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package config

import "sync"

type SqlManager struct {
    sqlMap map[string]string
    lock   sync.Mutex
}

var g_sql_mgr = SqlManager{sqlMap: map[string]string{}}

func RegisterSql(sqlId string, sql string) {
    g_sql_mgr.lock.Lock()
    defer g_sql_mgr.lock.Unlock()

    g_sql_mgr.sqlMap[sqlId] = sql
}

func FindSql(sqlId string) string {
    g_sql_mgr.lock.Lock()
    defer g_sql_mgr.lock.Unlock()

    return g_sql_mgr.sqlMap[sqlId]
}
