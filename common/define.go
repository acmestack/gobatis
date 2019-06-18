/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description:
 */

package common

//参数
//  idx:迭代数
//  bean:序列化后的值
//返回值:
//  打断迭代返回true
type IterFunc func(idx int64, bean interface{}) bool

const (
	FIELD_NAME = "xfield"
)
