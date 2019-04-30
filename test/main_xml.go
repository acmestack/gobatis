/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package test

var main_xml = `
<mapper namespace="test.UserInfoMapper">
    <select id="selectCount">
        SELECT count(*) FROM test_table
        <where>
            <if test="password != nil">
                password = #{password}
            </if>
        </where>
        and id &lt; #{id}
    </select>

    <sql id="test_field">
        id, username, password
    </sql>

    <select id="selectUser">
        SELECT <include refid="test_field"> </include> FROM test_table
        <where>
            <if test="{0} != nil">
                username = #{0}
            </if>
        </where>
    </select>

    <insert id="insertUser">
        INSERT INTO test_table
        VALUES(
        #{id},
        #{username},
        #{password}
        )
    </insert>

    <update id="updateUser">
        UPDATE test_table
        <set>
            <if test="id != -1"> id = #{id}, </if>
            <if test="username != nil"> username = #{username} </if>
            <if test="password != nil"> password = #{password}, </if>
        </set>
    </update>

    <delete id="deleteUser">
        DELETE FROM test_table
        <where>
            <if test="id != -1"> id = #{id} </if>
            <if test="username != nil"> AND username = #{username} </if>
            <if test="password != nil"> AND password = #{password} </if>
        </where>
    </delete>
</mapper>
`
