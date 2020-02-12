// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description: 

package template

import (
	"testing"
)

func TestReplace(t *testing.T) {
	old := "ab"
	new := "cdab"
	s, i := replace("fhaksfjlabdasdabdasljfabda", old, new, -1)
	t.Log(s, " ", i)
}

func TestFmt(t *testing.T) {
	s := getPlaceHolderKey(10)
	t.Log(s)
}
