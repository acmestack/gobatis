/*
 * Copyright (c) 2022, OpeningO
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

package factory

import (
	"github.com/xfali/loadbalance"
)

type LoadBalanceType int

const (
	LBRoundRobbin       LoadBalanceType = loadbalance.LBRoundRobbin
	LBRoundRobbinWeight LoadBalanceType = loadbalance.LBRoundRobbinWeight
	LBRandom            LoadBalanceType = loadbalance.LBRandom
	LBRandomWeight      LoadBalanceType = loadbalance.LBRandomWeight

	DefaultGroup = "default"
)

type Manager interface {
	Bind(action string, weight int, factory Factory)
	Select(action string) Factory
}

type SingleSource struct {
	fac Factory
}

func NewSingleSource(fac Factory) *SingleSource {
	return &SingleSource{fac: fac}
}

func (lb *SingleSource) Bind(action string, weight int, factory Factory) {
	lb.fac = factory
}

func (lb *SingleSource) Select(action string) Factory {
	return lb.fac
}

type DefaultMultiSource struct {
	lbType      int
	actionMaps  map[string]loadbalance.LoadBalance
	factoryMaps map[Factory]loadbalance.LoadBalance
}

func NewMultiSource(t LoadBalanceType) *DefaultMultiSource {
	return &DefaultMultiSource{
		actionMaps:  map[string]loadbalance.LoadBalance{},
		factoryMaps: map[Factory]loadbalance.LoadBalance{},
		lbType:      int(t),
	}
}

func (lb *DefaultMultiSource) Bind(action string, weight int, factory Factory) {
	if action == "" {
		action = DefaultGroup
	}

	if v, ok := lb.actionMaps[action]; ok {
		v.Add(weight, factory)
	} else {
		if f, ok := lb.factoryMaps[factory]; ok {
			lb.actionMaps[action] = f
			lb.factoryMaps[factory] = f
		} else {
			newlb := loadbalance.Create(lb.lbType)
			newlb.Add(weight, factory)
			lb.actionMaps[action] = newlb
			lb.factoryMaps[factory] = newlb
		}
	}
}

func (lb *DefaultMultiSource) Select(action string) Factory {
	if v, ok := lb.actionMaps[action]; ok {
		f := v.Select(nil)
		if f != nil {
			return f.(Factory)
		}
	}
	return nil
}
