/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description:
 */

package test

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/xfali/gobatis"
	"github.com/xfali/gobatis/datasource"
	"testing"
)

type SqlTest struct {
	Id       int64  `xfield:"id"`
	Username string `xfield:"username"`
	Password string `xfield:"password"`
}

var sessionMgr *gobatis.SessionManager

func init() {
	fac := gobatis.NewFactory(
		gobatis.SetMaxConn(100),
		gobatis.SetMaxIdleConn(50),
		gobatis.SetDataSource(&datasource.MysqlDataSource{
			Host:     "localhost",
			Port:     3306,
			DBName:   "test",
			Username: "root",
			Password: "123",
			Charset:  "utf8",
		}))
	var testV TestTable
	gobatis.RegisterModel(&testV)
	sessionMgr = gobatis.NewSessionManager(fac)
}

func TestSelectWithSimpleType(t *testing.T) {
	sql := `SELECT username, password FROM test_table WHERE id = #{0}`
	ret := SqlTest{}
	sessionMgr.NewSession().Select(sql).Param(1).Result(&ret)
	t.Logf("%v %v", ret.Username, ret.Password)
}

func TestSelectWithMap(t *testing.T) {
	sql := `SELECT username, password FROM test_table WHERE id = #{id}`
	ret := SqlTest{}
	sessionMgr.NewSession().Select(sql).Param(map[string]interface{}{"id": 1}).Result(&ret)
	t.Logf("%v %v", ret.Username, ret.Password)
}

func TestSelectWithStruct(t *testing.T) {
	sql := `SELECT username, password FROM test_table WHERE id = #{SqlTest.id}`
	ret := SqlTest{}
	sessionMgr.NewSession().Select(sql).Param(SqlTest{Id: 1}).Result(&ret)
	t.Logf("%v %v", ret.Username, ret.Password)
}

func TestSelectWithComplex(t *testing.T) {
	sql := `SELECT username, password FROM test_table WHERE id = #{SqlTest.id} AND username = #{0} AND password = #{pw}`
	ret := SqlTest{}
	sessionMgr.NewSession().Select(sql).Param(SqlTest{Id: 1}, "test_user", map[string]interface{}{"pw": "test_pw"}).Result(&ret)
	t.Logf("%v %v", ret.Username, ret.Password)
}
