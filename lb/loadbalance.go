// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description: 

package lb

import (
	"math/rand"
)

type LoadBalance interface {
	Select() interface{}
}

type BaseLoadBalance struct {
	invokers []interface{}
}

func (lb *BaseLoadBalance) AddFactory(invokers ...interface{}) {
	lb.invokers = append(lb.invokers, invokers...)
}

type RandomLoadBalance struct {
	BaseLoadBalance
	rand rand.Rand
}

func (lb *RandomLoadBalance) Select() interface{} {
	size := len(lb.invokers)
	if size == 0 {
		return nil
	}

	return lb.invokers[lb.rand.Intn(size-1)]
}

type RoundRobbinLoadBalance struct {
	BaseLoadBalance
	i int
}

func (lb *RoundRobbinLoadBalance) Select() interface{} {
	size := len(lb.invokers)
	if size == 0 {
		return nil
	}

	fac := lb.invokers[lb.i]
	lb.i = (lb.i + 1) % size
	return fac
}
