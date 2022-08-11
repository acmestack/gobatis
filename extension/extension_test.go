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
	"fmt"
	"github.com/acmestack/gobatis/executor"
	"github.com/acmestack/gobatis/reflection"
	"github.com/xfali/aop"
	"testing"
)

func TestExecutorExtension(t *testing.T) {
	ex := NewExecutorExtension(executor.NewDummyExecutor())
	ex.Extend(aop.PointCutRegExp("", "(.*?)", nil, nil), func(invocation aop.Invocation, params []interface{}) (ret []interface{}) {
		t.Log("params: ", fmt.Sprintln(params...))
		ret = invocation.Invoke(params)
		if len(ret) > 0 {
			t.Log("result: ", fmt.Sprintln(ret...))
		} else {
			t.Log("result: nil")
		}
		return ret
	})
	var i int = 0
	result, err := reflection.GetObjectInfo(&i)
	if err != nil {
		t.Fatal(err)
	}
	ex.Begin()
	ex.Exec(context.Background(), "update tbl where id = ? and name = ?", "hello", "world")
	ex.Query(context.Background(), result, "select * from tbl where id = ? and name = ?", "hello", "world")
	ex.Commit(false)
	ex.Rollback(false)
	ex.Close(false)
}
