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
	"strconv"
	"strings"
	"time"
)

type BaseMapper[T any] struct {
	SessMgr      *gobatis.SessionManager
	Ctx          context.Context
	Columns      []string
	ParamNameSeq int32
	ParamMap     map[string]any
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
		queryWrapper.init()
	}
	if userMapper.ParamMap == nil {
		userMapper.ParamMap = map[string]any{}
	}
	expression := queryWrapper.Expression
	build := strings.Builder{}
	for _, v := range expression {
		if paramValue, ok := v.(ParamValue); ok {
			queryWrapper.ParamNameSeq = queryWrapper.ParamNameSeq + 1
			mapping := constants.MAPPING + strconv.Itoa(queryWrapper.ParamNameSeq)
			userMapper.ParamMap[mapping] = paramValue.value
			build.WriteString(constants.HASH_LEFT_BRACE + mapping + constants.RIGHT_BRACE + constants.SPACE)
		} else {
			build.WriteString(v.(string) + constants.SPACE)
		}
	}
	sqlCondition := build.String()
	sqlId := constants.SELECT + constants.CONNECTION + strconv.Itoa(time.Now().Nanosecond())
	sql := constants.SELECT + constants.SPACE + constants.ASTERISK + constants.SPACE + constants.FROM + constants.SPACE + "test_table" +
		constants.SPACE + constants.WHERE + constants.SPACE + sqlCondition
	err := gobatis.RegisterSql(sqlId, sql)
	if err != nil {
		return nil, err
	}
	sess := userMapper.SessMgr.NewSession()
	var arr []T
	err = sess.Select(sqlId).Param(userMapper.ParamMap).Result(&arr)
	if err != nil {
		return nil, err
	}
	gobatis.UnregisterSql(sqlId)
	return arr, nil
}
