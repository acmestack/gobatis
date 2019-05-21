/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package gobatis

import (
    "github.com/xfali/gobatis/reflection"
    "reflect"
    "sync"
)

type ObjectInfo struct {
    obj  reflection.Object
}

type ObjectCache struct {
    objCache map[string]reflection.Object
    lock     sync.Mutex
}

var globalObjectCache = ObjectCache{
    objCache: map[string]reflection.Object{},
}

func findObject(bean interface{}) reflection.Object {
    classname := reflection.GetBeanClassName(bean)
    globalObjectCache.lock.Lock()
    defer globalObjectCache.lock.Unlock()

    return globalObjectCache.objCache[classname]
}

func cacheObject(obj reflection.Object) {
    globalObjectCache.lock.Lock()
    defer globalObjectCache.lock.Unlock()

    globalObjectCache.objCache[obj.GetClassName()] = obj
}

func ParseObject(bean interface{}) (ObjectInfo, error) {
    obj := findObject(bean)
    ret := ObjectInfo{}
    var err error
    if obj == nil {
        obj, err = reflection.GetObjectInfo(bean)
        if err != nil {
            return ObjectInfo{}, err
        }

        cacheObject(obj)
    }
    ret.obj = obj.New()
    ret.obj.ResetValue(reflection.ReflectValue(bean))
    return ret, nil
}

func (o *ObjectInfo) Deserialize(columns []string, values []interface{}) (interface{}, error) {
    obj := o.obj
    if o.obj.CanAddValue() {
        obj = o.obj.NewElem()
    }
    for i := range columns {
        if obj.CanSetField() {
            obj.SetField(columns[i], values[i])
        } else {
            obj.SetValue(reflect.ValueOf(values[0]))
            break
        }
    }
    if o.obj.CanAddValue() {
        o.obj.AddValue(obj.GetValue())
    }
    //已经add，不需要返回interface，需要注意的是IterFunc的参数为nil
    //return obj.GetValue().Interface(), nil
    return nil, nil
}
