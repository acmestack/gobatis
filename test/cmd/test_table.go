package test_package

import "time"

type TestTable struct {
    //TableName gobatis.ModelName `test_table`
    Id int `xfield:"id"`
    Username string `xfield:"username"`
    Password string `xfield:"password"`
    Createtime time.Time `xfield:"createtime"`
}
