

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
package benmarks

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/mixbee/mixbee/mixbee-eventbus/actor"
)

type ping struct{ val int }

func BenchmarkSyncTest(b *testing.B) {
	defer time.Sleep(10 * time.Microsecond)
	runtime.GOMAXPROCS(runtime.NumCPU())
	defer runtime.GOMAXPROCS(1)
	b.ReportAllocs()
	b.ResetTimer()
	props := actor.FromFunc(func(context actor.Context) {
		switch msg := context.Message().(type) {
		case *ping:
			val := msg.val
			context.Sender().Tell(&ping{val: val + 1})
		}
	})
	actora := actor.Spawn(props)
	iterations := int64(b.N)
	for i := int64(0); i < iterations; i++ {
		value := actora.RequestFuture(&ping{val: 1}, 50*time.Millisecond)
		res, err := value.Result()
		if err != nil {
			fmt.Printf("sync send msg error,%s,%d", err, res)
		}
	}
	b.StopTimer()
}
