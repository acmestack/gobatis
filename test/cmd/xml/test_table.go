package test_package

var CMD_TEST_XML =
`
<mapper namespace="test_package.TestTable">
    <sql id="columns_id">id,username,password,update_time</sql>

    <select id="selectTestTable">
        SELECT <include refid="columns_id"> </include> FROM test_table
        <where>
            <if test="id != nil and id != 0">AND id = #{id} </if>
            <if test="username != nil">AND username = #{username} </if>
            <if test="password != nil">AND password = #{password} </if>
            <if test="update_time != nil">AND update_time = #{update_time} </if>
        </where>
    </select>

    <select id="selectTestTableCount">
        SELECT COUNT(*) FROM test_table
        <where>
            <if test="id != nil and id != 0">AND id = #{id} </if>
            <if test="username != nil">AND username = #{username} </if>
            <if test="password != nil">AND password = #{password} </if>
            <if test="update_time != nil">AND update_time = #{update_time} </if>
        </where>
    </select>

    <insert id="insertTestTable">
        INSERT INTO test_table (id,username,password,update_time)
        VALUES(
        #{id},
        #{username},
        #{password},
        #{update_time}
        )
    </insert>

    <update id="updateTestTable">
        UPDATE test_table
        <set>
            <if test="id != nil and id != 0"> id = #{id} </if>
            <if test="username != nil"> username = #{username} </if>
            <if test="password != nil"> password = #{password} </if>
            <if test="update_time != nil"> update_time = #{update_time} </if>
        </set>
        WHERE id = #{id}
    </update>

    <delete id="deleteTestTable">
        DELETE FROM test_table
        <where>
            <if test="id != nil and id != 0">AND id = #{id} </if>
            <if test="username != nil">AND username = #{username} </if>
            <if test="password != nil">AND password = #{password} </if>
            <if test="update_time != nil">AND update_time = #{update_time} </if>
        </where>
    </delete>
</mapper>
`