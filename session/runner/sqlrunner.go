/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package runner

import (
    "github.com/xfali/gobatis/config"
    "github.com/xfali/gobatis/errors"
    "github.com/xfali/gobatis/factory"
    "github.com/xfali/gobatis/handler"
    "github.com/xfali/gobatis/logging"
    "github.com/xfali/gobatis/parsing/sqlparser"
    "github.com/xfali/gobatis/reflection"
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
    action        string
    metadata      *sqlparser.Metadata
    log           logging.LogFunc
}

func (this *SqlManager) newSqlRunner(action, sqlId string) *SqlRunner {
    if sql, ok := this.sqlmap[sqlId]; ok {
        return &SqlRunner{action: action, log: this.factory.LogFunc(), session: this.factory.CreateSession(), sql: sql}
    }
    return nil
}

func (this *SqlManager) Select(sqlId string) *SqlRunner {
    return this.newSqlRunner(sqlparser.SELECT, sqlId)
}

func (this *SqlManager) Update(sqlId string) *SqlRunner {
    return this.newSqlRunner(sqlparser.UPDATE, sqlId)
}

func (this *SqlManager) Delete(sqlId string) *SqlRunner {
    return this.newSqlRunner(sqlparser.DELETE, sqlId)
}

func (this *SqlManager) Insert(sqlId string) *SqlRunner {
    return this.newSqlRunner(sqlparser.INSERT, sqlId)
}

func (this *SqlRunner) Params(params ...interface{}) *SqlRunner {
    this.metadata = nil
    md, err := sqlparser.ParseWithParams(this.sql, params...)
    if err == nil {
        if this.action == md.Action {
            this.metadata = md
        } else {
            this.log(logging.WARN, "sql action not match expect %s get %s", this.action, md.Action)
        }
    } else {
        this.log(logging.WARN, "%s", err.Error())
    }
    return this
}

func (this *SqlRunner) ParamType(paramVar interface{}) *SqlRunner {
    this.metadata = nil
    ti, err := reflection.GetTableInfo(&paramVar)
    if err != nil {
        return this
    } else {
        this.log(logging.WARN, "%s", err.Error())
    }
    params := ti.MapValue()
    md, err := sqlparser.ParseWithParamMap(this.sql, params)
    if err == nil {
        if this.action == md.Action {
            this.metadata = md
        } else {
            this.log(logging.WARN, "sql action not match expect %s get %s", this.action, md.Action)
        }
    } else {
        this.log(logging.WARN, "%s", err.Error())
    }
    return this
}

func (this *SqlRunner) Result(bean interface{}) *SqlRunner {
    if this.metadata == nil {
        this.log(logging.WARN, "Sql Matadata is nil")
        return nil
    }

    mi := config.FindModelInfoOfBean(bean)
    if mi == nil {
        this.log(logging.WARN, errors.MODEL_NOT_REGISTER.Error())
        return nil
    }
    this.resultHandler = mi

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
        v, err := this.session.Select(this.resultHandler, this.metadata.PrepareSql, this.metadata.Params...)
        if err == nil {
            rv.Set(reflect.ValueOf(v))
        }
        break
    case reflect.Struct:
        v, err := this.session.SelectOne(this.resultHandler, this.metadata.PrepareSql, this.metadata.Params...)
        if err == nil {
            rv.Set(reflect.ValueOf(v))
        }
        break
    }
    return this
}
