package impl

import (
	"encoding/json"
	"fmt"
	"github.com/acmestack/gobatis"
	"github.com/acmestack/gobatis/datasource"
	"github.com/acmestack/gobatis/factory"
	"github.com/acmestack/gobatis/plus/mapper"
	"testing"
)

func TestUserMapperImpl_SelectList(t *testing.T) {
	mgr := gobatis.NewSessionManager(connect())
	userMapper := UserMapperImpl[TestTable]{mapper.BaseMapper[TestTable]{SessMgr: mgr}}
	queryWrapper := &mapper.QueryWrapper[TestTable]{}
	queryWrapper.Eq("username", "user1").Or().Eq("username", "user2").Select("username")
	list, _ := userMapper.SelectList(queryWrapper)
	marshal, _ := json.Marshal(list)
	fmt.Println(string(marshal))
}

func connect() factory.Factory {
	return gobatis.NewFactory(
		gobatis.SetMaxConn(100),
		gobatis.SetMaxIdleConn(50),
		gobatis.SetDataSource(&datasource.MysqlDataSource{
			Host:     "localhost",
			Port:     3306,
			DBName:   "test",
			Username: "root",
			Password: "123456",
			Charset:  "utf8",
		}))
}

type TestTable struct {
	TableName gobatis.TableName `test_table`
	Id        int64             `column:"id"`
	Username  string            `column:"username"`
	Password  string            `column:"password"`
}
