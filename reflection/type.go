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

package reflection

import (
	"reflect"
	"time"
)

var (
	EmptyString       string
	BoolDefault       bool
	ByteDefault       byte
	Complex64Default  complex64
	Complex128Default complex128
	Float32Default    float32
	Float64Default    float64
	Int64Default      int64
	Uint64Default     uint64
	Int32Default      int32
	Uint32Default     uint32
	Int16Default      int16
	Uint16Default     uint16
	Int8Default       int8
	Uint8Default      uint8
	IntDefault        int
	UintDefault       uint
	TimeDefault       time.Time
)

var (
	IntType   = reflect.TypeOf(IntDefault)
	Int8Type  = reflect.TypeOf(Int8Default)
	Int16Type = reflect.TypeOf(Int16Default)
	Int32Type = reflect.TypeOf(Int32Default)
	Int64Type = reflect.TypeOf(Int64Default)

	UintType   = reflect.TypeOf(UintDefault)
	Uint8Type  = reflect.TypeOf(Uint8Default)
	Uint16Type = reflect.TypeOf(Uint16Default)
	Uint32Type = reflect.TypeOf(Uint32Default)
	Uint64Type = reflect.TypeOf(Uint64Default)

	Float32Type = reflect.TypeOf(Float32Default)
	Float64Type = reflect.TypeOf(Float64Default)

	Complex64Type  = reflect.TypeOf(Complex64Default)
	Complex128Type = reflect.TypeOf(Complex128Default)

	StringType = reflect.TypeOf(EmptyString)
	BoolType   = reflect.TypeOf(BoolDefault)
	ByteType   = reflect.TypeOf(ByteDefault)
	BytesType  = reflect.SliceOf(ByteType)

	TimeType = reflect.TypeOf(TimeDefault)
)

var (
	IntKind   = IntType.Kind()
	Int8Kind  = Int8Type.Kind()
	Int16Kind = Int16Type.Kind()
	Int32Kind = Int32Type.Kind()
	Int64Kind = Int64Type.Kind()

	UintKind   = UintType.Kind()
	Uint8Kind  = Uint8Type.Kind()
	Uint16Kind = Uint16Type.Kind()
	Uint32Kind = Uint32Type.Kind()
	Uint64Kind = Uint64Type.Kind()

	Float32Kind = Float32Type.Kind()
	Float64Kind = Float64Type.Kind()

	Complex64Kind  = Complex64Type.Kind()
	Complex128Kind = Complex128Type.Kind()

	StringKind = StringType.Kind()
	BoolKind   = BoolType.Kind()
	ByteKind   = ByteType.Kind()
	BytesKind  = BytesType.Kind()

	TimeKind = TimeType.Kind()
)

var SqlType2GoType = map[string]reflect.Type{
	"int":                IntType,
	"integer":            IntType,
	"tinyint":            IntType,
	"smallint":           IntType,
	"mediumint":          IntType,
	"bigint":             Int64Type,
	"int unsigned":       UintType,
	"integer unsigned":   UintType,
	"tinyint unsigned":   UintType,
	"smallint unsigned":  UintType,
	"mediumint unsigned": UintType,
	"bigint unsigned":    Uint64Type,
	"bit":                Int8Type,
	"bool":               BoolType,
	"enum":               StringType,
	"set":                StringType,
	"varchar":            StringType,
	"char":               StringType,
	"tinytext":           StringType,
	"mediumtext":         StringType,
	"text":               StringType,
	"longtext":           StringType,
	"blob":               StringType,
	"tinyblob":           StringType,
	"mediumblob":         StringType,
	"longblob":           StringType,
	"date":               TimeType,
	"datetime":           TimeType,
	"timestamp":          TimeType,
	"time":               TimeType,
	"float":              Float64Type,
	"double":             Float64Type,
	"decimal":            Float64Type,
	"binary":             StringType,
	"varbinary":          StringType,
}
