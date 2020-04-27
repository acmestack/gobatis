// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description: 

package sqlmanager

import (
    "github.com/xfali/gobatis"
    "testing"
)

func TestManager(t *testing.T) {
    err := gobatis.ScanMapperFile("./x")
    if err != nil {
        t.Fatal(err)
    }
}
