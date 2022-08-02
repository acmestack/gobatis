/*
 * Licensed to the AcmeStack under one or more contributor license
 * agreements. See the NOTICE file distributed with this work for
 * additional information regarding copyright ownership.
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

package cache

import (
	"fmt"
	"github.com/acmestack/gobatis/parsing/sqlparser"
	"sort"
	"strings"
	"sync"
)

type MetadataCache struct {
	cache map[MetadataCacheKey]*sqlparser.Metadata
	lock  sync.Mutex
}

type MetadataCacheKey string

var gMetadataCache = MetadataCache{
	cache: map[MetadataCacheKey]*sqlparser.Metadata{},
}

func FindMetadata(key MetadataCacheKey) *sqlparser.Metadata {
	gMetadataCache.lock.Lock()
	defer gMetadataCache.lock.Unlock()

	return gMetadataCache.cache[key]
}

func CacheMetadata(key MetadataCacheKey, data *sqlparser.Metadata) {
	gMetadataCache.lock.Lock()
	defer gMetadataCache.lock.Unlock()

	gMetadataCache.cache[key] = data
}

func CalcKey(sql string, params map[string]interface{}) MetadataCacheKey {
	buf := strings.Builder{}
	buf.WriteString(sql)
	list := make([]string, len(params))
	i := 0
	for k := range params {
		list[i] = k
		i++
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i] > list[j]
	})
	for i := range list {
		buf.WriteString(list[i])
		buf.WriteString(fmt.Sprintf("%v", params[list[i]]))
	}
	return MetadataCacheKey(buf.String())
}
