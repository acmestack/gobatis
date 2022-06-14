/*
 * Copyright (c) 2022, AcmeStack
 * All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package util

import (
	"database/sql"
	"github.com/acmestack/gobatis/reflection"
	"reflect"
)

func ScanRows(rows *sql.Rows, result reflection.Object) int64 {
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
			if !deserialize(result, columns, values) {
				break
			}
			index++
		}
	}
	return index
}

func deserialize(result reflection.Object, columns []string, values []interface{}) bool {
	obj := result
	if result.CanAddValue() {
		obj = result.NewElem()
	}
	for i := range columns {
		if obj.CanSetField() {
			obj.SetField(columns[i], values[i])
		} else {
			obj.SetValue(reflect.ValueOf(values[0]))
			break
		}
	}
	if result.CanAddValue() {
		result.AddValue(obj.GetValue())
		return true
	}
	return false
}
