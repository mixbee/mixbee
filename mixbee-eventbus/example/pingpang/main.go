

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
package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/mixbee/mixbee/mixbee-eventbus/actor"
)

type ping struct{ val int }
type pingActor struct{}

var start, end int64

func (state *pingActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		fmt.Println("Started, initialize actor here")
	case *actor.Stopping:
		fmt.Println("Stopping, actor is about shut down")
	case *actor.Restarting:
		fmt.Println("Restarting, actor is about restart")
	case *ping:
		val := msg.val
		if val < 10000000 {
			context.Sender().Request(&ping{val: val + 1}, context.Self())
		} else {
			end = time.Now().UnixNano()
			fmt.Printf("%s end %d\n", context.Self().Id, end)
		}
	}
}
func main() {
	fmt.Printf("test pingpang")
	runtime.GOMAXPROCS(runtime.NumCPU())
	props := actor.FromProducer(func() actor.Actor { return &pingActor{} })
	actora := actor.Spawn(props)
	actorb := actor.Spawn(props)
	start = time.Now().UnixNano()
	fmt.Printf("begin time %d\n", start)
	actora.Request(&ping{val: 1}, actorb)

	time.Sleep(10 * time.Second)
	fmt.Println((end - start) / 1000000)
	actora.Stop()
	actorb.Stop()
}
