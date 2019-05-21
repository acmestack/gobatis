/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package gobatis

import (
    "context"
    "github.com/xfali/gobatis/common"
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
    //最后插入的自增id
    LastInsertId() int64
    //设置Context
    Context(ctx context.Context) Runner
}

type Session struct {
    ctx     context.Context
    log     logging.LogFunc
    session session.SqlSession
}

type BaseRunner struct {
    session        session.SqlSession
    sqlDynamicData parsing.DynamicData
    action         string
    metadata       *sqlparser.Metadata
    log            logging.LogFunc
    ctx            context.Context
    this           Runner
}

type SelectIterRunner struct {
    iterFunc common.IterFunc
    count    int64
    BaseRunner
}

type SelectRunner struct {
    iterFunc common.IterFunc
    count    int64
    BaseRunner
}

type InsertRunner struct {
    lastId int64
    BaseRunner
}

type UpdateRunner struct {
    BaseRunner
}

type DeleteRunner struct {
    BaseRunner
}

func getSql(sqlId string) *parsing.DynamicData {
    ret := FindSql(sqlId)
    //FIXME: 当没有查找到sqlId对应的sql语句，则尝试使用sqlId直接操作数据库
    //该设计可能需要设计一个更合理的方式
    if ret == nil {
        return &parsing.DynamicData{OriginData: sqlId}
    }
    return ret
}

//使用一个session操作数据库
func (this *SessionManager) NewSession() *Session {
    return &Session{
        ctx:     context.Background(),
        log:     this.factory.LogFunc(),
        session: this.factory.CreateSession(),
    }
}

func (this *Session) SetContext(ctx context.Context) *Session {
    this.ctx = ctx
    return this
}

func (this *Session) GetContext() context.Context {
    return this.ctx
}

//开启事务执行语句
//返回true则提交，返回false回滚
//抛出异常错误触发回滚
func (this *Session) Tx(txFunc func(session *Session) bool) {
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

func (this *Session) SelectWithIterFunc(sqlId string, iterFunc common.IterFunc) Runner {
    ret := &SelectIterRunner{}
    ret.action = sqlparser.SELECT
    ret.log = this.log
    ret.session = this.session
    ret.iterFunc = iterFunc
    ret.sqlDynamicData = *getSql(sqlId)
    ret.this = ret
    return ret
}

func (this *Session) Select(sql string) Runner {
    return createSelect(this.ctx, this.log, this.session, getSql(sql))
}

func (this *Session) Update(sql string) Runner {
    return createUpdate(this.ctx, this.log, this.session, getSql(sql))
}

func (this *Session) Delete(sql string) Runner {
    return createDelete(this.ctx, this.log, this.session, getSql(sql))
}

func (this *Session) Insert(sql string) Runner {
    return createInsert(this.ctx, this.log, this.session, getSql(sql))
}

func (this *BaseRunner) Param(params ...interface{}) Runner {
    paramMap := reflection.ParseParams(params...)
    //TODO: 使用缓存加速，避免每次都生成动态sql
    //测试发现性能提升非常有限，故取消
    //key := cache.CalcKey(this.sqlDynamicData.OriginData, paramMap)
    //md := cache.FindMetadata(key)
    //var err error
    //if md == nil {
    //    sqlStr := this.sqlDynamicData.ReplaceWithMap(paramMap)
    //    md, err = sqlparser.ParseWithParamMap(sqlStr, paramMap)
    //    if err == nil {
    //        cache.CacheMetadata(key, md)
    //    }
    //}

    sqlStr := this.sqlDynamicData.ReplaceWithMap(paramMap)
    md, err := sqlparser.ParseWithParamMap(sqlStr, paramMap)
    
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

//Context 设置执行的context
func (this *BaseRunner) Context(ctx context.Context) Runner {
    this.ctx = ctx
    return this.this
}

func (this *SelectIterRunner) myIterFunc(idx int64, bean interface{}) bool {
    this.count++
    return this.iterFunc(idx, bean)
}

func (this *SelectIterRunner) Result(bean interface{}) error {
    err := reflection.MustPtr(bean)
    if err != nil {
        this.log(logging.WARN, "Result error: %s\n", err.Error())
        return err
    }
    if this.metadata == nil {
        this.log(logging.WARN, "Sql Matadata is nil")
        return errors.RUNNER_NOT_READY
    }

    mi := FindModelInfoOfBean(bean)
    if mi == nil {
        this.log(logging.WARN, errors.MODEL_NOT_REGISTER.Error())
        return errors.RESULT_NAME_NOT_FOUND
    }
    return this.session.Query(this.ctx, mi, this.myIterFunc, this.metadata.PrepareSql, this.metadata.Params...)
}

func (this *SelectRunner) Result(bean interface{}) error {
    if this.metadata == nil {
        this.log(logging.WARN, "Sql Matadata is nil")
        return errors.RUNNER_NOT_READY
    }

    if reflection.IsNil(bean) {
        return errors.RESULT_POINTER_IS_NIL
    }

    //mi := FindModelInfoOfBean(bean)
    //if mi == nil {
    //    this.log(logging.WARN, errors.MODEL_NOT_REGISTER.Error())
    //    return errors.RESULT_NAME_NOT_FOUND
    //}
    obj, err := ParseObject(bean)
    if err != nil {
        return err
    }
    iterFunc := func(idx int64, bean interface{}) bool {
        return !obj.obj.CanAddValue()
    }
    return this.session.Query(this.ctx, &obj, iterFunc, this.metadata.PrepareSql, this.metadata.Params...)

}

func (this *InsertRunner) Result(bean interface{}) error {
    if this.metadata == nil {
        this.log(logging.WARN, "Sql Matadata is nil")
        return errors.RUNNER_NOT_READY
    }
    var rv reflect.Value
    if bean != nil {
        rv = reflect.ValueOf(bean)
        err := reflection.MustPtrValue(rv)
        rv = rv.Elem()
        if err != nil {
            return err
        }
    }
    i, id, err := this.session.Insert(this.ctx, this.metadata.PrepareSql, this.metadata.Params...)
    this.lastId = id
    if bean != nil {
        reflection.SetValue(rv, i)
    }
    return err
}

func (this *InsertRunner) LastInsertId() int64 {
    return this.lastId
}

func (this *UpdateRunner) Result(bean interface{}) error {
    if this.metadata == nil {
        this.log(logging.WARN, "Sql Matadata is nil")
        return errors.RUNNER_NOT_READY
    }
    var rv reflect.Value
    if bean != nil {
        rv = reflect.ValueOf(bean)
        err := reflection.MustPtrValue(rv)
        rv = rv.Elem()
        if err != nil {
            return err
        }
    }
    i, err := this.session.Update(this.ctx, this.metadata.PrepareSql, this.metadata.Params...)
    if bean != nil {
        reflection.SetValue(rv, i)
    }
    return err
}

func (this *DeleteRunner) Result(bean interface{}) error {
    if this.metadata == nil {
        this.log(logging.WARN, "Sql Matadata is nil")
        return errors.RUNNER_NOT_READY
    }
    var rv reflect.Value
    if bean != nil {
        rv = reflect.ValueOf(bean)
        err := reflection.MustPtrValue(rv)
        rv = rv.Elem()
        if err != nil {
            return err
        }
    }
    i, err := this.session.Delete(this.ctx, this.metadata.PrepareSql, this.metadata.Params...)
    if bean != nil {
        reflection.SetValue(rv, i)
    }
    return err
}

func (this *BaseRunner) Result(bean interface{}) error {
    //FAKE RETURN
    panic("Cannot be here")
    //return nil, nil
}

func (this *BaseRunner) LastInsertId() int64 {
    return -1
}

func createSelect(ctx context.Context, log logging.LogFunc, session session.SqlSession, sqlDynamic *parsing.DynamicData) Runner {
    ret := &SelectRunner{}
    ret.action = sqlparser.SELECT
    ret.log = log
    ret.session = session
    ret.sqlDynamicData = *sqlDynamic
    ret.ctx = ctx
    ret.this = ret
    return ret
}

func createUpdate(ctx context.Context, log logging.LogFunc, session session.SqlSession, sqlDynamic *parsing.DynamicData) Runner {
    ret := &UpdateRunner{}
    ret.action = sqlparser.UPDATE
    ret.log = log
    ret.session = session
    ret.sqlDynamicData = *sqlDynamic
    ret.ctx = ctx
    ret.this = ret
    return ret
}

func createDelete(ctx context.Context, log logging.LogFunc, session session.SqlSession, sqlDynamic *parsing.DynamicData) Runner {
    ret := &DeleteRunner{}
    ret.action = sqlparser.DELETE
    ret.log = log
    ret.session = session
    ret.sqlDynamicData = *sqlDynamic
    ret.ctx = ctx
    ret.this = ret
    return ret
}

func createInsert(ctx context.Context, log logging.LogFunc, session session.SqlSession, sqlDynamic *parsing.DynamicData) Runner {
    ret := &InsertRunner{}
    ret.action = sqlparser.INSERT
    ret.log = log
    ret.session = session
    ret.sqlDynamicData = *sqlDynamic
    ret.ctx = ctx
    ret.this = ret
    return ret
}
