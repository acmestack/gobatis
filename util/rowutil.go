/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package util

import (
    "database/sql"
    "github.com/xfali/gobatis"
    "github.com/xfali/gobatis/handler"
)

func ScanRows(rows *sql.Rows, handler handler.ResultHandler, iterFunc gobatis.IterFunc) int64 {
    columns, _ := rows.Columns()
    scanArgs := make([]interface{}, len(columns))
    values := make([]interface{}, len(columns))

    for j := range values {
        scanArgs[j] = &values[j]
    }
    var index int64 = 0
    for rows.Next() {
        if err := rows.Scan(scanArgs...); err == nil {
            //for _, col := range values {
            //    logging.Debug("%v", col)
            //}
            result, err := handler.Deserialize(columns, values)
            if err == nil {
                stop := iterFunc(index, result)
                if stop {
                    break
                }
            }
            index++
        }
    }
    return index
}
