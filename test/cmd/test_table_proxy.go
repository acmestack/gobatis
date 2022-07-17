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

package test

import (
	"github.com/acmestack/gobatis"
)

func init() {
	modelV := TestTable{}
	gobatis.RegisterModel(&modelV)
	//gobatis.RegisterTemplateFile("./template/test_table_mapper.tmpl")
}

func SelectTestTable(sess *gobatis.Session, model TestTable) ([]TestTable, error) {
	var dataList []TestTable
	err := sess.Select("test.selectTestTable").Param(model).Result(&dataList)
	return dataList, err
}

func SelectTestTableCount(sess *gobatis.Session, model TestTable) (int64, error) {
	var ret int64
	err := sess.Select("test.selectTestTableCount").Param(model).Result(&ret)
	return ret, err
}

func InsertTestTable(sess *gobatis.Session, model TestTable) (int64, int64, error) {
	var ret int64
	runner := sess.Insert("test.insertTestTable").Param(model)
	err := runner.Result(&ret)
	id := runner.LastInsertId()
	return ret, id, err
}

func InsertBatchTestTable(sess *gobatis.Session, models []TestTable) (int64, int64, error) {
	var ret int64
	runner := sess.Insert("test.insertBatchTestTable").Param(models)
	err := runner.Result(&ret)
	id := runner.LastInsertId()
	return ret, id, err
}

func UpdateTestTable(sess *gobatis.Session, model TestTable) (int64, error) {
	var ret int64
	err := sess.Update("test.updateTestTable").Param(model).Result(&ret)
	return ret, err
}

func DeleteTestTable(sess *gobatis.Session, model TestTable) (int64, error) {
	var ret int64
	err := sess.Delete("test.deleteTestTable").Param(model).Result(&ret)
	return ret, err
}
