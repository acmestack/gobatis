/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package gobatis

type SqlRunner int

func (this *SqlRunner) Select(sqlId string) *SqlRunner {

    return this
}

func (this *SqlRunner) Update(sqlId string) *SqlRunner {

    return this
}

func (this *SqlRunner) Delete(sqlId string) *SqlRunner {

    return this
}

func (this *SqlRunner) Insert(sqlId string) *SqlRunner {

    return this
}

func (this *SqlRunner) ParamType(params interface{}) *SqlRunner {

    return this
}

func (this *SqlRunner) Params(params ...interface{}) *SqlRunner {

    return this
}

func (this *SqlRunner) ResultType(bean interface{}) *SqlRunner {

    return this
}
