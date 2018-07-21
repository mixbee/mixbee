

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

type Ball struct {
	val int
}

var start, end int64

//func Benchmark_Division1(b *testing.B){
func main() {
	fmt.Printf("test performance")
	runtime.GOMAXPROCS(4)
	times := 10000000
	props := actor.FromFunc(func(context actor.Context) {
		switch msg := context.Message().(type) {

		case Ball:
			val := msg.val
			if val < times {
				context.Sender().Request(Ball{val: val + 1}, context.Self())
			} else {
				end = time.Now().UnixNano()
				fmt.Printf("end at time %d\n", end)
			}
		default:
		}
	})
	playerA, _ := actor.SpawnNamed(props, "playerA")
	playerB, _ := actor.SpawnNamed(props, "playerB")
	start = time.Now().UnixNano()
	fmt.Println("start at time:", start)
	playerA.Request(Ball{val: 1}, playerB)
	time.Sleep(10000 * time.Millisecond)
	fmt.Printf("run time:%d     elapsed time:%d ms", times, (end-start)/1000000)
}
