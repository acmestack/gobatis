/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package runner

import (
    "github.com/xfali/gobatis"
    "github.com/xfali/gobatis/config"
    "github.com/xfali/gobatis/errors"
    "github.com/xfali/gobatis/factory"
    "github.com/xfali/gobatis/logging"
    "github.com/xfali/gobatis/parsing"
    "github.com/xfali/gobatis/parsing/sqlparser"
    "github.com/xfali/gobatis/reflection"
    "github.com/xfali/gobatis/session"
    "reflect"
)

type SessionManager struct {
    factory factory.Factory
}

func NewSessionManager(factory factory.Factory) *SessionManager {
    return &SessionManager{factory: factory}
}

type Runner interface {
    //参数
    //注意：如果没有参数也必须调用
    //如果参数个数为1并且为struct，将解析struct获得参数
    //如果参数个数大于1并且全部为简单类型，或则个数为1且为简单类型，则使用这些参数
    Param(params ...interface{}) Runner
    //获得结果
    Result(bean interface{}) error
}

type RunnerSession struct {
    log     logging.LogFunc
    session session.Session
}

type BaseRunner struct {
    session        session.Session
    sqlDynamicData parsing.DynamicData
    action         string
    metadata       *sqlparser.Metadata
    log            logging.LogFunc
    this           Runner
}

type SelectIterRunner struct {
    iterFunc gobatis.IterFunc
    count    int64
    BaseRunner
}

type SelectRunner struct {
    iterFunc gobatis.IterFunc
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

func getSql(sqlId string) *parsing.DynamicData {
    ret := config.FindSql(sqlId)
    //FIXME: 当没有查找到sqlId对应的sql语句，则尝试使用sqlId直接操作数据库
    //该设计可能需要设计一个更合理的方式
    if ret == nil {
        return &parsing.DynamicData{OriginData: sqlId}
    }
    return ret
}

//使用一个session操作数据库
func (this *SessionManager) NewSession() *RunnerSession {
    return &RunnerSession{
        log:     this.factory.LogFunc(),
        session: this.factory.CreateSession(),
    }
}

//开启事务执行语句
//返回true则提交，返回false回滚
//抛出异常错误触发回滚
func (this *RunnerSession) Tx(txFunc func(session *RunnerSession) bool) {
    this.session.Begin()
    defer func() {
        if r := recover(); r != nil {
            this.session.Rollback()
            panic(r)
        }
    }()

    if txFunc(this) != true {
        this.session.Rollback()
    } else {
        this.session.Commit()
    }
}

func (this *RunnerSession) SelectWithIterFunc(sqlId string, iterFunc gobatis.IterFunc) Runner {
    ret := &SelectIterRunner{}
    ret.action = sqlparser.SELECT
    ret.log = this.log
    ret.session = this.session
    ret.iterFunc = iterFunc
    ret.sqlDynamicData = *getSql(sqlId)
    ret.this = ret
    return ret
}

func (this *RunnerSession) Select(sqlId string) Runner {
    return createSelect(this.log, this.session, getSql(sqlId))
}

func (this *RunnerSession) Update(sqlId string) Runner {
    return createUpdate(this.log, this.session, getSql(sqlId))
}

func (this *RunnerSession) Delete(sqlId string) Runner {
    return createDelete(this.log, this.session, getSql(sqlId))
}

func (this *RunnerSession) Insert(sqlId string) Runner {
    return createInsert(this.log, this.session, getSql(sqlId))
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
    //TODO: 使用缓存加速，避免每次都生成动态sql
    sqlStr := this.sqlDynamicData.ReplaceWithParams(params...)
    md, err := sqlparser.ParseWithParams(sqlStr, params...)
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
    //TODO: 使用缓存加速，避免每次都生成动态sql
    sqlStr := this.sqlDynamicData.ReplaceWithMap(params)
    md, err := sqlparser.ParseWithParamMap(sqlStr, params)
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

func createSelect(log logging.LogFunc, session session.Session, sqlDynamic *parsing.DynamicData) Runner {
    ret := &SelectRunner{}
    ret.action = sqlparser.SELECT
    ret.log = log
    ret.session = session
    ret.sqlDynamicData = *sqlDynamic
    ret.this = ret
    return ret
}

func createUpdate(log logging.LogFunc, session session.Session, sqlDynamic *parsing.DynamicData) Runner {
    ret := &UpdateRunner{}
    ret.action = sqlparser.UPDATE
    ret.log = log
    ret.session = session
    ret.sqlDynamicData = *sqlDynamic
    ret.this = ret
    return ret
}

func createDelete(log logging.LogFunc, session session.Session, sqlDynamic *parsing.DynamicData) Runner {
    ret := &DeleteRunner{}
    ret.action = sqlparser.DELETE
    ret.log = log
    ret.session = session
    ret.sqlDynamicData = *sqlDynamic
    ret.this = ret
    return ret
}

func createInsert(log logging.LogFunc, session session.Session, sqlDynamic *parsing.DynamicData) Runner {
    ret := &InsertRunner{}
    ret.action = sqlparser.INSERT
    ret.log = log
    ret.session = session
    ret.sqlDynamicData = *sqlDynamic
    ret.this = ret
    return ret
}
