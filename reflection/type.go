/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description:
 */

package reflection

import (
	"reflect"
	"time"
)

var (
	EMPTY_STRING       string
	BOOL_DEFAULT       bool
	BYTE_DEFAULT       byte
	COMPLEX64_DEFAULT  complex64
	COMPLEX128_DEFAULT complex128
	FLOAT32_DEFAULT    float32
	FLOAT64_DEFAULT    float64
	INT64_DEFAULT      int64
	UINT64_DEFAULT     uint64
	INT32_DEFAULT      int32
	UINT32_DEFAULT     uint32
	INT16_DEFAULT      int16
	UINT16_DEFAULT     uint16
	INT8_DEFAULT       int8
	UINT8_DEFAULT      uint8
	INT_DEFAULT        int
	UINT_DEFAULT       uint
	TIME_DEFAULT       time.Time
)

var (
	IntType   = reflect.TypeOf(INT_DEFAULT)
	Int8Type  = reflect.TypeOf(INT8_DEFAULT)
	Int16Type = reflect.TypeOf(INT16_DEFAULT)
	Int32Type = reflect.TypeOf(INT32_DEFAULT)
	Int64Type = reflect.TypeOf(INT64_DEFAULT)

	UintType   = reflect.TypeOf(UINT_DEFAULT)
	Uint8Type  = reflect.TypeOf(UINT8_DEFAULT)
	Uint16Type = reflect.TypeOf(UINT16_DEFAULT)
	Uint32Type = reflect.TypeOf(UINT32_DEFAULT)
	Uint64Type = reflect.TypeOf(UINT64_DEFAULT)

	Float32Type = reflect.TypeOf(FLOAT32_DEFAULT)
	Float64Type = reflect.TypeOf(FLOAT64_DEFAULT)

	Complex64Type  = reflect.TypeOf(COMPLEX64_DEFAULT)
	Complex128Type = reflect.TypeOf(COMPLEX128_DEFAULT)

	StringType = reflect.TypeOf(EMPTY_STRING)
	BoolType   = reflect.TypeOf(BOOL_DEFAULT)
	ByteType   = reflect.TypeOf(BYTE_DEFAULT)
	BytesType  = reflect.SliceOf(ByteType)

	TimeType = reflect.TypeOf(TIME_DEFAULT)
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
