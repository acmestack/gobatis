/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package test

import (
    "github.com/xfali/gobatis/config"
    "github.com/xfali/gobatis/parsing/xml"
    "testing"
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
    config.RegisterModel(&testV)
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
    config.RegisterModel(&testV)
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

func TestXml0(t *testing.T) {
    config.RegisterMapperData([]byte(main_xml))
    testV := TestTable{}
    t.Logf("selectUser %s\n", config.FindSql("selectUser").Replace(100))
    t.Logf("insertUser %s\n", config.FindSql("insertUser").Replace(testV))
    testV.Password = "pw"
    t.Logf("updateUser %s\n", config.FindSql("updateUser").Replace(testV))
    testV.Id = -1
    t.Logf("deleteUser %s\n", config.FindSql("deleteUser").Replace(testV))
}
