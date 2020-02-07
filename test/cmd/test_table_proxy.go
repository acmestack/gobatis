package test_package

import (
    "github.com/xfali/gobatis"
)

func init() {
    modelV := TestTable{}
    gobatis.RegisterModel(&modelV)
    //gobatis.RegisterMapperFile("./xml/test_table_mapper.xml")
}

func SelectTestTable(sess *gobatis.Session, model TestTable) ([]TestTable, error) {
    var dataList []TestTable
    err := sess.Select("selectTestTable").Param(model).Result(&dataList)
    return dataList, err
}

func SelectTestTableCount(sess *gobatis.Session, model TestTable) (int64, error) {
    var ret int64
    err := sess.Select("selectTestTableCount").Param(model).Result(&ret)
    return ret, err
}

func InsertTestTable(sess *gobatis.Session, model TestTable) (int64, int64, error) {
    var ret int64
    runner := sess.Insert("insertTestTable").Param(model)
    err := runner.Result(&ret)
    id := runner.LastInsertId()
    return ret, id, err
}

func UpdateTestTable(sess *gobatis.Session, model TestTable) (int64, error) {
    var ret int64
    err := sess.Update("updateTestTable").Param(model).Result(&ret)
    return ret, err
}

func DeleteTestTable(sess *gobatis.Session, model TestTable) (int64, error) {
    var ret int64
    err := sess.Delete("deleteTestTable").Param(model).Result(&ret)
    return ret, err
}

