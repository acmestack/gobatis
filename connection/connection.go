/*
 * Licensed to the AcmeStack under one or more contributor license
 * agreements. See the NOTICE file distributed with this work for
 * additional information regarding copyright ownership.
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

package connection

import (
	"context"
	"github.com/acmestack/gobatis/common"
	"github.com/acmestack/gobatis/reflection"
	"github.com/acmestack/gobatis/statement"
)

type Connection interface {
	Prepare(sql string) (statement.Statement, error)
	Query(ctx context.Context, result reflection.Object, sql string, params ...interface{}) error
	Exec(ctx context.Context, sql string, params ...interface{}) (common.Result, error)
}
