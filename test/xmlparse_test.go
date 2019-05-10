/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package test

import (
    "github.com/xfali/gobatis"
    "github.com/xfali/gobatis/parsing/xml"
    "testing"
    "time"
)

func TestXml(t *testing.T) {
    //xmlFile := os.Getenv("xmlFile")
    //m, err := xml.ParseFile(xmlFile)
    m, err := xml.Parse([]byte(test_xml))
    if err != nil {
        t.Fatal(err)
    }

    t.Log(m)
}

func TestXmlFormat(t *testing.T) {
    //xmlFile := os.Getenv("xmlFile")
    //m, err := xml.ParseFile(xmlFile)
    m, err := xml.Parse([]byte(test_xml))
    if err != nil {
        t.Fatal(err)
    }

    ret := m.Format()
    for k, v := range ret {
        t.Logf("id : %s, v : %s\n", k, v)
    }
}

func TestXmlDynamic(t *testing.T) {
    src := `SELECT <include refid="select_field"></include> FROM PERSON
        <where>
            <if test="id != nil">
                id = #{id}
            </if>
        </where>
        and name = #{name}`
    m, err := xml.ParseDynamic(src, nil)
    if err != nil {
        t.Fatal(err)
    }

    t.Log(m)
}

func TestXmlDynamic2(t *testing.T) {
    src := `SELECT <include refid="select_field"></include> FROM PERSON
        <where>
            <if test="{0} != nil">
                OR id = #{0}
            </if>
            <if test="{0} != nil">
                OR id = #{0}
            </if>
        </where>
        and name = #{1}`
    m, err := xml.ParseDynamic(src, nil)
    if err != nil {
        t.Fatal(err)
    }

    str := "12345678"
    t.Logf("%s %s", str[:3], str[3:])

    t.Log(m.Replace(100, 200))
}

func TestXmlDynamic2_2(t *testing.T) {
    src := `SELECT <include refid="select_field"></include> FROM PERSON
        <where>
            <if test="{0} == never_true">
                OR id = #{0}
            </if>
            <if test="{0} != nil">
                OR second_id = #{0}
            </if>
        </where>
        and name = #{1}`
    m, err := xml.ParseDynamic(src, nil)
    if err != nil {
        t.Fatal(err)
    }

    str := "12345678"
    t.Logf("%s %s", str[:3], str[3:])

    t.Log(m.Replace(100, 200))
}

func TestXmlDynamic3(t *testing.T) {
    src := `SELECT <include refid="select_field"></include> FROM PERSON
        <where>
            <if test="{0} == 100">
                and id = #{0}
            </if>
        </where>
        and name = #{1}`
    m, err := xml.ParseDynamic(src, nil)
    if err != nil {
        t.Fatal(err)
    }

    t.Logf("arg 100: %s\n", m.Replace(100))
    t.Logf("arg 200: %s\n", m.Replace(200))
}

func TestXmlDynamic4(t *testing.T) {
    testV := TestTable{
        Username:"testuser",
        Password:"testpw",
    }
    gobatis.RegisterModel(&testV)
    src := `SELECT <include refid="select_field"></include> FROM PERSON
        <where>
            <if test="username == testuser">
                name = #{username}
            </if>
        </where>`
    m, err := xml.ParseDynamic(src, nil)
    if err != nil {
        t.Fatal(err)
    }

    t.Logf("arg testV testuser: %s\n", m.Replace(testV))

    testV.Username = "error_user"
    t.Logf("arg testV error_user: %s\n", m.Replace(testV))
}

func TestXmlDynamic5(t *testing.T) {
    testV := TestTable{
        Username:"testuser",
        Password:"testpw",
    }
    gobatis.RegisterModel(&testV)
    src := `UPDATE PERSON
        <set>
            <if test="username != nil">
                username = #{username}
            </if>
            <if test="password != nil">
                password = #{password}
            </if>
        </set>`
    m, err := xml.ParseDynamic(src, nil)
    if err != nil {
        t.Fatal(err)
    }

    t.Logf("arg testV testuser: %s\n", m.Replace(testV))

    testV.Username = ""
    t.Logf("arg testV error_user: %s\n", m.Replace(testV))
}

func TestXmlDynamicIf0(t *testing.T) {
    src := `SELECT * FROM test_table
            <if test="{0} != nil">
                where id = #{0}
            </if> `
    m, err := xml.ParseDynamic(src, nil)
    if err != nil {
        t.Fatal(err)
    }

    t.Logf("arg nil : %s\n", m.Replace(""))

    t.Logf("arg 0 : %s\n", m.Replace(0))

    t.Logf("arg 1 : %s\n", m.Replace(1))
}

func TestXmlDynamicIf1(t *testing.T) {
    src := `SELECT * FROM test_table
            <if test="{0} != nil and {0} != 0">
                where id = #{0}
            </if> `
    m, err := xml.ParseDynamic(src, nil)
    if err != nil {
        t.Fatal(err)
    }

    t.Logf("arg nil : %s\n", m.Replace(""))

    t.Logf("arg 0 : %s\n", m.Replace(0))

    t.Logf("arg 1 : %s\n", m.Replace(1))
}

func TestXmlDynamicIf2(t *testing.T) {
    src := `SELECT * FROM test_table
            <if test="{0} != nil or {0} != 0">
                where id = #{0}
            </if> `
    m, err := xml.ParseDynamic(src, nil)
    if err != nil {
        t.Fatal(err)
    }

    t.Logf("arg nil : %s\n", m.Replace(""))

    t.Logf("arg 0 : %s\n", m.Replace(0))

    t.Logf("arg 1 : %s\n", m.Replace(1))
}

func TestXmlDynamicIf3(t *testing.T) {
    src := `SELECT * FROM test_table
            <if test="{0} != nil">
                where id = #{0}
            </if> `
    m, err := xml.ParseDynamic(src, nil)
    if err != nil {
        t.Fatal(err)
    }

    t.Logf("arg nil : %s\n", m.Replace(""))

    t.Logf("arg zero time : %s\n", m.Replace(time.Time{}))

    t.Logf("arg now : %s\n", m.Replace(time.Now()))
}

func TestXml0(t *testing.T) {
    gobatis.RegisterMapperData([]byte(main_xml))
    testV := TestTable{}
    t.Logf("selectUser %s\n", gobatis.FindSql("selectUser").Replace(100))
    t.Logf("insertUser %s\n", gobatis.FindSql("insertUser").Replace(testV))
    testV.Password = "pw"
    t.Logf("updateUser %s\n", gobatis.FindSql("updateUser").Replace(testV))
    testV.Id = -1
    t.Logf("deleteUser %s\n", gobatis.FindSql("deleteUser").Replace(testV))
}
