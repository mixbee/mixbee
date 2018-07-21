

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
	"math/rand"

	"github.com/mixbee/mixbee/mixbee-eventbus/actor"
)

type randomGroupRouter struct {
	GroupRouter
}

type randomPoolRouter struct {
	PoolRouter
}

type randomRouterState struct {
	routees *actor.PIDSet
	values  []actor.PID
}

func (state *randomRouterState) SetRoutees(routees *actor.PIDSet) {
	state.routees = routees
	state.values = routees.Values()
}

func (state *randomRouterState) GetRoutees() *actor.PIDSet {
	return state.routees
}

func (state *randomRouterState) RouteMessage(message interface{}) {
	l := len(state.values)
	r := rand.Intn(l)
	pid := state.values[r]
	pid.Tell(message)
}

func NewRandomPool(size int) *actor.Props {
	return actor.FromSpawnFunc(spawner(&randomPoolRouter{PoolRouter{PoolSize: size}}))
}

func NewRandomGroup(routees ...*actor.PID) *actor.Props {
	return actor.FromSpawnFunc(spawner(&randomGroupRouter{GroupRouter{Routees: actor.NewPIDSet(routees...)}}))
}

func (config *randomPoolRouter) CreateRouterState() Interface {
	return &randomRouterState{}
}

func (config *randomGroupRouter) CreateRouterState() Interface {
	return &randomRouterState{}
}
