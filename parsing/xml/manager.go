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

package xml

import (
	"github.com/xfali/gobatis/errors"
	"github.com/xfali/gobatis/logging"
	"github.com/xfali/gobatis/parsing"
	"github.com/xfali/gobatis/parsing/sqlparser"
	"sync"
)

type Manager struct {
	sqlMap map[string]*parsing.DynamicData
	lock   sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		sqlMap: map[string]*parsing.DynamicData{},
	}
}

func (m *Manager) RegisterData(data []byte) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	mapper, err := Parse(data)
	if err != nil {
		logging.Warn("register mapper data failed: %s err: %v\n", string(data), err)
		return err
	}

	return m.formatMapper(mapper)
}

func (m *Manager) RegisterFile(file string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	mapper, err := ParseFile(file)
	if err != nil {
		logging.Warn("register mapper file failed: %s err: %v\n", file, err)
		return err
	}

	return m.formatMapper(mapper)
}

func (m *Manager) formatMapper(mapper *Mapper) error {
	ret := mapper.Format()
	for k, v := range ret {
		if _, ok := m.sqlMap[k]; ok {
			return errors.SQL_ID_DUPLICATES
		} else {
			m.sqlMap[k] = v
		}
	}
	return nil
}

func (m *Manager) FindSqlParser(sqlId string) (sqlparser.SqlParser, bool) {
	m.lock.Lock()
	defer m.lock.Unlock()

	v, ok := m.sqlMap[sqlId]
	return v, ok
}

func (m *Manager) RegisterSql(sqlId string, sql string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.sqlMap[sqlId]; ok {
		return errors.SQL_ID_DUPLICATES
	} else {
		dd := &parsing.DynamicData{OriginData: sql}
		m.sqlMap[sqlId] = dd
	}
	return nil
}

func (m *Manager) UnregisterSql(sqlId string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	delete(m.sqlMap, sqlId)
}
