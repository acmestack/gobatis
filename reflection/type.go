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
    c_EMPTY_STRING       string
    c_BOOL_DEFAULT       bool
    c_BYTE_DEFAULT       byte
    c_COMPLEX64_DEFAULT  complex64
    c_COMPLEX128_DEFAULT complex128
    c_FLOAT32_DEFAULT    float32
    c_FLOAT64_DEFAULT    float64
    c_INT64_DEFAULT      int64
    c_UINT64_DEFAULT     uint64
    c_INT32_DEFAULT      int32
    c_UINT32_DEFAULT     uint32
    c_INT16_DEFAULT      int16
    c_UINT16_DEFAULT     uint16
    c_INT8_DEFAULT       int8
    c_UINT8_DEFAULT      uint8
    c_INT_DEFAULT        int
    c_UINT_DEFAULT       uint
    c_TIME_DEFAULT       time.Time
)

var (
    IntType   = reflect.TypeOf(c_INT_DEFAULT)
    Int8Type  = reflect.TypeOf(c_INT8_DEFAULT)
    Int16Type = reflect.TypeOf(c_INT16_DEFAULT)
    Int32Type = reflect.TypeOf(c_INT32_DEFAULT)
    Int64Type = reflect.TypeOf(c_INT64_DEFAULT)

    UintType   = reflect.TypeOf(c_UINT_DEFAULT)
    Uint8Type  = reflect.TypeOf(c_UINT8_DEFAULT)
    Uint16Type = reflect.TypeOf(c_UINT16_DEFAULT)
    Uint32Type = reflect.TypeOf(c_UINT32_DEFAULT)
    Uint64Type = reflect.TypeOf(c_UINT64_DEFAULT)

    Float32Type = reflect.TypeOf(c_FLOAT32_DEFAULT)
    Float64Type = reflect.TypeOf(c_FLOAT64_DEFAULT)

    Complex64Type  = reflect.TypeOf(c_COMPLEX64_DEFAULT)
    Complex128Type = reflect.TypeOf(c_COMPLEX128_DEFAULT)

    StringType = reflect.TypeOf(c_EMPTY_STRING)
    BoolType   = reflect.TypeOf(c_BOOL_DEFAULT)
    ByteType   = reflect.TypeOf(c_BYTE_DEFAULT)
    BytesType  = reflect.SliceOf(ByteType)

    TimeType = reflect.TypeOf(c_TIME_DEFAULT)
)


