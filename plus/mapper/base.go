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
	Save(entity T) int64

	SaveBatch(entities ...T) (int64, int64)

	UpdateById(entity T) int64

	SelectById(id any) (T, error)

	SelectBatchIds(queryWrapper *QueryWrapper[T]) ([]T, error)

	SelectOne(queryWrapper *QueryWrapper[T]) (T, error)

	SelectCount(queryWrapper *QueryWrapper[T]) (int64, error)

	SelectList(queryWrapper *QueryWrapper[T]) ([]T, error)

	DeleteById(id any) int64

	DeleteBatchIds(ids []any) int64
}
