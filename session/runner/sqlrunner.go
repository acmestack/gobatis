/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package runner

import (
    "github.com/xfali/gobatis/factory"
    "github.com/xfali/gobatis/handler"
    "github.com/xfali/gobatis/session"
    "reflect"
)

type SqlManager struct {
    sqlmap  map[string]string
    factory factory.Factory
}

func NewSqlManager(factory factory.Factory) *SqlManager {
    return &SqlManager{sqlmap: map[string]string{}, factory: factory}
}

func (this *SqlManager) RegisterSql(sqlId string, sql string) *SqlManager {
    this.sqlmap[sqlId] = sql
    return this
}

type SqlRunner struct {
    resultHandler handler.ResultHandler
    session       session.Session
    sql           string
    params        []interface{}
}

func (this *SqlManager) newSqlRunner(sqlId string) *SqlRunner {
    if sql, ok := this.sqlmap[sqlId]; ok {
        return &SqlRunner{session: this.factory.CreateSession(), sql: sql}
    }
    return nil
}

func (this *SqlManager) Select(sqlId string) *SqlRunner {
    return this.newSqlRunner(sqlId)
}

func (this *SqlManager) Update(sqlId string) *SqlRunner {
    return this.newSqlRunner(sqlId)
}

func (this *SqlManager) Delete(sqlId string) *SqlRunner {
    return this.newSqlRunner(sqlId)
}

func (this *SqlManager) Insert(sqlId string) *SqlRunner {
    return this.newSqlRunner(sqlId)
}

func (this *SqlRunner) Params(params ...interface{}) *SqlRunner {
    this.params = params
    return this
}

func (this *SqlRunner) Result(bean interface{}) *SqlRunner {
    rt := reflect.TypeOf(bean)
    if rt.Kind() != reflect.Ptr {
        return nil
    }

    rv := reflect.ValueOf(bean)
    rt = rt.Elem()
    rv = rv.Elem()

    switch rt.Kind() {
    case reflect.Slice:
        //FIXME: bean append in loop
        v, err := this.session.Select(this.resultHandler, this.sql, this.params...)
        if err == nil {
            rv.Set(reflect.ValueOf(v))
        }
        break
    case reflect.Struct:
        v, err := this.session.SelectOne(this.resultHandler, this.sql, this.params...)
        if err == nil {
            rv.Set(reflect.ValueOf(v))
        }
        break
    }
    return this
}
