package test_package

var CMD_TEST_XML =
`
<mapper namespace="test_package.TestTable">
    <sql id="columns_id">id,username,password,update_time</sql>

    <select id="selectTestTable">
        SELECT <include refid="columns_id"> </include> FROM test_table
        <where>
            <if test="{TestTable.id} != nil and {TestTable.id} != 0">AND id = #{TestTable.id} </if>
            <if test="{TestTable.username} != nil">AND username = #{TestTable.username} </if>
            <if test="{TestTable.password} != nil">AND password = #{TestTable.password} </if>
            <if test="{TestTable.update_time} != nil">AND update_time = #{TestTable.update_time} </if>
        </where>
    </select>

    <select id="selectTestTableCount">
        SELECT COUNT(*) FROM test_table
        <where>
            <if test="{TestTable.id} != nil and {TestTable.id} != 0">AND id = #{TestTable.id} </if>
            <if test="{TestTable.username} != nil">AND username = #{TestTable.username} </if>
            <if test="{TestTable.password} != nil">AND password = #{TestTable.password} </if>
            <if test="{TestTable.update_time} != nil">AND update_time = #{TestTable.update_time} </if>
        </where>
    </select>

    <insert id="insertTestTable">
        INSERT INTO test_table (id,username,password,update_time)
        VALUES(
        #{TestTable.id},
        #{TestTable.username},
        #{TestTable.password},
        #{TestTable.update_time}
        )
    </insert>

    <update id="updateTestTable">
        UPDATE test_table
        <set>
            <if test="{TestTable.username} != nil"> username = #{TestTable.username} </if>
            <if test="{TestTable.password} != nil"> password = #{TestTable.password} </if>
            <if test="{TestTable.update_time} != nil"> update_time = #{TestTable.update_time} </if>
        </set>
        WHERE id = #{TestTable.id}
    </update>

    <delete id="deleteTestTable">
        DELETE FROM test_table
        <where>
            <if test="{TestTable.id} != nil and {TestTable.id} != 0">AND id = #{TestTable.id} </if>
            <if test="{TestTable.username} != nil">AND username = #{TestTable.username} </if>
            <if test="{TestTable.password} != nil">AND password = #{TestTable.password} </if>
            <if test="{TestTable.update_time} != nil">AND update_time = #{TestTable.update_time} </if>
        </where>
    </delete>
</mapper>
`