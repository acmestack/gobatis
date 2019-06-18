/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description:
 */

package reflection

import (
	"encoding/json"
	"fmt"
	"github.com/xfali/gobatis/errors"
	"github.com/xfali/gobatis/logging"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func ReflectValue(bean interface{}) reflect.Value {
	return reflect.Indirect(reflect.ValueOf(bean))
}

func IsSimpleObject(bean interface{}) bool {
	rt := reflect.TypeOf(bean)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	return IsSimpleType(rt)
}

//是否是数据库使用的简单类型，注意不能是PTR
func IsSimpleType(t reflect.Type) bool {
	switch t.Kind() {
	case IntKind, Int8Kind, Int16Kind, Int32Kind, Int64Kind, UintKind, Uint8Kind, Uint16Kind, Uint32Kind, Uint64Kind,
		Float32Kind, Float64Kind, Complex64Kind, Complex128Kind, StringKind, BoolKind, ByteKind /*, BytesKind, TimeKind*/ :
		return true
	}

	if t.ConvertibleTo(BytesType) || t.ConvertibleTo(TimeType) {
		return true
	}
	return false
}

func SafeSetValue(f reflect.Value, v interface{}) bool {
	if err := MustPtrValue(f); err != nil {
		logging.Info("value cannot be set: %s\n", err.Error())
		return false
	}
	f = f.Elem()
	return SetValue(f, v)
}

func SetValue(f reflect.Value, v interface{}) bool {
	if v == nil {
		return false
	}

	hasAssigned := false
	rawValue := reflect.Indirect(reflect.ValueOf(v))
	rawValueType := reflect.TypeOf(rawValue.Interface())
	vv := reflect.ValueOf(rawValue.Interface())

	ft := f.Type()
	switch ft.Kind() {
	case reflect.Bool:
		switch rawValueType.Kind() {
		case reflect.Bool:
			hasAssigned = true
			f.SetBool(vv.Bool())
			break
		case reflect.Slice:
			if d, ok := vv.Interface().([]uint8); ok {
				hasAssigned = true
				f.SetBool(d[0] != 0)
			}
			break
		}
		break
	case reflect.String:
		switch rawValueType.Kind() {
		case reflect.String:
			hasAssigned = true
			f.SetString(vv.String())
			break
		case reflect.Slice:
			if d, ok := vv.Interface().([]uint8); ok {
				hasAssigned = true
				f.SetString(string(d))
			}
			break
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
			hasAssigned = true
			f.SetString(strconv.FormatUint(vv.Uint(), 10))
			break
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			hasAssigned = true
			f.SetString(strconv.FormatInt(vv.Int(), 10))
			break
		case reflect.Float64:
			hasAssigned = true
			f.SetString(strconv.FormatFloat(vv.Float(), 'g', -1, 64))
			break
		case reflect.Float32:
			hasAssigned = true
			f.SetString(strconv.FormatFloat(vv.Float(), 'g', -1, 32))
			break
		case reflect.Bool:
			hasAssigned = true
			f.SetString(strconv.FormatBool(vv.Bool()))
			break
		//case reflect.Struct:
		//    if ti, ok := v.(time.Time); ok {
		//        hasAssigned = true
		//        if ti.IsZero() {
		//            f.SetString("")
		//        } else {
		//            f.SetString(ti.String())
		//        }
		//    } else {
		//        hasAssigned = true
		//        f.SetString(fmt.Sprintf("%v", v))
		//    }
		default:
			hasAssigned = true
			f.SetString(fmt.Sprintf("%v", v))
		}
		break
	case reflect.Complex64, reflect.Complex128:
		switch rawValueType.Kind() {
		case reflect.Complex64, reflect.Complex128:
			hasAssigned = true
			f.SetComplex(vv.Complex())
			break
		case reflect.Slice:
			if rawValueType.ConvertibleTo(BytesType) {
				d := vv.Bytes()
				if len(d) > 0 {
					if f.CanAddr() {
						err := json.Unmarshal(d, f.Addr().Interface())
						if err != nil {
							return false
						}
					} else {
						x := reflect.New(ft)
						err := json.Unmarshal(d, x.Interface())
						if err != nil {
							return false
						}
						hasAssigned = true
						f.Set(x.Elem())
						break
					}
				}
			}
			break
		}
		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch rawValueType.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			hasAssigned = true
			f.SetInt(vv.Int())
			break
		case reflect.Slice:
			if d, ok := vv.Interface().([]uint8); ok {
				intV, err := strconv.ParseInt(string(d), 10, 64)
				if err == nil {
					hasAssigned = true
					f.SetInt(intV)
				}
			}
			break
		}
		break
	case reflect.Float32, reflect.Float64:
		switch rawValueType.Kind() {
		case reflect.Float32, reflect.Float64:
			hasAssigned = true
			f.SetFloat(vv.Float())
			break
		case reflect.Slice:
			if d, ok := vv.Interface().([]uint8); ok {
				floatV, err := strconv.ParseFloat(string(d), 64)
				if err == nil {
					hasAssigned = true
					f.SetFloat(floatV)
				}
			}
			break
		}
		break
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		switch rawValueType.Kind() {
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
			hasAssigned = true
			f.SetUint(vv.Uint())
			break
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			hasAssigned = true
			f.SetUint(uint64(vv.Int()))
			break
		case reflect.Slice:
			if d, ok := vv.Interface().([]uint8); ok {
				uintV, err := strconv.ParseUint(string(d), 10, 64)
				if err == nil {
					hasAssigned = true
					f.SetUint(uintV)
				}
			}
			break
		}
		break
	case reflect.Struct:
		fieldType := f.Type()
		if fieldType.ConvertibleTo(TimeType) {
			if rawValueType == TimeType {
				hasAssigned = true
				t := vv.Convert(TimeType).Interface().(time.Time)
				f.Set(reflect.ValueOf(t).Convert(fieldType))
			} else if rawValueType == IntType || rawValueType == Int64Type ||
				rawValueType == Int32Type {
				hasAssigned = true

				t := time.Unix(vv.Int(), 0)
				f.Set(reflect.ValueOf(t).Convert(fieldType))
			} else {
				if d, ok := vv.Interface().([]byte); ok {
					t, err := convert2Time(d, time.Local)
					if err == nil {
						hasAssigned = true
						f.Set(reflect.ValueOf(t).Convert(fieldType))
					}
				}
			}
		} else {
			f.Set(reflect.ValueOf(v))
		}
		break
	case reflect.Interface:
		hasAssigned = true
		f.Set(vv)
		break
	}

	return hasAssigned
}

const (
	zeroTime0 = "0000-00-00 00:00:00"
	zeroTime1 = "0001-01-01 00:00:00"
)

func convert2Time(data []byte, location *time.Location) (time.Time, error) {
	timeStr := strings.TrimSpace(string(data))
	var timeRet time.Time
	var err error
	if timeStr == zeroTime0 || timeStr == zeroTime1 {
	} else if !strings.ContainsAny(timeStr, "- :") {
		// time stamp
		sd, err := strconv.ParseInt(timeStr, 10, 64)
		if err == nil {
			timeRet = time.Unix(sd, 0)
		}
	} else if len(timeStr) > 19 && strings.Contains(timeStr, "-") {
		timeRet, err = time.ParseInLocation(time.RFC3339Nano, timeStr, location)
		if err != nil {
			timeRet, err = time.ParseInLocation("2006-01-02 15:04:05.999999999", timeStr, location)
		}
		if err != nil {
			timeRet, err = time.ParseInLocation("2006-01-02 15:04:05.9999999 Z07:00", timeStr, location)
		}
	} else if len(timeStr) == 19 && strings.Contains(timeStr, "-") {
		timeRet, err = time.ParseInLocation("2006-01-02 15:04:05", timeStr, location)
	} else if len(timeStr) == 10 && timeStr[4] == '-' && timeStr[7] == '-' {
		timeRet, err = time.ParseInLocation("2006-01-02", timeStr, location)
	}
	return timeRet, nil
}

func MustPtr(bean interface{}) error {
	return MustPtrValue(reflect.ValueOf(bean))
}

func MustPtrValue(beanValue reflect.Value) error {
	if beanValue.Kind() != reflect.Ptr {
		return errors.RESULT_ISNOT_POINTER
	} else if beanValue.Elem().Kind() == reflect.Ptr {
		return errors.RESULT_PTR_VALUE_IS_POINTER
	}
	return nil
}

func ToSlice(arr interface{}) []interface{} {
	v := reflect.ValueOf(arr)
	if v.Kind() != reflect.Slice {
		panic("toslice arr not slice")
	}
	l := v.Len()
	ret := make([]interface{}, l)
	for i := 0; i < l; i++ {
		ret[i] = v.Index(i).Interface()
	}
	return ret
}

func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	return false
}

func CanSet(i interface{}) bool {
	if i == nil {
		return false
	}
	vi := reflect.ValueOf(i)
	if MustPtrValue(vi) != nil {
		return false
	}
	if vi.IsNil() {
		return false
	}
	return true
}

func NewValue(t reflect.Type) reflect.Value {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Ptr {
		panic("Error Type")
	}
	return reflect.New(t).Elem()
}

func New(t reflect.Type) interface{} {
	return NewValue(t).Interface()
}
