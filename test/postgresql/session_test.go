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

package test

import (
	"github.com/acmestack/gobatis"
	"testing"
)

func TestSessionTpl(t *testing.T) {
	gobatis.RegisterTemplateFile("./postgresql.tpl")
	var param = TestTable{Id: 1, Username: "user", Password: "pw"}
	t.Run("select", func(t *testing.T) {
		mgr := gobatis.NewSessionManager(connect())
		sess := mgr.NewSession()
		var ret []TestTable
		sess.Select("selectTestTable").Param(param).Result(&ret)
		t.Log(ret)
	})

	t.Run("insert", func(t *testing.T) {
		mgr := gobatis.NewSessionManager(connect())
		sess := mgr.NewSession()
		var ret int
		err := sess.Insert("insertTestTable").Param(param).Result(&ret)
		t.Log(err)
		t.Log(ret)
	})

	t.Run("insertBatch", func(t *testing.T) {
		mgr := gobatis.NewSessionManager(connect())
		sess := mgr.NewSession()
		var ret int
		err := sess.Insert("insertBatchTestTable").Param([]TestTable{
			{Id: 13, Username: "user13", Password: "pw13"},
			{Id: 14, Username: "user14", Password: "pw14"},
		}).Result(&ret)
		t.Log(err)
		t.Log(ret)
	})

	t.Run("update", func(t *testing.T) {
		mgr := gobatis.NewSessionManager(connect())
		sess := mgr.NewSession()
		var ret int
		err := sess.Update("updateTestTable").Param(TestTable{Id: 1, Username: "user2", Password: "pw2"}).Result(&ret)
		t.Log(err)
		t.Log(ret)
	})

	t.Run("delete", func(t *testing.T) {
		mgr := gobatis.NewSessionManager(connect())
		sess := mgr.NewSession()
		var ret int
		sess.Delete("deleteTestTable").Param(TestTable{Id: 1}).Result(&ret)
		t.Log(ret)
	})
}

func TestSessionXml(t *testing.T) {
	gobatis.RegisterMapperFile("./postgresql.xml")
	var param = TestTable{Id: 1, Username: "user", Password: "pw"}
	t.Run("select", func(t *testing.T) {
		mgr := gobatis.NewSessionManager(connect())
		sess := mgr.NewSession()
		var ret []TestTable
		sess.Select("selectTestTable").Param(param).Result(&ret)
		t.Log(ret)
	})

	t.Run("insert", func(t *testing.T) {
		mgr := gobatis.NewSessionManager(connect())
		sess := mgr.NewSession()
		var ret int
		err := sess.Insert("insertTestTable").Param(param).Result(&ret)
		t.Log(err)
		t.Log(ret)
	})

	t.Run("update", func(t *testing.T) {
		mgr := gobatis.NewSessionManager(connect())
		sess := mgr.NewSession()
		var ret int
		err := sess.Update("updateTestTable").Param(TestTable{Id: 1, Username: "user2", Password: "pw2"}).Result(&ret)
		t.Log(err)
		t.Log(ret)
	})

	t.Run("delete", func(t *testing.T) {
		mgr := gobatis.NewSessionManager(connect())
		sess := mgr.NewSession()
		var ret int
		sess.Delete("deleteTestTable").Param(TestTable{Id: 1}).Result(&ret)
		t.Log(ret)
	})
}
