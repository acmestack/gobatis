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

type Runner interface {
    Param(params ...interface{}) Runner
    Result(bean interface{}) error
}

type BaseRunner struct {
    session  session.Session
    sql      string
    action   string
    metadata *sqlparser.Metadata
    log      logging.LogFunc
    this     Runner
}

type SelectIterRunner struct {
    iterFunc session.IterFunc
    count    int64
    BaseRunner
}

type SelectRunner struct {
    iterFunc session.IterFunc
    count    int64
    BaseRunner
}

type InsertRunner struct {
    BaseRunner
}

type UpdateRunner struct {
    BaseRunner
}

type DeleteRunner struct {
    BaseRunner
}

func (this *SqlManager) GetSql(sqlId string) string {
    return this.sqlmap[sqlId]
}

func (this *SqlManager) SelectWithIterFunc(sqlId string, iterFunc session.IterFunc) Runner {
    ret := &SelectIterRunner{}
    ret.action = sqlparser.SELECT
    ret.log = this.factory.LogFunc()
    ret.session = this.factory.CreateSession()
    ret.iterFunc = iterFunc
    ret.sql = this.GetSql(sqlId)
    ret.this = ret
    return ret
}

func (this *SqlManager) Select(sqlId string) Runner {
    ret := &SelectRunner{}
    ret.action = sqlparser.SELECT
    ret.log = this.factory.LogFunc()
    ret.session = this.factory.CreateSession()
    ret.sql = this.GetSql(sqlId)
    ret.this = ret
    return ret
}

func (this *SqlManager) Update(sqlId string) Runner {
    ret := &UpdateRunner{}
    ret.action = sqlparser.UPDATE
    ret.log = this.factory.LogFunc()
    ret.session = this.factory.CreateSession()
    ret.sql = this.GetSql(sqlId)
    ret.this = ret
    return ret
}

func (this *SqlManager) Delete(sqlId string) Runner {
    ret := &DeleteRunner{}
    ret.action = sqlparser.DELETE
    ret.log = this.factory.LogFunc()
    ret.session = this.factory.CreateSession()
    ret.sql = this.GetSql(sqlId)
    ret.this = ret
    return ret
}

func (this *SqlManager) Insert(sqlId string) Runner {
    ret := &InsertRunner{}
    ret.action = sqlparser.INSERT
    ret.log = this.factory.LogFunc()
    ret.session = this.factory.CreateSession()
    ret.sql = this.GetSql(sqlId)
    ret.this = ret
    return ret
}

func (this *BaseRunner) Param(params ...interface{}) Runner {
    if len(params) == 0 {
        return this.params()
    } else if len(params) == 1 {
        t := reflect.TypeOf(params[0])
        if t.Kind() == reflect.Ptr {
            t = t.Elem()
        }
        if t.Kind() == reflect.Struct {
            return this.paramType(params[0])
        } else {
            if reflection.IsSimpleType(params[0]) {
                return this.params(params...)
            }
        }
    } else {
        for _, v := range params {
            if !reflection.IsSimpleType(v) {
                this.log(logging.WARN, "Param error: expect simple type, but get other type")
                return this.this
            }
        }
        return this.params(params...)
    }
    return this.this
}


func (this *BaseRunner) params(params ...interface{}) Runner {
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
    return this.this
}

func (this *BaseRunner) paramType(paramVar interface{}) Runner {
    this.metadata = nil
    ti, err := reflection.GetTableInfo(paramVar)
    if err != nil {
        this.log(logging.WARN, "%s", err.Error())
        return this.this
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
    return this.this
}

func (this *SelectIterRunner) myIterFunc(idx int64, bean interface{}) bool {
    this.count++
    return this.iterFunc(idx, bean)
}

func (this *SelectIterRunner) Result(bean interface{}) error {
    err := checkBeanValue(reflect.ValueOf(bean))
    if err != nil {
        return err
    }
    if this.metadata == nil {
        this.log(logging.WARN, "Sql Matadata is nil")
        return errors.RUNNER_NOT_READY
    }

    mi := config.FindModelInfoOfBean(bean)
    if mi == nil {
        this.log(logging.WARN, errors.MODEL_NOT_REGISTER.Error())
        return errors.RESULT_NAME_NOT_FOUND
    }
    return this.session.Query(mi, this.myIterFunc, this.metadata.PrepareSql, this.metadata.Params...)
}

func (this *SelectRunner) Result(bean interface{}) error {
    if this.metadata == nil {
        this.log(logging.WARN, "Sql Matadata is nil")
        return errors.RUNNER_NOT_READY
    }

    mi := config.FindModelInfoOfBean(bean)
    if mi == nil {
        this.log(logging.WARN, errors.MODEL_NOT_REGISTER.Error())
        return errors.RESULT_NAME_NOT_FOUND
    }
    rt := reflect.TypeOf(bean)
    rv := reflect.ValueOf(bean)
    err := checkBeanValue(rv)
    if err != nil {
        return err
    }
    rt = rt.Elem()
    rv = rv.Elem()

    switch rt.Kind() {
    case reflect.Slice:
        //FIXME: bean append in loop
        retV := rv
        iterFunc := func(idx int64, bean interface{}) bool {
            retV = reflect.Append(retV, reflect.ValueOf(bean))
            return false
        }
        err := this.session.Query(mi, iterFunc, this.metadata.PrepareSql, this.metadata.Params...)
        if err == nil {
            rv.Set(retV)
        } else {
            return err
        }
        break
    default:
        v, err := this.session.SelectOne(mi, this.metadata.PrepareSql, this.metadata.Params...)
        if err == nil {
            retV := reflect.ValueOf(v)
            if retV.IsValid() {
                rv.Set(reflect.ValueOf(v))
            } else {
                return errors.RESULT_SELECT_EMPTY_VALUE
            }
        } else {
            return err
        }
        break
    }
    return nil
}

func (this *InsertRunner) Result(bean interface{}) error {
    if this.metadata == nil {
        this.log(logging.WARN, "Sql Matadata is nil")
        return errors.RUNNER_NOT_READY
    }
    var rv reflect.Value
    if bean != nil {
        rv = reflect.ValueOf(bean)
        err := checkBeanValue(rv)
        rv = rv.Elem()
        if err != nil {
            return err
        }
    }
    i := this.session.Insert(this.metadata.PrepareSql, this.metadata.Params...)
    if bean != nil {
        reflection.SetValue(rv, i)
    }
    return nil
}

func (this *UpdateRunner) Result(bean interface{}) error {
    if this.metadata == nil {
        this.log(logging.WARN, "Sql Matadata is nil")
        return errors.RUNNER_NOT_READY
    }
    var rv reflect.Value
    if bean != nil {
        rv = reflect.ValueOf(bean)
        err := checkBeanValue(rv)
        rv = rv.Elem()
        if err != nil {
            return err
        }
    }
    i := this.session.Update(this.metadata.PrepareSql, this.metadata.Params...)
    if bean != nil {
        reflection.SetValue(rv, i)
    }
    return nil
}

func (this *DeleteRunner) Result(bean interface{}) error {
    if this.metadata == nil {
        this.log(logging.WARN, "Sql Matadata is nil")
        return errors.RUNNER_NOT_READY
    }
    var rv reflect.Value
    if bean != nil {
        rv = reflect.ValueOf(bean)
        err := checkBeanValue(rv)
        rv = rv.Elem()
        if err != nil {
            return err
        }
    }
    i := this.session.Delete(this.metadata.PrepareSql, this.metadata.Params...)
    if bean != nil {
        reflection.SetValue(rv, i)
    }
    return nil
}

func (this *BaseRunner) Result(bean interface{}) error {
    //FAKE RETURN
    panic("Cannot be here")
    //return nil, nil
}

func checkBeanValue(beanValue reflect.Value) error {
    if beanValue.Kind() != reflect.Ptr {
        return errors.RESULT_ISNOT_POINTER
    } else if beanValue.Elem().Kind() == reflect.Ptr {
        return errors.RESULT_PTR_VALUE_IS_POINTER
    }
    return nil
}

func (this *BaseRunner) ResultBad(bean interface{}) *BaseRunner {
    panic("Cannot be here")
    if this.metadata == nil {
        this.log(logging.WARN, "Sql Matadata is nil")
        return nil
    }

    mi := config.FindModelInfoOfBean(bean)
    if mi == nil {
        this.log(logging.WARN, errors.MODEL_NOT_REGISTER.Error())
        return nil
    }

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
        v, err := this.session.Select(mi, this.metadata.PrepareSql, this.metadata.Params...)
        if err == nil {
            rv.Set(reflect.ValueOf(v))
        }
        break
    case reflect.Struct:
        v, err := this.session.SelectOne(mi, this.metadata.PrepareSql, this.metadata.Params...)
        if err == nil {
            retV := reflect.ValueOf(v)
            if retV.IsValid() {
                rv.Set(reflect.ValueOf(v))
            } else {
                return nil
            }
        }
        break
    }
    return this
}
