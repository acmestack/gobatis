/*
 * Copyright (c) 2022, AcmeStack
 * All rights reserved.
 *
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

var test_xml = `
<mapper namespace="test.UserInfoMapper">
    <select
            id="selectPerson"
            parameterType="int"
            parameterMap="deprecated"
            resultType="hashmap"
            resultMap="personResultMap"
            flushCache="false"
            useCache="true"
            timeout="10000"
            fetchSize="256"
            statementType="PREPARED"
            resultSetType="FORWARD_ONLY">
        SELECT * FROM PERSON WHERE ID = #{id}
    </select>

    <select
            id="selectCount"
            resultSetType="FORWARD_ONLY">
        SELECT count()* FROM PERSON
        <where>
            <if test="#{id} != nil">
                name = #{name}
            </if>
        </where>
        and id &lt; #{id}
    </select>

    <select
            id="selectCount"
            resultSetType="FORWARD_ONLY">
        SELECT <include refid="select_field"></include> FROM PERSON
        <where>
            <if test="#{id} != nil">
                and id = #{id}
            </if>
        </where>
        and name = #{name}
    </select>

    <insert
            id="insertAuthor"
            parameterType="domain.blog.Author"
            flushCache="true"
            statementType="PREPARED"
            keyProperty=""
            keyColumn=""
            useGeneratedKeys=""
            timeout="20">
        insert into Author (id,username,password,email,bio)
        values (#{id},#{username},#{password},#{email},#{bio})
    </insert>

    <sql id="select_field">
        id, username, password
    </sql>

    <update
            id="updateAuthor"
            parameterType="domain.blog.Author"
            flushCache="true"
            statementType="PREPARED"
            timeout="20">
        update Author set
        username = #{username},
        password = #{password},
        email = #{email},
        bio = #{bio}
        where id = #{id}
    </update>

    <delete
            id="deleteAuthor"
            parameterType="domain.blog.Author"
            flushCache="true"
            statementType="PREPARED"
            timeout="20">
        delete from Author where id = #{id}
    </delete>

    <sql id="fromClause">
        /*通过${}取调用者传递的属性值*/
        from ${tableName}
    </sql>
</mapper>

`
