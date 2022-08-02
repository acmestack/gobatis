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

package test

import (
	"github.com/acmestack/gobatis"
	"github.com/acmestack/gobatis/logging"
	"github.com/acmestack/gobatis/parsing/xml"
	"strings"
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
		Username: "testuser",
		Password: "testpw",
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
		Username: "testuser",
		Password: "testpw",
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

func TestXmlDynamicChoose(t *testing.T) {
	src := `SELECT * FROM BLOG WHERE state = 'ACTIVE'
<choose>
    <when test="{first} != nil">
      AND first = #{first}
    </when>
    <when test="{second} != nil">
      AND second = #{second}
    </when>
    <otherwise>
      AND third = 1
    </otherwise>
  </choose>
    `
	logging.SetLevel(logging.DEBUG)
	m, err := xml.ParseDynamic(src, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Run("choose first", func(t *testing.T) {
		params := map[string]interface{}{
			"first": "first",
		}
		ret := m.Replace(params)
		t.Logf("arg first : %s\n", ret)

		if strings.Index(ret, "AND first = #{first}") == -1 {
			t.FailNow()
		}
	})

	t.Run("choose second", func(t *testing.T) {
		params := map[string]interface{}{
			"second": "second",
		}
		ret := m.Replace(params)
		t.Logf("arg second : %s\n", ret)

		if strings.Index(ret, "AND second = #{second}") == -1 {
			t.FailNow()
		}
	})

	t.Run("both", func(t *testing.T) {
		params := map[string]interface{}{
			"first":  "first",
			"second": "second",
		}
		ret := m.Replace(params)
		t.Logf("arg both : %s\n", ret)

		if strings.Index(ret, "AND first = #{first}") == -1 {
			t.FailNow()
		}
	})

	t.Run("none", func(t *testing.T) {
		ret := m.Replace()
		t.Logf("arg none : %s\n", ret)

		if strings.Index(ret, "AND third = 1") == -1 {
			t.FailNow()
		}
	})
}

func TestXmlDynamicWhere(t *testing.T) {
	src := `SELECT * FROM PERSON
        <where>
            <choose>
                <when test="{first} != nil">
                    AND first = #{first}
                </when>
                <when test="{second} != nil">
                    AND second = #{second}
                </when>
                <otherwise>
                    AND third = 1
                </otherwise>
            </choose>
        </where>
`
	logging.SetLevel(logging.DEBUG)
	m, err := xml.ParseDynamic(src, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Run("choose first", func(t *testing.T) {
		params := map[string]interface{}{
			"first": "first",
		}
		ret := m.Replace(params)
		t.Logf("arg first : %s\n", ret)

		if strings.Index(ret, "where first = #{first}") == -1 {
			t.FailNow()
		}
	})

	t.Run("choose second", func(t *testing.T) {
		params := map[string]interface{}{
			"second": "second",
		}
		ret := m.Replace(params)
		t.Logf("arg second : %s\n", ret)

		if strings.Index(ret, "where second = #{second}") == -1 {
			t.FailNow()
		}
	})

	t.Run("both", func(t *testing.T) {
		params := map[string]interface{}{
			"first":  "first",
			"second": "second",
		}
		ret := m.Replace(params)
		t.Logf("arg both : %s\n", ret)

		if strings.Index(ret, "where first = #{first}") == -1 {
			t.FailNow()
		}
	})

	t.Run("none", func(t *testing.T) {
		ret := m.Replace()
		t.Logf("arg none : %s\n", ret)

		if strings.Index(ret, "where third = 1") == -1 {
			t.FailNow()
		}
	})
}

func TestXmlDynamicForeach(t *testing.T) {
	src := `SELECT * FROM PERSON WHERE IN
        <foreach item="item" index="index" collection="{0}"
            open="(" separator="," close=")">
            #{item}
        </foreach>
`
	logging.SetLevel(logging.DEBUG)
	m, err := xml.ParseDynamic(src, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Run("foreach first", func(t *testing.T) {
		params := []string{
			"first", "second",
		}
		ret, err := m.ParseMetadata("mysql", params)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("arg first : %s\n", ret)

	})
}

func TestXmlDynamicForeach2(t *testing.T) {
	src := `INSERT INTO TEST_TABLE(id, username, password) VALUES
        <foreach item="item" index="index" collection="{0}"
            open="(" separator="," close=")">
            #{item}
        </foreach>
`
	logging.SetLevel(logging.DEBUG)
	m, err := xml.ParseDynamic(src, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Run("foreach first", func(t *testing.T) {
		params := []string{
			"first", "second",
		}
		ret, err := m.ParseMetadata("mysql", params)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("arg first : %s\n", ret)

	})
}

func TestXmlDynamicForeach3(t *testing.T) {
	src := `INSERT INTO TEST_TABLE(id, username, password) VALUES
        <foreach item="item" index="index" collection="{0}"
            open="" separator="," close="">
            (#{item.testParseStruct.Username}, #{item.testParseStruct.Password})
        </foreach>
`
	logging.SetLevel(logging.DEBUG)
	m, err := xml.ParseDynamic(src, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Run("foreach first", func(t *testing.T) {
		params := []testParseStruct{
			{
				Username: "user1",
				Password: "pw1",
			},
			{
				Username: "user2",
				Password: "pw2",
			},
		}
		ret, err := m.ParseMetadata("mysql", params)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("arg first : %s\n", ret)

	})
}

func TestXmlDynamicIfForeach(t *testing.T) {
	src := `INSERT INTO TEST_TABLE(id, username, password) VALUES
		<if test="{0} != nil">
			<foreach item="item" index="index" collection="{0}"
				open="" separator="," close="">
				(#{item.testParseStruct.Username}, #{item.testParseStruct.Password})
			</foreach>
		</if>
`
	logging.SetLevel(logging.DEBUG)
	m, err := xml.ParseDynamic(src, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Run("foreach first", func(t *testing.T) {
		params := []testParseStruct{
			{
				Username: "user1",
				Password: "pw1",
			},
			{
				Username: "user2",
				Password: "pw2",
			},
		}
		ret, err := m.ParseMetadata("mysql", params)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("arg first : %s\n", ret)

	})
}
