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

type Base[T any] interface {
	Insert(entity T) int64

	InsertBatch(entities ...T) (int64, int64)

	DeleteById(id any) int64

	DeleteBatchIds(ids []any) int64

	UpdateById(entity T) int64

	SelectById(id any) T

	SelectBatchIds(ids []any) []T

	SelectOne(entity T) T

	SelectCount(entity T) int64

	SelectList(queryWrapper QueryWrapper[T]) ([]T, error)
}
