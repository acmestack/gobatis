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
	queryWrapper.Eq("username", "user123").Select("username")
	list, err := userMapper.SelectList(queryWrapper)
	if err != nil {
		fmt.Println(err.Error())
	}
	marshal, _ := json.Marshal(list)
	fmt.Println(string(marshal))
}

func TestUserMapperImpl_SelectOne(t *testing.T) {
	mgr := gobatis.NewSessionManager(connect())
	userMapper := UserMapperImpl[TestTable]{mapper.BaseMapper[TestTable]{SessMgr: mgr}}
	queryWrapper := &mapper.QueryWrapper[TestTable]{}
	queryWrapper.Eq("username", "user1").Select("username", "password")
	entity, err := userMapper.SelectOne(queryWrapper)
	if err != nil {
		fmt.Println(err.Error())
	}
	marshal, _ := json.Marshal(entity)
	fmt.Println(string(marshal))
}

func TestUserMapperImpl_SelectCount(t *testing.T) {
	mgr := gobatis.NewSessionManager(connect())
	userMapper := UserMapperImpl[TestTable]{mapper.BaseMapper[TestTable]{SessMgr: mgr}}
	queryWrapper := &mapper.QueryWrapper[TestTable]{}
	queryWrapper.Eq("username", "user123")
	count, err := userMapper.SelectCount(queryWrapper)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(count)
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
