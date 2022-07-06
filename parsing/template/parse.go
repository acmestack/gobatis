/*
 * Copyright (c) 2022, AcmeStack
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

package template

import (
	"io/ioutil"
	"strings"
	"sync"
	"text/template"

	"github.com/acmestack/gobatis/errors"
	"github.com/acmestack/gobatis/logging"
	"github.com/acmestack/gobatis/parsing/sqlparser"
)

const (
	namespaceTmplName = "namespace"
)

type Parser struct {
	//template
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

// ParseMetadata only use first param
func (p *Parser) ParseMetadata(driverName string, params ...interface{}) (*sqlparser.Metadata, error) {
	if p.tpl == nil {
		return nil, errors.ParseTemplateNilError
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

func (manager *Manager) RegisterData(data []byte) error {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	tpl := template.New("")
	tpl = tpl.Funcs(dummyFuncMap)
	tpl, err := tpl.Parse(string(data))
	if err != nil {
		logging.Warn("register template data failed: %s err: %v\n", string(data), err)
		return err
	}

	ns := getNamespace(tpl)
	tpls := tpl.Templates()
	for _, v := range tpls {
		if v.Name() != "" && v.Name() != namespaceTmplName {
			manager.sqlMap[ns+v.Name()] = &Parser{tpl: v}
		}
	}

	return nil
}

func (manager *Manager) RegisterFile(file string) error {
	manager.lock.Lock()
	defer manager.lock.Unlock()

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

	ns := getNamespace(tpl)
	tpls := tpl.Templates()
	for _, v := range tpls {
		if v.Name() != "" && v.Name() != namespaceTmplName {
			manager.sqlMap[ns+v.Name()] = &Parser{tpl: v}
		}
	}

	return nil
}

func getNamespace(tpl *template.Template) string {
	ns := strings.Builder{}
	nsTpl := tpl.Lookup(namespaceTmplName)
	if nsTpl != nil {
		err := nsTpl.Execute(&ns, nil)
		if err != nil {
			ns.Reset()
		}
	}

	ret := strings.TrimSpace(ns.String())

	if ret != "" {
		ret = ret + "."
	}
	return ret
}

func (manager *Manager) FindSqlParser(sqlId string) (*Parser, bool) {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	v, ok := manager.sqlMap[sqlId]
	return v, ok
}
