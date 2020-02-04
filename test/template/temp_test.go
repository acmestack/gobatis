// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description: 

package template

import (
	tmp2 "github.com/xfali/gobatis/parsing/template"
	"html/template"
	"os"
	"testing"
)

type TestTable struct {
	Id       int
	UserName string
	Password string
}

var driverName = "mysql"

func TestTemplate(t *testing.T) {
	tpl, err := template.ParseFiles("./sql.tpl")
	if err != nil {
		t.Fatal(err)
	}

	s := tpl.Templates()
	for _, v := range s {
		t.Log(v.Name())
	}

	var param = TestTable{Id: 1, UserName:"user", Password:"pw"}
	t.Run("select", func(t *testing.T) {
		tpl = tpl.Lookup("selectTestTable")
		if tpl == nil {
			t.Fatal("not found")
		}

		err = tpl.Execute(os.Stdout, param)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("insert", func(t *testing.T) {
		tpl = tpl.Lookup("insertTestTable")
		if tpl == nil {
			t.Fatal("not found")
		}

		err = tpl.Execute(os.Stdout, param)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("update", func(t *testing.T) {
		tpl = tpl.Lookup("updateTestTable")
		if tpl == nil {
			t.Fatal("not found")
		}

		err = tpl.Execute(os.Stdout, param)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("delete", func(t *testing.T) {
		tpl = tpl.Lookup("deleteTestTable")
		if tpl == nil {
			t.Fatal("not found")
		}

		err = tpl.Execute(os.Stdout, param)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestParser(t *testing.T) {
	mgr := tmp2.NewManager()
	mgr.RegisterFile("./sql.tpl")
	var param = TestTable{Id: 1, UserName:"user", Password:"pw"}
	t.Run("select", func(t *testing.T) {
		tmp, _ := mgr.FindSql("selectTestTable")
		md, err := tmp.ParseMetadata(driverName, param)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(md.PrepareSql)
	})

	t.Run("insert", func(t *testing.T) {
		tmp, _ := mgr.FindSql("insertTestTable")
		md, err := tmp.ParseMetadata(driverName, param)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(md.PrepareSql)
	})

	t.Run("update", func(t *testing.T) {
		tmp, _ := mgr.FindSql("updateTestTable")
		md, err := tmp.ParseMetadata(driverName, param)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(md.PrepareSql)
	})

	t.Run("delete", func(t *testing.T) {
		tmp, _ := mgr.FindSql("deleteTestTable")
		md, err := tmp.ParseMetadata(driverName, param)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(md.PrepareSql)
	})
}

