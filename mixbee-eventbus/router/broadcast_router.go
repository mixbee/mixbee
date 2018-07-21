

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

import "github.com/mixbee/mixbee/mixbee-eventbus/actor"

type broadcastGroupRouter struct {
	GroupRouter
}

type broadcastPoolRouter struct {
	PoolRouter
}

type broadcastRouterState struct {
	routees *actor.PIDSet
}

func (state *broadcastRouterState) SetRoutees(routees *actor.PIDSet) {
	state.routees = routees
}

func (state *broadcastRouterState) GetRoutees() *actor.PIDSet {
	return state.routees
}

func (state *broadcastRouterState) RouteMessage(message interface{}) {
	state.routees.ForEach(func(i int, pid actor.PID) {
		pid.Tell(message)
	})
}

func NewBroadcastPool(size int) *actor.Props {
	return actor.FromSpawnFunc(spawner(&broadcastPoolRouter{PoolRouter{PoolSize: size}}))
}

func NewBroadcastGroup(routees ...*actor.PID) *actor.Props {
	return actor.FromSpawnFunc(spawner(&broadcastGroupRouter{GroupRouter{Routees: actor.NewPIDSet(routees...)}}))
}

func (config *broadcastPoolRouter) CreateRouterState() Interface {
	return &broadcastRouterState{}
}

func (config *broadcastGroupRouter) CreateRouterState() Interface {
	return &broadcastRouterState{}
}
