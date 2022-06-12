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
	"errors"
	_ "github.com/lib/pq"
	"github.com/xfali/gobatis"
	"github.com/xfali/gobatis/datasource"
	"testing"
)

var sessionMgr *gobatis.SessionManager

func init() {
	fac := gobatis.NewFactory(
		gobatis.SetMaxConn(100),
		gobatis.SetMaxIdleConn(50),
		gobatis.SetDataSource(&datasource.PostgreDataSource{
			Host:     "localhost",
			Port:     5432,
			DBName:   "testdb",
			Username: "test",
			Password: "test",
		}))
	sessionMgr = gobatis.NewSessionManager(fac)
}

func TestSelect(t *testing.T) {
	ret, err := SelectTestTable(sessionMgr.NewSession(), TestTable{Username: "test_user"})

	if err != nil {
		t.Fatal(err)
	}

	for _, v := range ret {
		t.Logf("value: %v", v)
	}
}

func TestSelectCount(t *testing.T) {
	//all count
	ret, err := SelectTestTableCount(sessionMgr.NewSession(), TestTable{})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("count is %d", ret)
}

func TestInsert(t *testing.T) {
	ret, id, err := InsertTestTable(sessionMgr.NewSession(), TestTable{Username: "test_insert_user", Password: "test_pw"})
	if err != nil {
		t.Log(err)
	}
	t.Logf("insert ret is %d, id is %d", ret, id)
}

func TestInsertBatch(t *testing.T) {
	ret, id, err := InsertBatchTestTable(sessionMgr.NewSession(), []TestTable{
		{Username: "test_insert_user4", Password: "test_pw4"},
		{Username: "test_insert_user5", Password: "test_pw5"},
	})
	if err != nil {
		t.Log(err)
	}
	t.Logf("insert ret is %d, id is %d", ret, id)
}

func TestUpdate(t *testing.T) {
	ret, err := UpdateTestTable(sessionMgr.NewSession(), TestTable{Id: 1, Username: "test_insert_user", Password: "update_pw"})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("update ret is %d", ret)
}

func TestDelete(t *testing.T) {
	//Same as: DELETE FROM test_table WHERE username = 'test_insert_user'
	ret, err := DeleteTestTable(sessionMgr.NewSession(), TestTable{Username: "test_insert_user"})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("delete ret is %d", ret)
}

func TestSessionTpl(t *testing.T) {
	gobatis.RegisterTemplateFile("./template/test_table_mapper.tmpl")
	var param = TestTable{Id: 1, Username: "user", Password: "pw"}
	t.Run("select", func(t *testing.T) {
		sess := sessionMgr.NewSession()
		ret, _ := SelectTestTable(sess, TestTable{Id: 1})
		t.Log(ret)
	})

	t.Run("count", func(t *testing.T) {
		sess := sessionMgr.NewSession()
		ret, _ := SelectTestTableCount(sess, TestTable{Id: 1})
		t.Log(ret)
	})

	t.Run("insert", func(t *testing.T) {
		sess := sessionMgr.NewSession()
		ret, _, err := InsertTestTable(sess, param)
		t.Log(err)
		t.Log(ret)
	})

	t.Run("insertBatch", func(t *testing.T) {
		sess := sessionMgr.NewSession()
		ret, _, err := InsertBatchTestTable(sess, []TestTable{
			{Username: "test_insert_user14", Password: "test_pw14"},
			{Username: "test_insert_user15", Password: "test_pw15"},
		})
		t.Log(err)
		t.Log(ret)
	})

	t.Run("update", func(t *testing.T) {
		sess := sessionMgr.NewSession()
		ret, err := UpdateTestTable(sess, TestTable{Id: 1, Username: "user2", Password: "pw2"})
		t.Log(err)
		t.Log(ret)
	})

	t.Run("delete", func(t *testing.T) {
		sess := sessionMgr.NewSession()
		ret, err := DeleteTestTable(sess, TestTable{Id: 1})
		t.Log(ret)
		t.Log(err)
	})
}

func TestSessionXml(t *testing.T) {
	gobatis.RegisterMapperFile("./xml/test_table_mapper.xml")
	var param = TestTable{Id: 1, Username: "user", Password: "pw"}
	t.Run("select", func(t *testing.T) {
		sess := sessionMgr.NewSession()
		ret, _ := SelectTestTable(sess, TestTable{Id: 1})
		t.Log(ret)
	})

	t.Run("count", func(t *testing.T) {
		sess := sessionMgr.NewSession()
		ret, _ := SelectTestTableCount(sess, TestTable{Id: 1})
		t.Log(ret)
	})

	t.Run("insert", func(t *testing.T) {
		sess := sessionMgr.NewSession()
		ret, _, err := InsertTestTable(sess, param)
		t.Log(err)
		t.Log(ret)
	})

	t.Run("insertBatch", func(t *testing.T) {
		sess := sessionMgr.NewSession()
		ret, _, err := InsertBatchTestTable(sess, []TestTable{
			{Username: "test_insert_user14", Password: "test_pw14"},
			{Username: "test_insert_user15", Password: "test_pw15"},
		})
		t.Log(err)
		t.Log(ret)
	})

	t.Run("update", func(t *testing.T) {
		sess := sessionMgr.NewSession()
		ret, err := UpdateTestTable(sess, TestTable{Id: 1, Username: "user2", Password: "pw2"})
		t.Log(err)
		t.Log(ret)
	})

	t.Run("delete", func(t *testing.T) {
		sess := sessionMgr.NewSession()
		ret, err := DeleteTestTable(sess, TestTable{Id: 1})
		t.Log(ret)
		t.Log(err)
	})
}

func TestTxSuccess(t *testing.T) {
	sessionMgr.NewSession().Tx(func(s *gobatis.Session) error {
		_, _, err := InsertTestTable(s, TestTable{Username: "tx_user", Password: "tx_pw"})
		t.Log(err)
		//select all
		ret, err := SelectTestTable(s, TestTable{})
		if err != nil {
			t.Fatal(err)
		}
		for _, v := range ret {
			t.Logf("value: %v", v)
		}
		//will commit
		return nil
	})

	//select all
	ret, _ := SelectTestTable(sessionMgr.NewSession(), TestTable{Username: "tx_user"})

	for _, v := range ret {
		t.Logf("value out tx: %v", v)
	}
}

func TestTxFail(t *testing.T) {
	sessionMgr.NewSession().Tx(func(s *gobatis.Session) error {
		_, _, err := InsertTestTable(s, TestTable{Username: "tx_user", Password: "tx_pw"})
		t.Log(err)
		////select all
		//ret, err := SelectTestTable(s, TestTable{})
		//if err != nil {
		//	t.Fatal(err)
		//}
		//for _, v := range ret {
		//	t.Logf("value: %v", v)
		//}
		//will commit
		return errors.New("rollback")
	})

	//select all
	ret, _ := SelectTestTable(sessionMgr.NewSession(), TestTable{Username: "tx_user"})

	for _, v := range ret {
		t.Logf("value out tx: %v\n", v)
	}
}
