/*
 * Copyright (c) 2022, OpeningO
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
	"github.com/xfali/gobatis/common"
	"github.com/xfali/gobatis/errors"
	"github.com/xfali/gobatis/logging"
	"reflect"
)

const (
	OBJECT_UNKNOWN = iota
	OBJECT_SIMPLETYPE
	OBJECT_STRUCT
	OBJECT_SLICE
	OBJECT_MAP

	OBJECT_COSTOM = 50000
)

type Object interface {
	Kind() int
	//生成克隆空对象
	New() Object
	//获得对象的元素
	NewElem() Object
	//设置字段
	SetField(name string, v interface{})
	//添加元素值
	AddValue(v reflect.Value)
	//获得对象名称
	GetClassName() string
	//是否能设置field
	CanSetField() bool
	//是否能添加值
	CanAddValue() bool

	//获得值
	NewValue() reflect.Value

	//是否能够设置
	CanSet(v reflect.Value) bool
	//设置值
	SetValue(v reflect.Value)
	//获得值
	GetValue() reflect.Value
	//变换value对象
	ResetValue(v reflect.Value)
}

var modelNameType reflect.Type

func SetModelNameType(mtype reflect.Type) {
	modelNameType = mtype
}

type Setable struct {
	//值
	Value reflect.Value
}

type Newable struct {
	Type reflect.Type
}

type StructInfo struct {
	//包含pkg的名称
	ClassName string
	//Model名称（目前用于xml解析是struct的前缀：#{x.username} 中的x）
	Name string
	//表字段和实体字段映射关系
	FieldNameMap map[string]string

	Setable

	Newable
}

type SliceInfo struct {
	//包含pkg的名称
	ClassName string
	Elem      Object

	Setable
	Newable
}

type SimpleTypeInfo struct {
	//包含pkg的名称
	ClassName string

	Setable
	Newable
}

type MapInfo struct {
	//包含pkg的名称
	ClassName string
	//元素类型
	ElemType reflect.Type

	Setable
	Newable
}

func (o *Setable) CanSet(v reflect.Value) bool {
	if o.Value.Kind() != v.Kind() {
		logging.Warn("Set value failed")
		return false
	}
	destClass := GetTypeClassName(o.Value.Type())
	srcClass := GetTypeClassName(v.Type())
	if destClass != srcClass {
		logging.Warn("different type: %s %s\n", destClass, srcClass)
		return false
	}
	return true
}

func (o *Setable) SetValue(v reflect.Value) {
	if !o.Value.IsValid() {
		o.Value = v
	}

	if o.CanSet(v) {
		o.Value.Set(v)
	}
}

func (o *Setable) ResetValue(v reflect.Value) {
	o.Value = v
}

func (o *Setable) GetValue() reflect.Value {
	return o.Value
}

func (o *Newable) NewValue() reflect.Value {
	return reflect.New(o.Type).Elem()
}

func (o *StructInfo) New() Object {
	ret := &StructInfo{
		ClassName:    o.ClassName,
		Name:         o.Name,
		FieldNameMap: o.FieldNameMap,
	}
	ret.Type = o.Type
	ret.Value = reflect.New(o.Type).Elem()
	return ret
}

func (o *StructInfo) NewElem() Object {
	return nil
}

func (o *StructInfo) SetField(name string, ov interface{}) {
	fieldName := o.FieldNameMap[name]
	if fieldName != "" {
		f := o.Value.FieldByName(fieldName)
		if f.IsValid() {
			SetValue(f, ov)
		}
	}
}

func (o *StructInfo) AddValue(v reflect.Value) {

}

func (o *StructInfo) GetClassName() string {
	return o.ClassName
}

func (o *StructInfo) Kind() int {
	return OBJECT_STRUCT
}

func (o *StructInfo) CanSetField() bool {
	return true
}

func (o *StructInfo) CanAddValue() bool {
	return false
}

func (o *SliceInfo) New() Object {
	ret := &SliceInfo{
		Elem: o.Elem.New(),
	}
	ret.Type = o.Type
	ret.Value = reflect.New(o.Type).Elem()
	return ret
}

func (o *SliceInfo) NewElem() Object {
	return o.Elem.New()
}

func (o *SliceInfo) SetField(name string, v interface{}) {
	logging.Info("slice not support SetField")
}

func (o *SliceInfo) AddValue(v reflect.Value) {
	if !o.Elem.CanSet(v) {
		logging.Warn("Add value failed, different kind")
		return
	}
	newValue := reflect.Append(o.Value, v)
	//直接设置，不使用o.SetValue，效率更高
	o.Value.Set(newValue)
}

func (o *SliceInfo) GetClassName() string {
	return o.ClassName
}

func (o *SliceInfo) Kind() int {
	return OBJECT_SLICE
}

func (o *SliceInfo) CanSetField() bool {
	return false
}

func (o *SliceInfo) CanAddValue() bool {
	return true
}

func (o *SimpleTypeInfo) New() Object {
	ret := &SimpleTypeInfo{
		ClassName: o.ClassName,
	}
	ret.Type = o.Type
	ret.Value = reflect.New(o.Type).Elem()
	return ret
}

func (o *SimpleTypeInfo) NewElem() Object {
	return nil
}

func (o *SimpleTypeInfo) SetField(name string, ov interface{}) {
	SetValue(o.Value, ov)
}

func (o *SimpleTypeInfo) AddValue(v reflect.Value) {

}

func (o *SimpleTypeInfo) GetClassName() string {
	return o.ClassName
}

//直接返回true，需要通过SetValue判断
func (o *SimpleTypeInfo) CanSet(v reflect.Value) bool {
	return true
}

func (o *SimpleTypeInfo) SetValue(v reflect.Value) {
	if !o.Value.IsValid() {
		o.Value = v
	}

	if !SetValue(o.Value, v.Interface()) {
		logging.Warn("SimpleTypeInfo SetValue failed")
	}
}

func (o *SimpleTypeInfo) Kind() int {
	return OBJECT_SIMPLETYPE
}

func (o *SimpleTypeInfo) CanSetField() bool {
	return false
}

func (o *SimpleTypeInfo) CanAddValue() bool {
	return false
}

func (o *MapInfo) CanSet(v reflect.Value) bool {
	if o.Value.Kind() != v.Kind() {
		logging.Warn("Set value failed")
		return false
	}
	rt := v.Type()
	if rt.Key().Kind() != reflect.String || rt.Elem().Kind() != reflect.Interface {
		logging.Warn("Map type support map[string]interface{} only")
		return false
	}
	return true
}

//TODO: 目前仅支持map[string]interface{}，需增加其他类型支持
func (o *MapInfo) SetValue(v reflect.Value) {
	if o.CanSet(v) {
		o.Value.Set(v)
	}
}

func (o *MapInfo) New() Object {
	ret := &MapInfo{
		ClassName: o.ClassName,
		ElemType:  o.ElemType,
	}
	ret.Value = reflect.New(o.Type).Elem()
	return ret
}

//FIXME: return nil，需要对map元素解析
func (o *MapInfo) NewElem() Object {
	return nil
}

func (o *MapInfo) SetField(name string, ov interface{}) {
	v := reflect.New(o.ElemType).Elem()
	if SetValue(v, ov) {
		o.Value.SetMapIndex(reflect.ValueOf(name), v)
	}
}

func (o *MapInfo) AddValue(v reflect.Value) {

}

func (o *MapInfo) CanSetField() bool {
	return true
}

func (o *MapInfo) CanAddValue() bool {
	return false
}

func (o *MapInfo) Kind() int {
	return OBJECT_MAP
}

func (o *MapInfo) GetClassName() string {
	return o.ClassName
}

func GetObjectInfo(model interface{}) (Object, error) {
	rt := reflect.TypeOf(model)
	rv := reflect.ValueOf(model)
	if err := MustPtrValue(rv); err != nil {
		return nil, err
	}

	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
		rv = rv.Elem()
	}
	return GetReflectObjectInfo(rt, rv)
}

func GetReflectObjectInfo(rt reflect.Type, rv reflect.Value) (Object, error) {
	if IsSimpleType(rt) {
		return GetReflectSimpleTypeInfo(rt, rv)
	}
	switch rt.Kind() {
	case reflect.Struct:
		return GetReflectStructInfo(rt, rv)
	case reflect.Slice:
		return GetReflectSliceInfo(rt, rv)
	case reflect.Map:
		return GetReflectMapInfo(rt, rv)
	}
	return nil, errors.OBJECT_NOT_SUPPORT
}

func GetReflectSimpleTypeInfo(rt reflect.Type, rv reflect.Value) (Object, error) {
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	ret := SimpleTypeInfo{
		ClassName: GetTypeClassName(rt),
	}
	ret.Type = rt
	ret.Value = rv
	return &ret, nil
}

func GetReflectSliceInfo(rt reflect.Type, rv reflect.Value) (Object, error) {
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
		rv = rv.Elem()
	}

	kind := rt.Kind()

	if kind != reflect.Slice {
		return nil, errors.PARSE_OBJECT_NOT_SLICE
	}
	//获得元素类型
	et := rt.Elem()
	ev := reflect.New(et).Elem()

	elemObj, err := GetReflectObjectInfo(et, ev)
	if err != nil {
		return nil, err
	}
	if elemObj.CanAddValue() {
		return nil, errors.SLICE_SLICE_NOT_SUPPORT
	}
	ret := SliceInfo{Elem: elemObj, ClassName: GetTypeClassName(rt)}
	ret.Type = rt
	ret.Value = rv
	return &ret, nil
}

func GetReflectMapInfo(rt reflect.Type, rv reflect.Value) (Object, error) {
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
		rv = rv.Elem()
	}

	kind := rt.Kind()

	if kind != reflect.Map {
		return nil, errors.PARSE_OBJECT_NOT_MAP
	}

	if rt.Key().Kind() != reflect.String {
		return nil, errors.GET_OBJECTINFO_FAILED
	}

	//TODO: 目前仅支持map[string]interface{}，需增加其他类型支持
	if rt.Elem().Kind() != reflect.Interface {
		logging.Warn("Map type support map[string]interface{} only, but get map[%v]%v \n", rt.Key(), rt.Elem())
		return nil, errors.GET_OBJECTINFO_FAILED
	}

	ret := MapInfo{ElemType: rt.Elem(), ClassName: GetTypeClassName(rt)}
	ret.Type = rt
	ret.Value = rv
	return &ret, nil
}

//GetStructInfo 解析结构体，使用：
//1、如果结构体中含有gobatis.ModelName类型的字段，则：
// a)、如果含有tag，则使用tag作为tablename；
// b)、如果不含有tag，则使用fieldName作为tablename。
//2、如果结构体中不含有gobatis.ModelName类型的字段，则使用结构体名称作为tablename
//3、如果结构体中含有xfield的tag，则：
// a）、如果tag为‘-’，则不进行columne与field的映射；
// b）、如果tag不为‘-’使用tag name作为column名称与field映射。
//4、如果结构体中不含有xfield的tag，则使用field name作为column名称与field映射
//5、如果字段的tag为‘-’，则不进行columne与field的映射；
func GetStructInfo(bean interface{}) (*StructInfo, error) {
	return GetReflectStructInfo(reflect.TypeOf(bean), reflect.ValueOf(bean))
}

func GetReflectStructInfo(rt reflect.Type, rv reflect.Value) (*StructInfo, error) {
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
		if rv.IsValid() {
			rv = rv.Elem()
		}
	}

	kind := rt.Kind()

	if kind != reflect.Struct {
		return nil, errors.PARSE_OBJECT_NOT_STRUCT
	}
	objInfo := StructInfo{
		FieldNameMap: map[string]string{},
	}
	objInfo.Type = rt
	objInfo.Value = rv
	//Default name is struct name
	objInfo.Name = rt.Name()
	objInfo.ClassName = GetTypeClassName(rt)

	//字段解析
	for i, j := 0, rt.NumField(); i < j; i++ {
		rtf := rt.Field(i)

		//if rtf.Type == modelNameType {
		//    if rtf.Tag != "" {
		//        objInfo.Name = string(rtf.Tag)
		//    } else {
		//        objInfo.Name = rtf.Name
		//    }
		//    continue
		//}

		//没有tag,表字段名与实体字段名一致
		if rtf.Tag == "" {
			objInfo.FieldNameMap[rtf.Name] = rtf.Name
			continue
		}

		if rtf.Tag == "-" {
			continue
		}

		fieldName := rtf.Name
		tagName := rtf.Tag.Get(common.FIELD_NAME)
		if tagName == "-" {
			continue
		} else if tagName != "" {
			fieldName = tagName
		}
		objInfo.FieldNameMap[fieldName] = rtf.Name
		continue
	}
	return &objInfo, nil
}

func (ti *StructInfo) MapValue() map[string]interface{} {
	paramMap := map[string]interface{}{}
	ti.FillMapValue(&paramMap)
	return paramMap
}

func (ti *StructInfo) FillMapValue(paramMap *map[string]interface{}) {
	for k, v := range ti.FieldNameMap {
		f := ti.Value.FieldByName(v)
		if !f.CanInterface() {
			f = reflect.Indirect(f)
		}
		(*paramMap)[k] = f.Interface()
	}
	//(*paramMap)["tablename"] = ti.Name
}

func GetBeanClassName(model interface{}) string {
	rt := reflect.TypeOf(model)
	return GetTypeClassName(rt)
}

func GetTypeClassName(rt reflect.Type) string {
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	return rt.String()
}
