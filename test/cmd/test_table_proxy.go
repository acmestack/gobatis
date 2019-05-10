package test_package

import (
    "context"
    "github.com/xfali/gobatis"
    "github.com/xfali/gobatis/test/cmd/xml"
)

type TestTableCallProxy gobatis.Session

func init() {
    modelV := TestTable{}
    gobatis.RegisterModel(&modelV)
    gobatis.RegisterMapperData([]byte(test_package.CMD_TEST_XML))
}

func New(proxyMrg *gobatis.SessionManager) *TestTableCallProxy {
    return (*TestTableCallProxy)(proxyMrg.NewSession())
}

func (proxy *TestTableCallProxy) Tx(txFunc func(s *TestTableCallProxy) bool) {
    sess := (*gobatis.Session)(proxy)
    sess.Tx(func(session *gobatis.Session) bool {
        return txFunc(proxy)
    })
}

func (proxy *TestTableCallProxy) SelectTestTable(model TestTable) []TestTable {
    var dataList []TestTable
    (*gobatis.Session)(proxy).Select("selectTestTable").Context(context.Background()).Param(model).Result(&dataList)
    return dataList
}

func (proxy *TestTableCallProxy) SelectTestTableCount(model TestTable) int64 {
    var ret int64
    (*gobatis.Session)(proxy).Select("selectTestTableCount").Context(context.Background()).Param(model).Result(&ret)
    return ret
}

func (proxy *TestTableCallProxy) InsertTestTable(model TestTable) (int64, int64) {
    var ret int64
    runner := (*gobatis.Session)(proxy).Insert("insertTestTable").Context(context.Background()).Param(model)
    runner.Result(&ret)
    id := runner.LastInsertId()
    return ret, id
}

func (proxy *TestTableCallProxy) UpdateTestTable(model TestTable) int64 {
    var ret int64
    (*gobatis.Session)(proxy).Update("updateTestTable").Context(context.Background()).Param(model).Result(&ret)
    return ret
}

func (proxy *TestTableCallProxy) DeleteTestTable(model TestTable) int64 {
    var ret int64
    (*gobatis.Session)(proxy).Delete("deleteTestTable").Context(context.Background()).Param(model).Result(&ret)
    return ret
}
