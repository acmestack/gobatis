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

func NewSingleSource(factory Factory) *SingleSource {
	return &SingleSource{fac: factory}
}

func (singleDs *SingleSource) Bind(action string, weight int, factory Factory) {
	singleDs.fac = factory
}

func (singleDs *SingleSource) Select(action string) Factory {
	return singleDs.fac
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

func (multiDs *DefaultMultiSource) Bind(action string, weight int, factory Factory) {
	if action == "" {
		action = DefaultGroup
	}

	if v, ok := multiDs.actionMaps[action]; ok {
		v.Add(weight, factory)
	} else {
		if f, ok := multiDs.factoryMaps[factory]; ok {
			multiDs.actionMaps[action] = f
			multiDs.factoryMaps[factory] = f
		} else {
			newlyMds := loadbalance.Create(multiDs.lbType)
			newlyMds.Add(weight, factory)
			multiDs.actionMaps[action] = newlyMds
			multiDs.factoryMaps[factory] = newlyMds
		}
	}
}

func (multiDs *DefaultMultiSource) Select(action string) Factory {
	if v, ok := multiDs.actionMaps[action]; ok {
		f := v.Select(nil)
		if f != nil {
			return f.(Factory)
		}
	}
	return nil
}
