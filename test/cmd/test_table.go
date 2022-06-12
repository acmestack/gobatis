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

package test

import (
	"github.com/xfali/gobatis"
	"time"
)

type TestTable struct {
	//TableName gobatis.ModelName `test_table`
	Id         int       `xfield:"id"`
	Username   string    `xfield:"username"`
	Password   string    `xfield:"password"`
	Createtime time.Time `xfield:"createtime"`
}

func (m *TestTable) Select(sess *gobatis.Session) ([]TestTable, error) {
	return SelectTestTable(sess, *m)
}

func (m *TestTable) Count(sess *gobatis.Session) (int64, error) {
	return SelectTestTableCount(sess, *m)
}

func (m *TestTable) Insert(sess *gobatis.Session) (int64, int64, error) {
	return InsertTestTable(sess, *m)
}

func (m *TestTable) Update(sess *gobatis.Session) (int64, error) {
	return UpdateTestTable(sess, *m)
}

func (m *TestTable) Delete(sess *gobatis.Session) (int64, error) {
	return DeleteTestTable(sess, *m)
}
