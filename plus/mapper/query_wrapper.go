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
	"github.com/acmestack/gobatis/builder"
	"github.com/acmestack/gobatis/plus/constants"
	"reflect"
)

type QueryWrapper[T any] struct {
	Columns           []string
	Entity            *T
	SqlBuild          *builder.SQLFragment
	TableName         string
	Expression        []any
	ParamNameSeq      int
	LastConditionType string
}

func (queryWrapper *QueryWrapper[T]) Eq(column string, val any) Wrapper[T] {
	queryWrapper.setCondition(column, val, constants.Eq)
	return queryWrapper
}

func (queryWrapper *QueryWrapper[T]) Ne(column string, val any) Wrapper[T] {
	queryWrapper.setCondition(column, val, constants.Ne)
	return queryWrapper
}

func (queryWrapper *QueryWrapper[T]) Gt(column string, val any) Wrapper[T] {
	queryWrapper.setCondition(column, val, constants.Gt)
	return queryWrapper
}

func (queryWrapper *QueryWrapper[T]) Ge(column string, val any) Wrapper[T] {
	queryWrapper.setCondition(column, val, constants.Ge)
	return queryWrapper
}

func (queryWrapper *QueryWrapper[T]) Lt(column string, val any) Wrapper[T] {
	queryWrapper.setCondition(column, val, constants.Lt)
	return queryWrapper
}

func (queryWrapper *QueryWrapper[T]) Le(column string, val any) Wrapper[T] {
	queryWrapper.setCondition(column, val, constants.Le)
	return queryWrapper
}

func (queryWrapper *QueryWrapper[T]) Like(column string, val any) Wrapper[T] {
	s := val.(string)
	queryWrapper.setCondition(column, "%"+s+"%", constants.Like)
	return queryWrapper
}

func (queryWrapper *QueryWrapper[T]) NotLike(column string, val any) Wrapper[T] {
	s := val.(string)
	queryWrapper.setCondition(column, "%"+s+"%", constants.Not+constants.Like)
	return queryWrapper
}

func (queryWrapper *QueryWrapper[T]) LikeLeft(column string, val any) Wrapper[T] {
	s := val.(string)
	queryWrapper.setCondition(column, "%"+s, constants.Like)
	return queryWrapper
}

func (queryWrapper *QueryWrapper[T]) LikeRight(column string, val any) Wrapper[T] {
	s := val.(string)
	queryWrapper.setCondition(column, s+"%", constants.Like)
	return queryWrapper
}

func (queryWrapper *QueryWrapper[T]) And() Wrapper[T] {
	queryWrapper.Expression = append(queryWrapper.Expression, constants.Eq)
	queryWrapper.LastConditionType = constants.Eq
	return queryWrapper
}

func (queryWrapper *QueryWrapper[T]) Or() Wrapper[T] {
	queryWrapper.Expression = append(queryWrapper.Expression, constants.Or)
	queryWrapper.LastConditionType = constants.Or
	return queryWrapper
}

func (queryWrapper *QueryWrapper[T]) Select(columns ...string) Wrapper[T] {
	queryWrapper.SqlBuild.Select(columns...)
	return queryWrapper
}

func (queryWrapper *QueryWrapper[T]) init() {
	if queryWrapper.Entity == nil {
		queryWrapper.Entity = new(T)
	}
	if queryWrapper.TableName == "" {
		queryWrapper.setTableName()
	}
}

type ParamValue struct {
	value any
}

func (queryWrapper *QueryWrapper[T]) setCondition(column string, val any, conditionType string) {

	if queryWrapper.LastConditionType != constants.And && queryWrapper.LastConditionType != constants.Or && len(queryWrapper.Expression) > 0 {
		queryWrapper.Expression = append(queryWrapper.Expression, constants.And)
	}

	queryWrapper.Expression = append(queryWrapper.Expression, column)

	queryWrapper.Expression = append(queryWrapper.Expression, conditionType)

	queryWrapper.Expression = append(queryWrapper.Expression, ParamValue{val})

}

func setField(entityValueRef reflect.Value, field reflect.StructField, val any) {
	ft := field.Type
	switch ft.Kind() {
	case reflect.String:
		entityValueRef.FieldByName(field.Name).SetString(val.(string))
	case reflect.Int:
		i := val.(int)
		entityValueRef.FieldByName(field.Name).SetInt(int64(i))
	}
}

func (queryWrapper *QueryWrapper[T]) setTableName() {
	// todo The future is through annotations get the tableName
	entityRef := reflect.TypeOf(queryWrapper.Entity).Elem()
	tableName := entityRef.Field(0).Tag
	queryWrapper.TableName = string(tableName)

	queryWrapper.checkColumns()

	queryWrapper.SqlBuild = builder.Select(queryWrapper.Columns...).From(string(tableName))
}

func (queryWrapper *QueryWrapper[T]) checkColumns() {
	if len(queryWrapper.Columns) == 0 {
		entityRef := reflect.TypeOf(queryWrapper.Entity).Elem()
		numField := entityRef.NumField()
		for i := 0; i < numField; i++ {
			field := entityRef.Field(i)
			filedName := field.Tag.Get("xfield")
			if filedName != "" {
				queryWrapper.Columns = append(queryWrapper.Columns, filedName)
			}
		}
	}
}
