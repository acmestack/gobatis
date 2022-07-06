/*
 * Copyright (c) 2022, AcmeStack
 * All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package gobatis

import (
	"context"

	"github.com/acmestack/gobatis/errors"
	"github.com/acmestack/gobatis/factory"
	"github.com/acmestack/gobatis/logging"
	"github.com/acmestack/gobatis/parsing/sqlparser"
	"github.com/acmestack/gobatis/reflection"
	"github.com/acmestack/gobatis/session"
)

type SessionManager struct {
	factory       factory.Factory
	ParserFactory ParserFactory
}

func NewSessionManager(factory factory.Factory) *SessionManager {
	return &SessionManager{
		factory:       factory,
		ParserFactory: DynamicParserFactory,
	}
}

type Runner interface {
	// Param 参数
	// 注意：如果没有参数也必须调用
	// 如果参数个数为1并且为struct，将解析struct获得参数
	// 如果参数个数大于1并且全部为简单类型，或则个数为1且为简单类型，则使用这些参数
	Param(params ...any) Runner
	// Result 获得结果
	Result(bean any) error
	// LastInsertId 最后插入的自增id
	LastInsertId() int64
	// Context 设置Context
	Context(ctx context.Context) Runner
}

type Session struct {
	ctx           context.Context
	log           logging.LogFunc
	session       session.SqlSession
	driver        string
	ParserFactory ParserFactory
}

type BaseRunner struct {
	session   session.SqlSession
	sqlParser sqlparser.SqlParser
	action    string
	metadata  *sqlparser.Metadata
	log       logging.LogFunc
	driver    string
	ctx       context.Context
	runner    Runner
}

type SelectRunner struct {
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

type ExecRunner struct {
	BaseRunner
}

// NewSession 使用一个session操作数据库
func (sessionManager *SessionManager) NewSession() *Session {
	return &Session{
		ctx:           context.Background(),
		log:           sessionManager.factory.LogFunc(),
		session:       sessionManager.factory.CreateSession(),
		driver:        sessionManager.factory.GetDataSource().DriverName(),
		ParserFactory: sessionManager.ParserFactory,
	}
}

// Context 包含session的context
func (sessionManager *SessionManager) Context(ctx context.Context) context.Context {
	sess := &Session{
		ctx:           ctx,
		log:           sessionManager.factory.LogFunc(),
		session:       sessionManager.factory.CreateSession(),
		driver:        sessionManager.factory.GetDataSource().DriverName(),
		ParserFactory: sessionManager.ParserFactory,
	}
	return context.WithValue(ctx, ContextSessionKey, sess)
}

func WithSession(ctx context.Context, sess *Session) context.Context {
	return context.WithValue(ctx, ContextSessionKey, sess)
}

func FindSession(ctx context.Context) *Session {
	if ctx == nil {
		return nil
	}
	return ctx.Value(ContextSessionKey).(*Session)
}

func (sessionManager *SessionManager) Close() error {
	return sessionManager.factory.Close()
}

// SetParserFactory 修改sql解析器创建者
func (sessionManager *SessionManager) SetParserFactory(fac ParserFactory) {
	sessionManager.ParserFactory = fac
}

func (session *Session) SetContext(ctx context.Context) *Session {
	session.ctx = ctx
	return session
}

func (session *Session) GetContext() context.Context {
	return session.ctx
}

// SetParserFactory 修改sql解析器创建者
func (session *Session) SetParserFactory(fac ParserFactory) {
	session.ParserFactory = fac
}

// Tx 开启事务执行语句
// 返回nil则提交，返回error回滚
// 抛出异常错误触发回滚
func (session *Session) Tx(txFunc func(session *Session) error) (err error) {
	e1 := session.session.Begin()
	if e1 != nil {
		return e1
	}
	defer func(err *error) {
		if r := recover(); r != nil {
			*err = session.session.Rollback()
			panic(r)
		}
	}(&err)

	if fnErr := txFunc(session); fnErr != nil {
		e := session.session.Rollback()
		if e != nil {
			session.log(logging.WARN, "Rollback error: %v , business error: %v\n", e, fnErr)
		}
		return fnErr
	} else {
		return session.session.Commit()
	}
}

func (session *Session) Select(sql string) Runner {
	return session.createSelect(session.findSqlParser(sql))
}

func (session *Session) Update(sql string) Runner {
	return session.createUpdate(session.findSqlParser(sql))
}

func (session *Session) Delete(sql string) Runner {
	return session.createDelete(session.findSqlParser(sql))
}

func (session *Session) Insert(sql string) Runner {
	return session.createInsert(session.findSqlParser(sql))
}

func (session *Session) Exec(sql string) Runner {
	return session.createExec(session.findSqlParser(sql))
}

func (baseRunner *BaseRunner) Param(params ...any) Runner {
	//TODO: 使用缓存加速，避免每次都生成动态sql
	//测试发现性能提升非常有限，故取消
	//key := cache.CalcKey(baseRunner.sqlDynamicData.OriginData, paramMap)
	//md := cache.FindMetadata(key)
	//var err error
	//if md == nil {
	//    md, err := baseRunner.sqlParser.Parse(params...)
	//    if err == nil {
	//        cache.CacheMetadata(key, md)
	//    }
	//}

	if baseRunner.sqlParser == nil {
		baseRunner.log(logging.WARN, errors.ParseParserNilError.Error())
		return baseRunner
	}

	md, err := baseRunner.sqlParser.ParseMetadata(baseRunner.driver, params...)

	if err == nil {
		if baseRunner.action == "" || baseRunner.action == md.Action {
			baseRunner.metadata = md
		} else {
			//allow different action
			baseRunner.log(logging.WARN, "sql action not match expect %s get %s", baseRunner.action, md.Action)
			baseRunner.metadata = md
		}
	} else {
		baseRunner.log(logging.WARN, err.Error())
	}
	return baseRunner.runner
}

//Context 设置执行的context
func (baseRunner *BaseRunner) Context(ctx context.Context) Runner {
	baseRunner.ctx = ctx
	return baseRunner.runner
}

func (selectRunner *SelectRunner) Result(bean any) error {
	if selectRunner.metadata == nil {
		selectRunner.log(logging.WARN, "Sql Metadata is nil")
		return errors.RunnerNotReady
	}

	if reflection.IsNil(bean) {
		return errors.ResultPointerIsNil
	}

	obj, err := ParseObject(bean)
	if err != nil {
		return err
	}
	return selectRunner.session.Query(selectRunner.ctx, obj, selectRunner.metadata.PrepareSql, selectRunner.metadata.Params...)

}

func (insertRunner *InsertRunner) Result(bean any) error {
	if insertRunner.metadata == nil {
		insertRunner.log(logging.WARN, "Sql Metadata is nil")
		return errors.RunnerNotReady
	}
	i, id, err := insertRunner.session.Insert(insertRunner.ctx, insertRunner.metadata.PrepareSql, insertRunner.metadata.Params...)
	insertRunner.lastId = id
	if reflection.CanSet(bean) {
		reflection.SetValue(reflection.ReflectValue(bean), i)
	}
	return err
}

func (insertRunner *InsertRunner) LastInsertId() int64 {
	return insertRunner.lastId
}

func (updateRunner *UpdateRunner) Result(bean any) error {
	if updateRunner.metadata == nil {
		updateRunner.log(logging.WARN, "Sql Metadata is nil")
		return errors.RunnerNotReady
	}
	i, err := updateRunner.session.Update(updateRunner.ctx, updateRunner.metadata.PrepareSql, updateRunner.metadata.Params...)
	if reflection.CanSet(bean) {
		reflection.SetValue(reflection.ReflectValue(bean), i)
	}
	return err
}

func (execRunner *ExecRunner) Result(bean any) error {
	if execRunner.metadata == nil {
		execRunner.log(logging.WARN, "Sql Metadata is nil")
		return errors.RunnerNotReady
	}
	i, err := execRunner.session.Update(execRunner.ctx, execRunner.metadata.PrepareSql, execRunner.metadata.Params...)
	if reflection.CanSet(bean) {
		reflection.SetValue(reflection.ReflectValue(bean), i)
	}
	return err
}

func (deleteRunner *DeleteRunner) Result(bean any) error {
	if deleteRunner.metadata == nil {
		deleteRunner.log(logging.WARN, "Sql Metadata is nil")
		return errors.RunnerNotReady
	}
	i, err := deleteRunner.session.Delete(deleteRunner.ctx, deleteRunner.metadata.PrepareSql, deleteRunner.metadata.Params...)
	if reflection.CanSet(bean) {
		reflection.SetValue(reflection.ReflectValue(bean), i)
	}
	return err
}

func (baseRunner *BaseRunner) Result(bean any) error {
	//FAKE RETURN
	panic("Cannot be here")
	//return nil, nil
}

func (baseRunner *BaseRunner) LastInsertId() int64 {
	return -1
}

func (session *Session) createSelect(parser sqlparser.SqlParser) Runner {
	ret := &SelectRunner{}
	ret.action = sqlparser.SELECT
	ret.log = session.log
	ret.session = session.session
	ret.sqlParser = parser
	ret.ctx = session.ctx
	ret.driver = session.driver
	ret.runner = ret
	return ret
}

func (session *Session) createUpdate(parser sqlparser.SqlParser) Runner {
	ret := &UpdateRunner{}
	ret.action = sqlparser.UPDATE
	ret.log = session.log
	ret.session = session.session
	ret.sqlParser = parser
	ret.ctx = session.ctx
	ret.driver = session.driver
	ret.runner = ret
	return ret
}

func (session *Session) createDelete(parser sqlparser.SqlParser) Runner {
	ret := &DeleteRunner{}
	ret.action = sqlparser.DELETE
	ret.log = session.log
	ret.session = session.session
	ret.sqlParser = parser
	ret.ctx = session.ctx
	ret.driver = session.driver
	ret.runner = ret
	return ret
}

func (session *Session) createInsert(parser sqlparser.SqlParser) Runner {
	ret := &InsertRunner{}
	ret.action = sqlparser.INSERT
	ret.log = session.log
	ret.session = session.session
	ret.sqlParser = parser
	ret.ctx = session.ctx
	ret.driver = session.driver
	ret.runner = ret
	return ret
}

func (session *Session) createExec(parser sqlparser.SqlParser) Runner {
	ret := &ExecRunner{}
	ret.action = ""
	ret.log = session.log
	ret.session = session.session
	ret.sqlParser = parser
	ret.ctx = session.ctx
	ret.driver = session.driver
	ret.runner = ret
	return ret
}

func (session *Session) findSqlParser(sqlId string) sqlparser.SqlParser {
	ret, ok := FindDynamicSqlParser(sqlId)
	if !ok {
		ret, ok = FindTemplateSqlParser(sqlId)
	}
	//FIXME: 当没有查找到sqlId对应的sql语句，则尝试使用sqlId直接操作数据库
	//该设计可能需要设计一个更合理的方式
	if !ok {
		d, err := session.ParserFactory(sqlId)
		if err != nil {
			session.log(logging.WARN, err.Error())
			return nil
		}
		return d
	}
	return ret
}
