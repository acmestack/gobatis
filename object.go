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
    "sync"
)

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

func ParseObject(bean interface{}) (reflection.Object, error) {
    obj := findObject(bean)
    var err error
    if obj == nil {
        obj, err = reflection.GetObjectInfo(bean)
        if err != nil {
            return nil, err
        }

        cacheObject(obj)
    }
    obj = obj.New()
    obj.ResetValue(reflection.ReflectValue(bean))
    return obj, nil
}
