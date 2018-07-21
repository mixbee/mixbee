

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

package zmqremote

import "github.com/mixbee/mixbee/mixbee-eventbus/actor"

type EndpointTerminatedEvent struct {
	Address string
}

type EndpointConnectedEvent struct {
	Address string
}

type remoteWatch struct {
	Watcher *actor.PID
	Watchee *actor.PID
}

type remoteUnwatch struct {
	Watcher *actor.PID
	Watchee *actor.PID
}

type remoteDeliver struct {
	header       actor.ReadonlyMessageHeader
	message      interface{}
	target       *actor.PID
	sender       *actor.PID
	serializerID int32
}

type remoteTerminate struct {
	Watcher *actor.PID
	Watchee *actor.PID
}

type JsonMessage struct {
	TypeName string
	Json     string
}

var (
	stopMessage interface{} = &actor.Stop{}
)

var (
	ActorPidRespErr         interface{} = &ActorPidResponse{StatusCode: ResponseStatusCodeERROR.ToInt32()}
	ActorPidRespTimeout     interface{} = &ActorPidResponse{StatusCode: ResponseStatusCodeTIMEOUT.ToInt32()}
	ActorPidRespUnavailable interface{} = &ActorPidResponse{StatusCode: ResponseStatusCodeUNAVAILABLE.ToInt32()}
)
