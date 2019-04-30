/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
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
