package test_package

import "time"

type TestTable struct {
    //TableName gobatis.ModelName `test_table`
    Id         int64     `xfield:"id"`
    Username   string    `xfield:"username"`
    Password   string    `xfield:"password"`
    UpdateTime time.Time `xfield:"update_time"`
}
