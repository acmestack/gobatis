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

package mapper

import (
	"context"
	"github.com/acmestack/gobatis"
	constants "github.com/acmestack/gobatis/plus/constants"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type BaseMapper[T any] struct {
	SessMgr      *gobatis.SessionManager
	Ctx          context.Context
	Columns      []string
	ParamNameSeq int32
}

func (userMapper *BaseMapper[T]) Insert(entity T) int64 {
	return 0
}

func (userMapper *BaseMapper[T]) InsertBatch(entities ...T) (int64, int64) {
	return 0, 0
}
func (userMapper *BaseMapper[T]) DeleteById(id any) int64 {
	return 0
}
func (userMapper *BaseMapper[T]) DeleteBatchIds(ids []any) int64 {
	return 0
}
func (userMapper *BaseMapper[T]) UpdateById(entity T) int64 {
	return 0
}
func (userMapper *BaseMapper[T]) SelectById(id any) T {
	return *new(T)
}
func (userMapper *BaseMapper[T]) SelectBatchIds(ids []any) []T {
	var arr []T
	return arr
}
func (userMapper *BaseMapper[T]) SelectOne(entity T) T {
	return *new(T)
}
func (userMapper *BaseMapper[T]) SelectCount(entity T) int64 {
	return 0
}

func (userMapper *BaseMapper[T]) SelectList(queryWrapper *QueryWrapper[T]) ([]T, error) {
	if queryWrapper == nil {
		queryWrapper = &QueryWrapper[T]{}
	}

	sqlCondition, paramMap := userMapper.buildCondition(queryWrapper)

	sqlId, sql := userMapper.buildSelectSql(queryWrapper, sqlCondition)

	err := gobatis.RegisterSql(sqlId, sql)
	if err != nil {
		return nil, err
	}

	sess := userMapper.SessMgr.NewSession()
	var arr []T
	err = sess.Select(sqlId).Param(paramMap).Result(&arr)
	if err != nil {
		return nil, err
	}

	// delete sqlId
	gobatis.UnregisterSql(sqlId)
	return arr, nil
}

func (userMapper *BaseMapper[T]) buildCondition(queryWrapper *QueryWrapper[T]) (string, map[string]any) {
	var paramMap = map[string]any{}
	expression := queryWrapper.Expression
	build := strings.Builder{}
	for _, v := range expression {
		if paramValue, ok := v.(ParamValue); ok {
			queryWrapper.ParamNameSeq = queryWrapper.ParamNameSeq + 1
			mapping := constants.MAPPING + strconv.Itoa(queryWrapper.ParamNameSeq)
			paramMap[mapping] = paramValue.value
			build.WriteString(constants.HASH_LEFT_BRACE + mapping + constants.RIGHT_BRACE + constants.SPACE)
		} else {
			build.WriteString(v.(string) + constants.SPACE)
		}
	}
	return build.String(), paramMap
}

func (userMapper *BaseMapper[T]) buildSelectSql(queryWrapper *QueryWrapper[T], sqlCondition string) (string, string) {

	tableName := userMapper.getTableName()

	sqlId := buildSqlId(constants.SELECT)

	var sqlFirstPart string
	if len(queryWrapper.Columns) > 0 {
		columns := strings.Join(queryWrapper.Columns, ",")
		// For example: select username,password from table
		sqlFirstPart = buildSelectSqlFirstPart(columns, tableName)
	} else {
		// For example: select * from table
		sqlFirstPart = buildSelectSqlFirstPart(constants.ASTERISK, tableName)
	}

	var sql string
	if len(queryWrapper.Expression) > 0 {
		sql = sqlFirstPart + constants.SPACE + constants.WHERE + constants.SPACE + sqlCondition
	} else {
		sql = sqlFirstPart
	}

	return sqlId, sql
}

func (userMapper *BaseMapper[T]) getTableName() string {
	entityRef := reflect.TypeOf(new(T)).Elem()
	tableNameTag := entityRef.Field(0).Tag
	tableName := string(tableNameTag)
	return tableName
}

func buildSqlId(sqlType string) string {
	sqlId := sqlType + constants.CONNECTION + strconv.Itoa(time.Now().Nanosecond())
	return sqlId
}

func buildSelectSqlFirstPart(columns string, tableName string) string {
	return constants.SELECT + constants.SPACE + columns + constants.SPACE + constants.FROM + constants.SPACE + tableName
}
