/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package cache

import (
    "fmt"
    "github.com/xfali/gobatis/parsing/sqlparser"
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
    for k, v := range params {
        buf.WriteString(k)
        buf.WriteString(fmt.Sprintf("%v", v))
    }
    return MetadataCacheKey(buf.String())
}
