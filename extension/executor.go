/*
 * Copyright (C) 2022, Xiongfa Li.
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

package extension

import (
	"context"
	"github.com/acmestack/gobatis/common"
	"github.com/acmestack/gobatis/executor"
	"github.com/acmestack/gobatis/reflection"
	"github.com/xfali/aop"
	"runtime"
	"strings"
)

type executorEx struct {
	proxy aop.Proxy
}

func NewExecutorExtension(e executor.Executor) *executorEx {
	ret := &executorEx{
		proxy: aop.New(e),
	}
	return ret
}

func (e *executorEx) Extend(pointCut aop.PointCut, advice aop.Advice) Extension {
	e.proxy.AddAdvisor(pointCut, advice)
	return e
}

func (e *executorEx) Close(rollback bool) {
	e.proxy.Call(caller(), rollback)
}

func (e *executorEx) Query(ctx context.Context, result reflection.Object, sql string, params ...interface{}) error {
	r, err := e.proxy.Call(caller(), ctx, result, sql, params)
	if err != nil {
		return err
	}
	if r[0] == nil {
		return nil
	}
	return r[0].(error)
}

func (e *executorEx) Exec(ctx context.Context, sql string, params ...interface{}) (res common.Result, rerr error) {
	r, err := e.proxy.Call(caller(), ctx, sql, params)
	if err != nil {
		return nil, err
	}
	if r[0] != nil {
		res = r[0].(common.Result)
	}
	if r[1] != nil {
		rerr = r[1].(error)
	}
	return
}

func (e *executorEx) Begin() error {
	r, err := e.proxy.Call(caller())
	if err != nil {
		return err
	}
	if r[0] == nil {
		return nil
	}
	return r[0].(error)
}

func (e *executorEx) Commit(require bool) error {
	r, err := e.proxy.Call(caller(), require)
	if err != nil {
		return err
	}
	if r[0] == nil {
		return nil
	}
	return r[0].(error)
}

func (e *executorEx) Rollback(require bool) error {
	r, err := e.proxy.Call(caller(), require)
	if err != nil {
		return err
	}
	if r[0] == nil {
		return nil
	}
	return r[0].(error)
}

func caller() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	name := f.Name()
	if i := strings.LastIndex(name, "."); i != -1 {
		return name[i+1:]
	}
	return name
}
