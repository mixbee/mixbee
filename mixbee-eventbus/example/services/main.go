

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
	"strconv"
	"time"

	"github.com/mixbee/mixbee/mixbee-eventbus/actor"
	msg "github.com/mixbee/mixbee/mixbee-eventbus/example/services/messages"
	"github.com/mixbee/mixbee/mixbee-eventbus/example/services/servicea"
	"github.com/mixbee/mixbee/mixbee-eventbus/example/services/serviceb"
)

func main() {
	sva := actor.FromProducer(func() actor.Actor { return &servicea.ServiceA{} })
	svb := actor.FromProducer(func() actor.Actor { return &serviceb.ServiceB{} })

	pipA, _ := actor.SpawnNamed(sva, "serviceA")
	pipB, _ := actor.SpawnNamed(svb, "serviceB")

	pipA.Request(&msg.ServiceARequest{"TEST A"}, pipB)

	pipB.Request(&msg.ServiceBRequest{"TEST B"}, pipA)
	time.Sleep(2 * time.Second)

	f := pipA.RequestFuture(1, 50*time.Microsecond)
	result, err := f.Result()
	if err != nil {
		fmt.Println("errors:", err.Error())
	}
	fmt.Println("get sync call result :" + strconv.Itoa(result.(int)))

}
