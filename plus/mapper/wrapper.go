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

type Wrapper[T any] interface {
	Eq(column string, val any) Wrapper[T]

	Ne(column string, val any) Wrapper[T]

	Gt(column string, val any) Wrapper[T]

	Ge(column string, val any) Wrapper[T]

	Lt(column string, val any) Wrapper[T]

	Le(column string, val any) Wrapper[T]

	Like(column string, val any) Wrapper[T]

	NotLike(column string, val any) Wrapper[T]

	LikeLeft(column string, val any) Wrapper[T]

	LikeRight(column string, val1 any) Wrapper[T]

	And() Wrapper[T]

	Or() Wrapper[T]

	Select(columns ...string) Wrapper[T]
}
