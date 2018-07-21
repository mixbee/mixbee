

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
	"time"

	"github.com/mixbee/mixbee/mixbee-eventbus/example/serverclient/client"
	"github.com/mixbee/mixbee/mixbee-eventbus/example/serverclient/server"
)

func main() {
	server := &server.Server{}
	client := &client.Client{}
	serverPID := server.Start()
	result, err := client.SyncCall(serverPID)
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	fmt.Println(result)
	fmt.Println("###################################")

	clientPID := client.AsyncCall(serverPID)

	time.Sleep(1 * time.Second)
	server.Stop(serverPID)
	client.Stop(clientPID)
}
