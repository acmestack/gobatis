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

package session

import (
	"context"
	"github.com/acmestack/gobatis/reflection"
)

type SqlSession interface {
	Close(rollback bool)

	Query(ctx context.Context, result reflection.Object, sql string, params ...interface{}) error

	Insert(ctx context.Context, sql string, params ...interface{}) (int64, int64, error)

	Update(ctx context.Context, sql string, params ...interface{}) (int64, error)

	Delete(ctx context.Context, sql string, params ...interface{}) (int64, error)

	Begin() error

	Commit() error

	Rollback() error
}
