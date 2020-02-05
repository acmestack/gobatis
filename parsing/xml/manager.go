// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description: 

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
