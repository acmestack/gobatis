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
    "github.com/xfali/GoBatis/errors"
    "github.com/xfali/GoBatis/logging"
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

func (s *MysqlStatement) Query(params ...interface{}) error {
    stmt := (*sql.Stmt)(s)
    rows, err := stmt.Query(params)
    if err != nil {
        return errors.STATEMENT_QUERY_ERROR
    }
    defer rows.Close()

    columns, _ := rows.Columns()
    scanArgs := make([]interface{}, len(columns))
    values := make([][]byte, len(columns))

    for j := range values {
        scanArgs[j] = &values[j]
    }
    for rows.Next() {
        if err := rows.Scan(scanArgs); err == nil {
            for _, col := range values {
                logging.Debug("%v", col)
            }
        }
    }
    return nil
}

func (s *MysqlStatement) Exec(params ...interface{}) (int64, error) {
    stmt := (*sql.Stmt)(s)
    result, err := stmt.Exec(params)
    if err != nil {
        return 0, errors.STATEMENT_QUERY_ERROR
    }
    ret, err := result.RowsAffected()
    if err != nil {
        return 0, errors.STATEMENT_QUERY_ERROR
    }
    return ret, nil
}

func (s *MysqlStatement) row2Slice(rows *sql.Rows, fields []string, bean interface{}) ([]interface{}, error) {
    scanResults := make([]interface{}, len(fields))
    for i := 0; i < len(fields); i++ {
        var cell interface{}
        scanResults[i] = &cell
    }
    if err := rows.Scan(scanResults...); err != nil {
        return nil, err
    }

    return scanResults, nil
}
