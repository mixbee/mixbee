

/***************************************************
Copyright 2016 https://github.com/AsynkronIT/protoactor-go

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*****************************************************/
package router

import (
	"github.com/mixbee/mixbee/mixbee-eventbus/actor"
)

type RouterType int

const (
	GroupRouterType RouterType = iota
	PoolRouterType
)

type RouterConfig interface {
	RouterType() RouterType
	OnStarted(context actor.Context, props *actor.Props, router Interface)
	CreateRouterState() Interface
}

type GroupRouter struct {
	Routees *actor.PIDSet
}

type PoolRouter struct {
	PoolSize int
}

func (config *GroupRouter) OnStarted(context actor.Context, props *actor.Props, router Interface) {
	config.Routees.ForEach(func(i int, pid actor.PID) {
		context.Watch(&pid)
	})
	router.SetRoutees(config.Routees)
}

func (config *GroupRouter) RouterType() RouterType {
	return GroupRouterType
}

func (config *PoolRouter) OnStarted(context actor.Context, props *actor.Props, router Interface) {
	var routees actor.PIDSet
	for i := 0; i < config.PoolSize; i++ {
		routees.Add(context.Spawn(props))
	}
	router.SetRoutees(&routees)
}

func (config *PoolRouter) RouterType() RouterType {
	return PoolRouterType
}

func spawner(config RouterConfig) actor.SpawnFunc {
	return func(id string, props *actor.Props, parent *actor.PID) (*actor.PID, error) {
		return spawn(id, config, props, parent)
	}
}
