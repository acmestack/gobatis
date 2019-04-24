/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package test

import (
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
            <if test="#{id} != nil">
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
            <if test="#{0} != nil">
                id = #{0}
            </if>
        </where>
        and name = #{1}`
    m, err := xml.ParseDynamic(src, nil)
    if err != nil {
        t.Fatal(err)
    }

    t.Log(m.Replace(100, 200))
}

func TestXml0(t *testing.T) {
}
