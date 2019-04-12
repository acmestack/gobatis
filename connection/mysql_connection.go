/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package connection

import (
    "database/sql"
    "github.com/xfali/gobatis/errors"
    "github.com/xfali/gobatis/handler"
    "github.com/xfali/gobatis/logging"
)

type MysqlConnection sql.DB
type MysqlStatement sql.Stmt

func (c *MysqlConnection) Prepare(sqlStr string) (Statement, error) {
    db := (*sql.DB)(c)
    s, err := db.Prepare(sqlStr)
    if err != nil {
        return nil, errors.CONNECTION_PREPARE_ERROR
    }
    return (*MysqlStatement)(s), nil
}

func (s *MysqlStatement) Query(handler handler.ResultHandler, iterFunc IterFunc, params ...interface{}) error {
    stmt := (*sql.Stmt)(s)
    rows, err := stmt.Query(params...)
    if err != nil {
        return errors.STATEMENT_QUERY_ERROR
    }
    defer rows.Close()

    columns, _ := rows.Columns()
    scanArgs := make([]interface{}, len(columns))
    values := make([]interface{}, len(columns))

    for j := range values {
        scanArgs[j] = &values[j]
    }
    var index int64 = 0
    for rows.Next() {
        if err := rows.Scan(scanArgs); err == nil {
            for _, col := range values {
                logging.Debug("%v", col)
            }
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
    return nil
}

func (s *MysqlStatement) Exec(params ...interface{}) (int64, error) {
    stmt := (*sql.Stmt)(s)
    result, err := stmt.Exec(params...)
    if err != nil {
        return 0, errors.STATEMENT_QUERY_ERROR
    }
    ret, err := result.RowsAffected()
    if err != nil {
        return 0, errors.STATEMENT_QUERY_ERROR
    }
    return ret, nil
}
