# Gobatis

## 介绍

Gobatis是一个golang的ORM框架，类似Java的Mybatis。支持直接执行sql语句以及简单的动态sql。

建议配合[gobatis-cmd](https://github.com/xfali/gobatis-cmd)自动代码、sql生成工具使用。

## 使用


### 1、配置数据库，获得SessionManager

```
func InitDB() *runner.SessionManager {
    fac := factory.DefaultFactory{
        Host:     "localhost",
        Port:     3306,
        DBName:   "test",
        Username: "root",
        Password: "123",
        Charset:  "utf8",

        MaxConn:     1000,
        MaxIdleConn: 500,

        Log: logging.DefaultLogf,
    }
    fac.Init()
    return runner.NewSessionManager(&fac)
}
```

### 2、定义Model

使用tag（"xfield"）定义struct，tag指定数据库表中的column name。

```
type TestTable struct {
    //指定table name
    TestTable gobatis.ModelName "test_table"
    //指定表字段id
    Id        int64             `xfield:"id"`
    //指定表字段username
    Username  string            `xfield:"username"`
    //指定表字段password
    Password  string            `xfield:"password"`
}
```

### 3、注册Model

注意：只有注册后的Model才能正确的序列化和反序列化；

Model注册后，Model的切片也能正确序列化和反序列化。
```
func init() {
    var model TestTable
    config.RegisterModel(&model)
}
```

### 4、调用

```
func Run() {
    //初始化db并获得Session Manager
    mgr := InitDB()
    
    //获得session
    session := mgr.NewSession()
    
    ret := TestTable{}
    
    //使用session查询
    session.Select("select * from test_table where id = ${0}").Param(100).Result(&ret)
    
    fmt.printf("%v\n", ret)
}
```

### 5、说明

1. ${}表示直接替换，#{}防止sql注入
2. 与Mybatis类似，语句中${0}、${1}、${2}...${n} 对应的是Param方法中对应的不定参数，最终替换和调用底层Driver
3. Param方法接受简单类型的不定参数（string、int、time、float等），如果参数仅为1个时，可以传递struct，底层自动解析struct获得参数，用法为：

```
param := TestTable{Username:"test_user"}
ret := TestTable{}
session.Select("select * from test_table where username = #{username}").Param(param).Result(&ret)
```

### 6、事务

使用
```
    mgr.NewSession().Tx(func(session *runner.RunnerSession) bool {
        ret := 0
        session.Insert("insert_id").Param(testV).Result(&ret)
        
        t.Logf("ret %d\n", ret)
        
        session.Select("select_id").Param().Result(&testList)
        
        for _, v := range  testList {
            t.Logf("data: %v", v)
        }
        //commit
        return true
    })
```
1. 当参数的func返回true，则提交
2. 当参数的func返回false，则回滚
3. 当参数的func内抛出panic，则回滚

### 7、xml

gobatis支持xml的sql解析及动态sql

1. 注册xml

```
config.RegisterMapperData([]byte(main_xml))
```

或
    
```
config.RegisterMapperFile(filePath)
```

2. xml示例

```
<mapper namespace="test_package.TestTable">
    <sql id="columns_id">id,username,password,update_time</sql>

    <select id="selectTestTable">
        SELECT <include refid="columns_id"> </include> FROM test_table
        <where>
            <if test="id != -1">AND id = #{id} </if>
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
            <if test="id != -1"> id = #{id} </if>
            <if test="username != nil"> username = #{username} </if>
            <if test="password != nil"> password = #{password} </if>
            <if test="update_time != nil"> update_time = #{update_time} </if>
        </set>
        WHERE id = #{id}
    </update>

    <delete id="deleteTestTable">
        DELETE FROM test_table
        <where>
            <if test="id != -1">AND id = #{id} </if>
            <if test="username != nil">AND username = #{username} </if>
            <if test="password != nil">AND password = #{password} </if>
            <if test="update_time != nil">AND update_time = #{update_time} </if>
        </where>
    </delete>
</mapper>
```