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

package executor

import (
	"context"
	"github.com/acmestack/gobatis/common"
	"github.com/acmestack/gobatis/reflection"
)

type dummyExecutor struct{}

func NewDummyExecutor() dummyExecutor {
	return dummyExecutor{}
}

func (e dummyExecutor) Close(rollback bool) {

}

func (e dummyExecutor) Query(ctx context.Context, result reflection.Object, sql string, params ...interface{}) error {
	return nil
}

func (e dummyExecutor) Exec(ctx context.Context, sql string, params ...interface{}) (common.Result, error) {
	return nil, nil
}

func (e dummyExecutor) Begin() error {
	return nil
}

func (e dummyExecutor) Commit(require bool) error {
	return nil
}

func (e dummyExecutor) Rollback(require bool) error {
	return nil
}
