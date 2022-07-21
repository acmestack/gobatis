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

	"github.com/acmestack/gobatis/common"
	"github.com/acmestack/gobatis/errors"
	"github.com/acmestack/gobatis/logging"
)

const (
	ObjectUnknown = iota
	ObjectSimpletype
	ObjectStruct
	ObjectSlice
	ObjectMap

	ObjectCustom = 50000
)

type Object interface {
	Kind() int
	// New 生成克隆空对象
	New() Object
	// NewElem 获得对象的元素
	NewElem() Object
	// SetField 设置字段
	SetField(name string, v interface{})
	// AddValue 添加元素值
	AddValue(v reflect.Value)
	// GetClassName 获得对象名称
	GetClassName() string
	// CanSetField 是否能设置field
	CanSetField() bool
	// CanAddValue 是否能添加值
	CanAddValue() bool

	// NewValue 获得值
	NewValue() reflect.Value

	// CanSet 是否能够设置
	CanSet(v reflect.Value) bool
	// SetValue 设置值
	SetValue(v reflect.Value)
	// GetValue 获得值
	GetValue() reflect.Value
	// ResetValue 变换value对象
	ResetValue(v reflect.Value)
}

var modelNameType reflect.Type

func SetModelNameType(mtype reflect.Type) {
	modelNameType = mtype
}

type Settable struct {
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

	Settable

	Newable
}

type SliceInfo struct {
	//包含pkg的名称
	ClassName string
	Elem      Object

	Settable
	Newable
}

type SimpleTypeInfo struct {
	//包含pkg的名称
	ClassName string

	Settable
	Newable
}

type MapInfo struct {
	//包含pkg的名称
	ClassName string
	//元素类型
	ElemType reflect.Type

	Settable
	Newable
}

func (settable *Settable) CanSet(v reflect.Value) bool {
	if settable.Value.Kind() != v.Kind() {
		logging.Warn("Set value failed")
		return false
	}
	destClass := GetTypeClassName(settable.Value.Type())
	srcClass := GetTypeClassName(v.Type())
	if destClass != srcClass {
		logging.Warn("different type: %settable %settable\n", destClass, srcClass)
		return false
	}
	return true
}

func (settable *Settable) SetValue(v reflect.Value) {
	if !settable.Value.IsValid() {
		settable.Value = v
	}

	if settable.CanSet(v) {
		settable.Value.Set(v)
	}
}

func (settable *Settable) ResetValue(v reflect.Value) {
	settable.Value = v
}

func (settable *Settable) GetValue() reflect.Value {
	return settable.Value
}

func (newable *Newable) NewValue() reflect.Value {
	return reflect.New(newable.Type).Elem()
}

func (structInfo *StructInfo) New() Object {
	ret := &StructInfo{
		ClassName:    structInfo.ClassName,
		Name:         structInfo.Name,
		FieldNameMap: structInfo.FieldNameMap,
	}
	ret.Type = structInfo.Type
	ret.Value = reflect.New(structInfo.Type).Elem()
	return ret
}

func (structInfo *StructInfo) NewElem() Object {
	return nil
}

func (structInfo *StructInfo) SetField(name string, ov interface{}) {
	fieldName := structInfo.FieldNameMap[name]
	if fieldName != "" {
		f := structInfo.Value.FieldByName(fieldName)
		if f.IsValid() {
			SetValue(f, ov)
		}
	}
}

func (structInfo *StructInfo) AddValue(v reflect.Value) {

}

func (structInfo *StructInfo) GetClassName() string {
	return structInfo.ClassName
}

func (structInfo *StructInfo) Kind() int {
	return ObjectStruct
}

func (structInfo *StructInfo) CanSetField() bool {
	return true
}

func (structInfo *StructInfo) CanAddValue() bool {
	return false
}

func (sliceInfo *SliceInfo) New() Object {
	ret := &SliceInfo{
		Elem: sliceInfo.Elem.New(),
	}
	ret.Type = sliceInfo.Type
	ret.Value = reflect.New(sliceInfo.Type).Elem()
	return ret
}

func (sliceInfo *SliceInfo) NewElem() Object {
	return sliceInfo.Elem.New()
}

func (sliceInfo *SliceInfo) SetField(name string, v interface{}) {
	logging.Info("slice not support SetField")
}

func (sliceInfo *SliceInfo) AddValue(v reflect.Value) {
	if !sliceInfo.Elem.CanSet(v) {
		logging.Warn("Add value failed, different kind")
		return
	}
	newValue := reflect.Append(sliceInfo.Value, v)
	//直接设置，不使用o.SetValue，效率更高
	sliceInfo.Value.Set(newValue)
}

func (sliceInfo *SliceInfo) GetClassName() string {
	return sliceInfo.ClassName
}

func (sliceInfo *SliceInfo) Kind() int {
	return ObjectSlice
}

func (sliceInfo *SliceInfo) CanSetField() bool {
	return false
}

func (sliceInfo *SliceInfo) CanAddValue() bool {
	return true
}

func (simpleTypeInfo *SimpleTypeInfo) New() Object {
	ret := &SimpleTypeInfo{
		ClassName: simpleTypeInfo.ClassName,
	}
	ret.Type = simpleTypeInfo.Type
	ret.Value = reflect.New(simpleTypeInfo.Type).Elem()
	return ret
}

func (simpleTypeInfo *SimpleTypeInfo) NewElem() Object {
	return nil
}

func (simpleTypeInfo *SimpleTypeInfo) SetField(name string, ov interface{}) {
	SetValue(simpleTypeInfo.Value, ov)
}

func (simpleTypeInfo *SimpleTypeInfo) AddValue(v reflect.Value) {

}

func (simpleTypeInfo *SimpleTypeInfo) GetClassName() string {
	return simpleTypeInfo.ClassName
}

// CanSet 直接返回true，需要通过SetValue判断
func (simpleTypeInfo *SimpleTypeInfo) CanSet(v reflect.Value) bool {
	return true
}

func (simpleTypeInfo *SimpleTypeInfo) SetValue(v reflect.Value) {
	if !simpleTypeInfo.Value.IsValid() {
		simpleTypeInfo.Value = v
	}

	if !SetValue(simpleTypeInfo.Value, v.Interface()) {
		logging.Warn("SimpleTypeInfo SetValue failed")
	}
}

func (simpleTypeInfo *SimpleTypeInfo) Kind() int {
	return ObjectSimpletype
}

func (simpleTypeInfo *SimpleTypeInfo) CanSetField() bool {
	return false
}

func (simpleTypeInfo *SimpleTypeInfo) CanAddValue() bool {
	return false
}

func (mapInfo *MapInfo) CanSet(v reflect.Value) bool {
	if mapInfo.Value.Kind() != v.Kind() {
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

// SetValue TODO: 目前仅支持map[string]interface{}，需增加其他类型支持
func (mapInfo *MapInfo) SetValue(v reflect.Value) {
	if mapInfo.CanSet(v) {
		mapInfo.Value.Set(v)
	}
}

func (mapInfo *MapInfo) New() Object {
	ret := &MapInfo{
		ClassName: mapInfo.ClassName,
		ElemType:  mapInfo.ElemType,
	}
	ret.Value = reflect.New(mapInfo.Type).Elem()
	return ret
}

// NewElem FIXME: return nil，需要对map元素解析
func (mapInfo *MapInfo) NewElem() Object {
	return nil
}

func (mapInfo *MapInfo) SetField(name string, ov interface{}) {
	v := reflect.New(mapInfo.ElemType).Elem()
	if SetValue(v, ov) {
		mapInfo.Value.SetMapIndex(reflect.ValueOf(name), v)
	}
}

func (mapInfo *MapInfo) AddValue(v reflect.Value) {

}

func (mapInfo *MapInfo) CanSetField() bool {
	return true
}

func (mapInfo *MapInfo) CanAddValue() bool {
	return false
}

func (mapInfo *MapInfo) Kind() int {
	return ObjectMap
}

func (mapInfo *MapInfo) GetClassName() string {
	return mapInfo.ClassName
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
	return nil, errors.ObjectNotSupport
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
		return nil, errors.ParseObjectNotSlice
	}
	//获得元素类型
	et := rt.Elem()
	ev := reflect.New(et).Elem()

	elemObj, err := GetReflectObjectInfo(et, ev)
	if err != nil {
		return nil, err
	}
	if elemObj.CanAddValue() {
		return nil, errors.SliceSliceNotSupport
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
		return nil, errors.ParseObjectNotMap
	}

	if rt.Key().Kind() != reflect.String {
		return nil, errors.GetObjectinfoFailed
	}

	//TODO: 目前仅支持map[string]interface{}，需增加其他类型支持
	if rt.Elem().Kind() != reflect.Interface {
		logging.Warn("Map type support map[string]interface{} only, but get map[%v]%v \n", rt.Key(), rt.Elem())
		return nil, errors.GetObjectinfoFailed
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
//3、如果结构体中含有column的tag，则：
// a）、如果tag为‘-’，则不进行columne与field的映射；
// b）、如果tag不为‘-’使用tag name作为column名称与field映射。
//4、如果结构体中不含有column的tag，则使用field name作为column名称与field映射
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
		return nil, errors.ParseObjectNotStruct
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
		tagName := rtf.Tag.Get(common.ColumnName)
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

func (structInfo *StructInfo) MapValue() map[string]interface{} {
	paramMap := map[string]interface{}{}
	structInfo.FillMapValue(&paramMap)
	return paramMap
}

func (structInfo *StructInfo) FillMapValue(paramMap *map[string]interface{}) {
	for k, v := range structInfo.FieldNameMap {
		f := structInfo.Value.FieldByName(v)
		if !f.CanInterface() {
			f = reflect.Indirect(f)
		}
		(*paramMap)[k] = f.Interface()
	}
	//(*paramMap)["tablename"] = structInfo.Name
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
