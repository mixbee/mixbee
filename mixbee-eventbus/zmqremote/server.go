

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

import (
	"github.com/mixbee/mixbee/mixbee-eventbus/actor"
	"github.com/mixbee/mixbee/mixbee-eventbus/log"
	zmq "github.com/pebbe/zmq4"
)

var (
	edpReader *endpointReader
	conn      *zmq.Socket
)

func Start(address string) {

	actor.ProcessRegistry.RegisterAddressResolver(remoteHandler)
	actor.ProcessRegistry.Address = address

	spawnActivatorActor()
	startEndpointManager()

	edpReader = &endpointReader{}

	conn, _ = zmq.NewSocket(zmq.ROUTER)
	err := conn.Bind("tcp://" + address)
	if err != nil {
		plog.Error("failed to Bind", log.Error(err))
	}
	plog.Info("Starting Proto.Actor server", log.String("address", address))
	go func() {
		edpReader.Receive(conn)
	}()
}

func Shutdonw() {
	edpReader.suspend(true)
	stopEndpointManager()
	stopActivatorActor()
	conn.Close()
}
