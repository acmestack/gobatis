// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description: 

package template

import (
	"github.com/xfali/gobatis/errors"
	"github.com/xfali/gobatis/logging"
	"github.com/xfali/gobatis/parsing/sqlparser"
	"io/ioutil"
	"strings"
	"sync"
	"text/template"
)

type Parser struct {
	tpl *template.Template
}

func CreateParser(data []byte) (*Parser, error) {
	tpl := template.New("")
	tpl = tpl.Funcs(dummyFuncMap)
	tpl, err := tpl.Parse(string(data))
	if err != nil {
		return nil, err
	}
	return &Parser{tpl: tpl}, nil
}

//only use first param
func (p *Parser) ParseMetadata(driverName string, params ...interface{}) (*sqlparser.Metadata, error) {
	if p.tpl == nil {
		return nil, errors.PARSE_TEMPLATE_NIL_ERROR
	}
	b := strings.Builder{}
	var param interface{} = nil
	if len(params) == 1 {
		param = params[0]
	} else {
		param = params
	}
	dynamic := selectDynamic(driverName)
	tpl := p.tpl.Funcs(dynamic.getFuncMap())
	err := tpl.Execute(&b, param)
	if err != nil {
		return nil, err
	}

	ret := &sqlparser.Metadata{}
	sql := strings.TrimSpace(b.String())
	action := sql[:6]
	action = strings.ToLower(action)
	ret.Action = action
	ret.PrepareSql, ret.Params = dynamic.format(sql)

	return ret, nil
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
	tpl = tpl.Funcs(dummyFuncMap)
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
	data, err := ioutil.ReadFile(file)
	if err != nil {
		logging.Warn("register template file failed: %s err: %v\n", file, err)
		return err
	}
	tpl = tpl.Funcs(dummyFuncMap)
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
