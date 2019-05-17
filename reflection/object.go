/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package reflection

import (
    "github.com/xfali/gobatis/common"
    "github.com/xfali/gobatis/errors"
    "github.com/xfali/gobatis/logging"
    "reflect"
)

type Object interface {
    NewValue() reflect.Value
    SetValue(v reflect.Value)
    SetField(name string, v interface{})
    Add(v reflect.Value)
    GetClassName() string
}

var modelNameType reflect.Type

func SetModelNameType(mtype reflect.Type) {
    modelNameType = mtype
}

type Setable struct {
    //值
    Value reflect.Value
}

type StructInfo struct {
    //包含pkg的名称
    ClassName string
    //Model名称（目前用于xml解析是struct的前缀：#{x.username} 中的x）
    Name string
    //表字段和实体字段映射关系
    FieldNameMap map[string]string
    //类型
    Type reflect.Type

    Setable
}

type SliceInfo struct {
    elem Object
    Setable
}

type SimpleTypeInfo struct {
    //包含pkg的名称
    ClassName string
    //类型
    Type reflect.Type

    Setable
}

type MapInfo struct {
    //包含pkg的名称
    ClassName string
    //类型
    Type reflect.Type

    Setable
}

func (o *Setable) SetValue(v reflect.Value) {
    if o.Value.Kind() != v.Kind() {
        logging.Warn("Set value failed")
        return
    }
    o.Value.Set(v)
}

func newStructInfo() *StructInfo {
    return &StructInfo{
        FieldNameMap: map[string]string{},
    }
}

func (o *StructInfo) NewValue() reflect.Value {
    return reflect.New(o.Type).Elem()
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

func (o *StructInfo) Add(v reflect.Value) {

}

func (o *StructInfo) GetClassName() string {
    return o.ClassName
}

func (o *SliceInfo) NewValue() reflect.Value {
    return o.elem.NewValue()
}

func (o *SliceInfo) SetField(name string, v interface{}) {

}

func (o *SliceInfo) Add(v reflect.Value) {
    //FIXME: 可能需要重新设置Value值
    o.Value = reflect.Append(o.Value, v)
}

func (o *SliceInfo) GetClassName() string {
    return o.elem.GetClassName()
}

func (o *SimpleTypeInfo) NewValue() reflect.Value {
    return reflect.New(o.Type).Elem()
}

func (o *SimpleTypeInfo) SetField(name string, ov interface{}) {
    SetValue(o.Value, ov)
}

func (o *SimpleTypeInfo) Add(v reflect.Value) {

}

func (o *SimpleTypeInfo) GetClassName() string {
    return o.ClassName
}

//TODO: 目前仅支持map[string]interface{}，需增加其他类型支持
func (o *MapInfo) SetValue(v reflect.Value) {
    if o.Value.Kind() != v.Kind() {
        logging.Warn("Set value failed")
        return
    }
    rt := v.Type()
    if rt.Key().Kind() != reflect.String || rt.Elem().Kind() != reflect.Interface {
        logging.Warn("Map type support map[string]interface{} only")
        return
    }

    o.Value.Set(v)
}

func (o *MapInfo) NewValue() reflect.Value {
    return reflect.New(o.Type).Elem()
}

func (o *MapInfo) SetField(name string, ov interface{}) {
    v := o.NewValue()
    if SetValue(v, ov) {
        o.Value.SetMapIndex(reflect.ValueOf(name), v)
    }
}

func (o *MapInfo) Add(v reflect.Value) {

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
        Type:      rt,
        ClassName: getTypeClassName(rt),
    }
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
        return nil, errors.PARSE_TABLEINFO_NOT_SLICE
    }

    //获得元素类型
    rt = rt.Elem()

    info, err := GetReflectObjectInfo(rt, reflect.Value{})
    if err != nil {
        return nil, errors.GET_OBJECTINFO_FAILED
    }
    ret := SliceInfo{elem: info}
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
        return nil, errors.PARSE_TABLEINFO_NOT_MAP
    }

    if rt.Key().Kind() != reflect.String {
        return nil, errors.GET_OBJECTINFO_FAILED
    }

    //TODO: 目前仅支持map[string]interface{}，需增加其他类型支持
    if rt.Elem().Kind() != reflect.Interface {
        logging.Warn("Map type support map[string]interface{} only, but get map[%v]%v \n", rt.Key(), rt.Elem())
        return nil, errors.GET_OBJECTINFO_FAILED
    }

    ret := MapInfo{Type: rt.Elem(), ClassName: getTypeClassName(rt)}
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
        return nil, errors.PARSE_TABLEINFO_NOT_STRUCT
    }
    objInfo := newStructInfo()
    objInfo.Type = rt
    objInfo.Value = rv
    //Default name is struct name
    objInfo.Name = rt.Name()
    objInfo.ClassName = rt.PkgPath() + "/" + objInfo.Name

    //字段解析
    for i, j := 0, rt.NumField(); i < j; i++ {
        rtf := rt.Field(i)

        if rtf.Type == modelNameType {
            if rtf.Tag != "" {
                objInfo.Name = string(rtf.Tag)
            } else {
                objInfo.Name = rtf.Name
            }
            continue
        }

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
    return objInfo, nil
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
    return getTypeClassName(rt)
}

func getTypeClassName(rt reflect.Type) string {
    if rt.Kind() == reflect.Ptr {
        rt = rt.Elem()
    }

    if rt.Kind() == reflect.Slice {
        rt = rt.Elem()
    }
    path := rt.PkgPath()
    if path == "" {
        return rt.Name()
    } else {
        return path + "/" + rt.Name()
    }
}
