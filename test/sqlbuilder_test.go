/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package test

import (
    "github.com/xfali/gobatis/builder"
    "strings"
    "testing"
)

func TestSqlBuilderSelect(t *testing.T) {
    hook := func(f *builder.SQLFragment) *builder.SQLFragment {
        t.Log(f)
        return f
    }
    t.Run("once call", func(t *testing.T) {
        str := builder.Select("A.test1", "B.test2").
            Hook(hook).
            From("test_a").
            Hook(hook).
            Where("id = 1").
            Hook(hook).
            And().
            Hook(hook).
            Where("name=2").
            Hook(hook).
            GroupBy("name").
            Hook(hook).
            OrderBy("name").
            Hook(hook).
            Decs().
            Hook(hook).
            String()
        t.Log(str)

        if strings.TrimSpace(str) != `SELECT A.test1, B.test2 FROM test_a WHERE id = 1 AND name=2 GROUP BY name ORDER BY name DECS` {
            t.FailNow()
        }
    })

    t.Run("multi call", func(t *testing.T) {
        str := builder.Select("A.test1", "B.test2").
            Hook(hook).
            Select("test3").
            Hook(hook).
            From("test_a AS A").
            Hook(hook).
            From("test_b AS B").
            Hook(hook).
            Where("id = 1").
            Hook(hook).
            And().
            Hook(hook).
            Where("name=2").
            Hook(hook).
            GroupBy("name").
            Hook(hook).
            OrderBy("name").
            Hook(hook).
            Decs().
            Hook(hook).
            String()
        t.Log(str)

        if strings.TrimSpace(str) != `SELECT A.test1, B.test2, test3 FROM test_a AS A, test_b AS B WHERE id = 1 AND name=2 GROUP BY name ORDER BY name DECS` {
            t.FailNow()
        }
    })
}

func TestSqlBuilderInsert(t *testing.T) {
    str := builder.InsertInto("test_table").
        IntoColumns("a", "b").
        IntoColumns("c").
        IntoValues("#{0}, #{1}").
        IntoValues("#{3}").
        String()
    t.Log(str)

    if strings.TrimSpace(str) != `INSERT INTO test_table (a, b, c) VALUES(#{0}, #{1}, #{3})` {
        t.FailNow()
    }
}

func TestSqlBuilderUpdate(t *testing.T) {
    str := builder.Update("test_table").
        Set("a", "#{0}").
        Set("b", "#{1}").
        Where("id = #{3}").
        Or().
        Where("name = #{4}").
        String()
    t.Log(str)
    if strings.TrimSpace(str) != `UPDATE test_table SET a = #{0} , b = #{1} WHERE id = #{3} OR name = #{4}` {
        t.FailNow()
    }
}

func TestSqlBuilderDelete(t *testing.T) {
    str := builder.DeleteFrom("test_table").
        Where("id = #{3}").
        Or().
        Where("name = #{4}").
        String()
    t.Log(str)
    if strings.TrimSpace(str) != `DELETE FROM test_table WHERE id = #{3} OR name = #{4}` {
        t.FailNow()
    }
}

func TestSqlBuilderWhere(t *testing.T) {
    str := builder.DeleteFrom("test_table").
        Or().
        Where("name = #{4}").
        String()
    t.Log(str)

    if strings.TrimSpace(str) != `DELETE FROM test_table WHERE name = #{4}` {
        t.FailNow()
    }
}

func TestSqlBuilderError(t *testing.T) {
    f := builder.InsertInto("test_table")
    f.IntoColumns("a").IntoValues("#{0}").IntoValues("a").IntoColumns("a")

    t.Log(f.String())
}
