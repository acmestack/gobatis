// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description: 

package template

import (
	"github.com/xfali/gobatis/logging"
	"github.com/xfali/gobatis/parsing/sqlparser"
	"html/template"
	"io/ioutil"
	"strings"
	"sync"
)

type Parser struct {
	tpl *template.Template
}

//only use first param
func (p *Parser) ParseMetadata(driverName string, params ...interface{}) (*sqlparser.Metadata, error) {
	b := strings.Builder{}
	var param interface{} = nil
	if len(params) > 0 {
		param = params[0]
	}

	err := p.tpl.Execute(&b, param)
	if err != nil {
		return nil, err
	}

	ret := &sqlparser.Metadata{}
	sql := strings.TrimSpace(b.String())
	action := sql[:6]
	action = strings.ToLower(action)
	ret.Action = action
	ret.PrepareSql = sql
	ret.Params = nil

	return ret, nil
}

func updateSet(sets ... string) string {
	b := strings.Builder{}
	for _, v := range sets {
		if len(v) > 0 {
			b.WriteString(strings.TrimSpace(v))
			b.WriteString(",")
		}
	}
	setStr := b.String()
	if len(setStr) == 0 {
		return ""
	} else {
		return " SET " + setStr[:len(setStr)-1]
	}
}

type Manager struct {
	sqlMap map[string]*Parser
	lock   sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		sqlMap: map[string]*Parser{},
	}
}

func (m *Manager) RegisterData(data []byte) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	tpl := template.New("")
	tpl = tpl.Funcs(template.FuncMap{
		"updateSet": updateSet,
	})
	tpl, err := tpl.Parse(string(data))
	if err != nil {
		logging.Warn("register template data failed: %s err: %v\n", string(data), err)
		return err
	}

	tpls := tpl.Templates()
	for _, v := range tpls {
		m.sqlMap[v.Name()] = &Parser{v}
	}

	return nil
}

func (m *Manager) RegisterFile(file string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	tpl := template.New("")
	tpl = tpl.Funcs(template.FuncMap{
		"updateSet": updateSet,
	})
	data, err := ioutil.ReadFile(file)
	if err != nil {
		logging.Warn("register template file failed: %s err: %v\n", file, err)
		return err
	}
	tpl, err = tpl.Parse(string(data))
	if err != nil {
		logging.Warn("register template file failed: %s err: %v\n", file, err)
		return err
	}

	tpls := tpl.Templates()
	for _, v := range tpls {
		m.sqlMap[v.Name()] = &Parser{v}
	}

	return nil
}

func (m *Manager) FindSqlParser(sqlId string) (*Parser, bool) {
	m.lock.Lock()
	defer m.lock.Unlock()

	v, ok := m.sqlMap[sqlId]
	return v, ok
}
