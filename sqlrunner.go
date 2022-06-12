/*
 * Copyright (c) 2022, OpeningO
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
	"github.com/xfali/gobatis/errors"
	"github.com/xfali/gobatis/factory"
	"github.com/xfali/gobatis/logging"
	"github.com/xfali/gobatis/parsing/sqlparser"
	"github.com/xfali/gobatis/reflection"
	"github.com/xfali/gobatis/session"
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
	// 参数
	// 注意：如果没有参数也必须调用
	// 如果参数个数为1并且为struct，将解析struct获得参数
	// 如果参数个数大于1并且全部为简单类型，或则个数为1且为简单类型，则使用这些参数
	Param(params ...interface{}) Runner
	// 获得结果
	Result(bean interface{}) error
	// 最后插入的自增id
	LastInsertId() int64
	// 设置Context
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
	this      Runner
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

// 使用一个session操作数据库
func (this *SessionManager) NewSession() *Session {
	return &Session{
		ctx:           context.Background(),
		log:           this.factory.LogFunc(),
		session:       this.factory.CreateSession(),
		driver:        this.factory.GetDataSource().DriverName(),
		ParserFactory: this.ParserFactory,
	}
}

// 包含session的context
func (this *SessionManager) Context(ctx context.Context) context.Context {
	sess := &Session{
		ctx:           ctx,
		log:           this.factory.LogFunc(),
		session:       this.factory.CreateSession(),
		driver:        this.factory.GetDataSource().DriverName(),
		ParserFactory: this.ParserFactory,
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

func (this *SessionManager) Close() error {
	return this.factory.Close()
}

// 修改sql解析器创建者
func (this *SessionManager) SetParserFactory(fac ParserFactory) {
	this.ParserFactory = fac
}

func (this *Session) SetContext(ctx context.Context) *Session {
	this.ctx = ctx
	return this
}

func (this *Session) GetContext() context.Context {
	return this.ctx
}

// 修改sql解析器创建者
func (this *Session) SetParserFactory(fac ParserFactory) {
	this.ParserFactory = fac
}

// 开启事务执行语句
// 返回nil则提交，返回error回滚
// 抛出异常错误触发回滚
func (this *Session) Tx(txFunc func(session *Session) error) (err error) {
	e1 := this.session.Begin()
	if e1 != nil {
		return e1
	}
	defer func(err *error) {
		if r := recover(); r != nil {
			*err = this.session.Rollback()
			panic(r)
		}
	}(&err)

	if fnErr := txFunc(this); fnErr != nil {
		e := this.session.Rollback()
		if e != nil {
			this.log(logging.WARN, "Rollback error: %v , business error: %v\n", e, fnErr)
		}
		return fnErr
	} else {
		return this.session.Commit()
	}
}

func (this *Session) Select(sql string) Runner {
	return this.createSelect(this.findSqlParser(sql))
}

func (this *Session) Update(sql string) Runner {
	return this.createUpdate(this.findSqlParser(sql))
}

func (this *Session) Delete(sql string) Runner {
	return this.createDelete(this.findSqlParser(sql))
}

func (this *Session) Insert(sql string) Runner {
	return this.createInsert(this.findSqlParser(sql))
}

func (this *Session) Exec(sql string) Runner {
	return this.createExec(this.findSqlParser(sql))
}

func (this *BaseRunner) Param(params ...interface{}) Runner {
	//TODO: 使用缓存加速，避免每次都生成动态sql
	//测试发现性能提升非常有限，故取消
	//key := cache.CalcKey(this.sqlDynamicData.OriginData, paramMap)
	//md := cache.FindMetadata(key)
	//var err error
	//if md == nil {
	//    md, err := this.sqlParser.Parse(params...)
	//    if err == nil {
	//        cache.CacheMetadata(key, md)
	//    }
	//}

	if this.sqlParser == nil {
		this.log(logging.WARN, errors.PARSE_PARSER_NIL_ERROR.Error())
		return this
	}

	md, err := this.sqlParser.ParseMetadata(this.driver, params...)

	if err == nil {
		if this.action == "" || this.action == md.Action {
			this.metadata = md
		} else {
			//allow different action
			this.log(logging.WARN, "sql action not match expect %s get %s", this.action, md.Action)
			this.metadata = md
		}
	} else {
		this.log(logging.WARN, err.Error())
	}
	return this.this
}

//Context 设置执行的context
func (this *BaseRunner) Context(ctx context.Context) Runner {
	this.ctx = ctx
	return this.this
}

func (this *SelectRunner) Result(bean interface{}) error {
	if this.metadata == nil {
		this.log(logging.WARN, "Sql Matadata is nil")
		return errors.RUNNER_NOT_READY
	}

	if reflection.IsNil(bean) {
		return errors.RESULT_POINTER_IS_NIL
	}

	obj, err := ParseObject(bean)
	if err != nil {
		return err
	}
	return this.session.Query(this.ctx, obj, this.metadata.PrepareSql, this.metadata.Params...)

}

func (this *InsertRunner) Result(bean interface{}) error {
	if this.metadata == nil {
		this.log(logging.WARN, "Sql Matadata is nil")
		return errors.RUNNER_NOT_READY
	}
	i, id, err := this.session.Insert(this.ctx, this.metadata.PrepareSql, this.metadata.Params...)
	this.lastId = id
	if reflection.CanSet(bean) {
		reflection.SetValue(reflection.ReflectValue(bean), i)
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
	i, err := this.session.Update(this.ctx, this.metadata.PrepareSql, this.metadata.Params...)
	if reflection.CanSet(bean) {
		reflection.SetValue(reflection.ReflectValue(bean), i)
	}
	return err
}

func (this *ExecRunner) Result(bean interface{}) error {
	if this.metadata == nil {
		this.log(logging.WARN, "Sql Matadata is nil")
		return errors.RUNNER_NOT_READY
	}
	i, err := this.session.Update(this.ctx, this.metadata.PrepareSql, this.metadata.Params...)
	if reflection.CanSet(bean) {
		reflection.SetValue(reflection.ReflectValue(bean), i)
	}
	return err
}

func (this *DeleteRunner) Result(bean interface{}) error {
	if this.metadata == nil {
		this.log(logging.WARN, "Sql Matadata is nil")
		return errors.RUNNER_NOT_READY
	}
	i, err := this.session.Delete(this.ctx, this.metadata.PrepareSql, this.metadata.Params...)
	if reflection.CanSet(bean) {
		reflection.SetValue(reflection.ReflectValue(bean), i)
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

func (this *Session) createSelect(parser sqlparser.SqlParser) Runner {
	ret := &SelectRunner{}
	ret.action = sqlparser.SELECT
	ret.log = this.log
	ret.session = this.session
	ret.sqlParser = parser
	ret.ctx = this.ctx
	ret.driver = this.driver
	ret.this = ret
	return ret
}

func (this *Session) createUpdate(parser sqlparser.SqlParser) Runner {
	ret := &UpdateRunner{}
	ret.action = sqlparser.UPDATE
	ret.log = this.log
	ret.session = this.session
	ret.sqlParser = parser
	ret.ctx = this.ctx
	ret.driver = this.driver
	ret.this = ret
	return ret
}

func (this *Session) createDelete(parser sqlparser.SqlParser) Runner {
	ret := &DeleteRunner{}
	ret.action = sqlparser.DELETE
	ret.log = this.log
	ret.session = this.session
	ret.sqlParser = parser
	ret.ctx = this.ctx
	ret.driver = this.driver
	ret.this = ret
	return ret
}

func (this *Session) createInsert(parser sqlparser.SqlParser) Runner {
	ret := &InsertRunner{}
	ret.action = sqlparser.INSERT
	ret.log = this.log
	ret.session = this.session
	ret.sqlParser = parser
	ret.ctx = this.ctx
	ret.driver = this.driver
	ret.this = ret
	return ret
}

func (this *Session) createExec(parser sqlparser.SqlParser) Runner {
	ret := &ExecRunner{}
	ret.action = ""
	ret.log = this.log
	ret.session = this.session
	ret.sqlParser = parser
	ret.ctx = this.ctx
	ret.driver = this.driver
	ret.this = ret
	return ret
}

func (this *Session) findSqlParser(sqlId string) sqlparser.SqlParser {
	ret, ok := FindDynamicSqlParser(sqlId)
	if !ok {
		ret, ok = FindTemplateSqlParser(sqlId)
	}
	//FIXME: 当没有查找到sqlId对应的sql语句，则尝试使用sqlId直接操作数据库
	//该设计可能需要设计一个更合理的方式
	if !ok {
		d, err := this.ParserFactory(sqlId)
		if err != nil {
			this.log(logging.WARN, err.Error())
			return nil
		}
		return d
	}
	return ret
}
