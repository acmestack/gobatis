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

package gobatis

import (
	"github.com/acmestack/gobatis/errors"
	"github.com/acmestack/gobatis/reflection"
	"sync"
)

type ModelName string

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

// 注册struct模型，模型描述了column和field之间的关联关系；
// 目前已非必要条件
func RegisterModel(model interface{}) error {
	return RegisterModelWithName("", model)
}

func RegisterModelWithName(name string, model interface{}) error {
	tableInfo, err := reflection.GetObjectInfo(model)
	if err != nil {
		return errors.PARSE_MODEL_TABLEINFO_FAILED
	}

	globalObjectCache.lock.Lock()
	defer globalObjectCache.lock.Unlock()

	if name == "" {
		name = tableInfo.GetClassName()
	}
	globalObjectCache.objCache[name] = tableInfo
	return nil
}
