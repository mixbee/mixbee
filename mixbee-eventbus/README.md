[中文版](./README_cn.md)

# The Go Language Implementation of Actor Model

[inter-process communication](#inter-process-communication)

[inter-node communication](#inter-node-communication)

[ONT's Signature Verification Test](#signature-verification-test)

# About Actor Model

Actor is a model of parallel computation model in computer science that treats "actors" as the the universal primitives of concurrent computation. In response to a message that it receives, an actor can make local decisions: make local decisions, create more actor, send more messages, and determine how to respond to the next message received.

The actor model adopts the philosophy that everything is an actor. This is similar to the everthing is an object philosophy used by some object-oriented programming languages. Although softwares including object-oriented language in sequence, Actor model executes in parellel in essence.

All actor has one (and only one) Mailbox. Mailbox is like a small message queue. Once the sender send message, the message will be pushed into the Mailbox. The sequence of push is determined by the sequence of sending. Mailbox also has many type of implementation, the default is FIFO. However, the implementation could be different according to the priority of pop.

![actor](/resources/actor.png)

## Actor VS Channel

![actor](/resources/actors.png)


![channel](/resources/channel.png)

## Advantages of Actor Model：

1. The volume of actor's Mailbox is unlimited, which won't interrupt the processing to writing.

2. All messages of each single actor share the same mailbox(channel)

3. Actor do not need to care about the message writer, and could decouple the logic between each module.

4. Actor can be deployed on different nodes.


## Disadvantage of Actor Model

1. Since the Actor model is designed to be an asynchronous model, the efficiency of synchronization is not very high.



## Create Actors

Props has provided the basis of declaration of how to create Actors. The following example defines the Actor Propos by defining the declaration of the function processes messge.

```go
var props Props = actor.FromFunc(func(c Context) {
	// process messages
})
```

Besides, The interface of Actor could be implemented by creating a structure and defining a Recive function.

```go
type MyActor struct {}

func (a *MyActor) Receive(c Context) {
	// process messages
}

var props Props = actor.FromProducer(func() Actor { return &MyActor{} })
```

Spawn and SpawnNamed make use of the given props to create the execution instance of Actor. Once the Actor is started, it begins to process the received message. Use the unique name specified by ststem to start the actor like: 
```go
pid := actor.Spawn(props)
```
The return value is an unique PID. You could start Actor if you want to name PID on your own.

Once an actor is started, a new email address will be created and related to the PID. The messages will be sent to the address, and processed bht actor.


## Process Messages

Actor Processes messages by Receive function which is defined: 

```go
Receive(c actor.Context)
```
The system will make sure that the function would only be called synchronously. Hence user don't need to figure out any additional protection measure.

## Communicate with other actors

PID is the main interface to send actors messages. And PID.Tell fucntion is used to send messages to the PID asynchronously.

```go
pid.Tell("Hello World")
```
According to different business requirement, the communication between actors could be carried out synchronously or asynchronously. And Actors align with PID whenever they communicate.

When PID.Request or PID.RequestFutre is used to send messages, the actor receiving messages will response to the sender by Context.Sender function, which returns the PID of sender.

In terms of synchronous communication, Actor uses Future to implment it. Actor will wait for the result before carry out the next step.

User could use RequestFuture function to send message to actor and wait for result. The function will return a Future:

```go
f := actor.RequestFuture(pid,"Hello", 50 * time.Millisecond)
res, err := f.Result() // waits for pid to reply */
```

# Inter-process Communication  
## Performance
### Asynchronous call

Protoactor-go can currently pass about 2000,000 messages per second between 2 actors, and can guarantee the order of those messages. 

```text
/app/go/bin/go build -o "/tmp/Build performanceTest.go and rungo" 
/app/gopath/src/github.com/ontio/ontology-eventbus/example/performanceTest.go
start at time: 1516953710985385134
end at time 1516953716291953904
run time:10000000     elapsed time:5306 ms
```

### Serial synchronization call

Protoactor-go can now pass more than 500,000 messages per second between client and server in a serial synchronous call.

```text
goos: linux
goarch: amd64
pkg: github.com/ontio/ontology-eventbus/example/benchmark
benchmark                  iter                 time/iter          bytes alloc          allocs
---------                  ----                ---------           -----------          ------
BenchmarkSyncTest-4   	 1000000	      1967 ns/op	     432 B/op	      13 allocs/op
testing: BenchmarkSyncTest-4 left GOMAXPROCS set to 1
BenchmarkSyncTest-4   	 1000000	      1987 ns/op	     432 B/op	      13 allocs/op
testing: BenchmarkSyncTest-4 left GOMAXPROCS set to 1
BenchmarkSyncTest-4   	 1000000	      1952 ns/op	     432 B/op	      13 allocs/op
testing: BenchmarkSyncTest-4 left GOMAXPROCS set to 1
BenchmarkSyncTest-4   	 1000000	      1975 ns/op	     432 B/op	      13 allocs/op
testing: BenchmarkSyncTest-4 left GOMAXPROCS set to 1
BenchmarkSyncTest-4   	 1000000	      1987 ns/op	     432 B/op	      13 allocs/op
testing: BenchmarkSyncTest-4 left GOMAXPROCS set to 1
PASS
ok  	github.com/ontio/ontology-eventbus/example/benchmarks	10.984s

```

## Hello world

```go
type Hello struct{ Who string }
type HelloActor struct{}

func (state *HelloActor) Receive(context actor.Context) {
    switch msg := context.Message().(type) {
    case Hello:
        fmt.Printf("Hello %v\n", msg.Who)
    }
}

func main() {
    props := actor.FromProducer(func() actor.Actor { return &HelloActor{} })
    pid := actor.Spawn(props)
    pid.Tell(Hello{Who: "Roger"})
    console.ReadLine()
}
```


## Two actors communicates each other

This example describes how to perform asynchronous communication between two actors. It mainly defines the behavior of the actor after receiving messages(Receive), including the processing method and the actor to which the processed message is sent. The asynchronous communication ensures the utilization of the actor.

```go
type ping struct{ val int }
type pingActor struct{}

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
			end := time.Now().UnixNano()
			fmt.Printf("%s end %d\n", context.Self().Id, end)
		}
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	props := actor.FromProducer(func() actor.Actor { return &pingActor{} })
	actora := actor.Spawn(props)
	actorb := actor.Spawn(props)
	fmt.Printf("begin time %d\n", time.Now().UnixNano())
	actora.Request(&ping{val: 1}, actorb)
	time.Sleep(10 * time.Second)
	actora.Stop()
	actorb.Stop()
	console.ReadLine()
}
```

## Server/Client synchronization call

This example mainly describes how to communicate with the actor (server) synchronously. The client sends the request message to the actor and waits for the actor to return the result. The request may need multiple actor to cooperate and complete. Asynchronous communication in the above example is used for processing between multiple actors, and the final processing result will be returned to the client.



#### message.go
```go
type Request struct {
	Who string
}

type Response struct {
	Welcome string
}
```

#### server.go
```go
type Server struct {}

func (server *Server) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		fmt.Println("Started, initialize server actor here")
	case *actor.Stopping:
		fmt.Println("Stopping, actor is about shut down")
	case *actor.Restarting:
		fmt.Println("Restarting, actor is about restart")
	case *message.Request:
		fmt.Println("Receive message", msg.Who)
		context.Sender().Request(&message.Response{Welcome: "Welcome!"}, context.Self())
	}
}

func (server *Server) Start() *actor.PID{
	props := actor.FromProducer(func() actor.Actor { return &Server{} })
	pid := actor.Spawn(props)
	return pid
}

func (server *Server) Stop(pid *actor.PID) {
	pid.Stop()
}
```

#### client.go
```go
type Client struct {}


//Call the server synchronously
func (client *Client) SyncCall(serverPID *actor.PID) (interface{}, error) {
	future := serverPID.RequestFuture(&message.Request{Who: "Ontology"}, 10*time.Second)
	result, err := future.Result()
	return result, err
}
```

#### main.go
```go
func main() {
	server := &server.Server{}
	client := &client.Client{}
	serverPID := server.Start()
	result, err := client.SyncCall(serverPID)
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	fmt.Println(result)
}
```
## EventHub
Actor can perform broadcast and subscribe operations through EventHub, support ALL, ROUNDROBIN, RANDOM broadcast mode

### Example
```go
package main

import (
	"github.com/ontio/ontology-eventbus/eventhub"
	"fmt"
	"github.com/ontio/ontology-eventbus/actor"

	"time"
)


type PubMessage struct{
	message string
}

type ResponseMessage struct{
	message string
}

func main() {

	eh:= eventhub.GlobalEventHub
	subprops := actor.FromFunc(func(context actor.Context) {
		switch msg := context.Message().(type) {

		case PubMessage:
			fmt.Println(context.Self().Id + " get message "+msg.message)
			context.Sender().Request(ResponseMessage{"response message from "+context.Self().Id },context.Self())

		default:

		}
	})

	pubprops := actor.FromFunc(func(context actor.Context) {
		switch msg := context.Message().(type) {

		case ResponseMessage:
			fmt.Println(context.Self().Id + " get message "+msg.message)
			//context.Sender().Request(ResponseMessage{"response message from "+context.Self().Id },context.Self())

		default:
			//fmt.Println("unknown message type")
		}
	})


	publisher, _ := actor.SpawnNamed(pubprops, "publisher")
	sub1, _ := actor.SpawnNamed(subprops, "sub1")
	sub2, _ := actor.SpawnNamed(subprops, "sub2")
	sub3, _ := actor.SpawnNamed(subprops, "sub3")

	topic:= "TEST"
	eh.Subscribe(topic,sub1)
	eh.Subscribe(topic,sub2)
	eh.Subscribe(topic,sub3)

	event := eventhub.Event{Publisher:publisher,Message:PubMessage{"hello fellows"},Topic:topic,Policy:eventhub.PublishPolicyAll}
	eh.Publish(event)
	time.Sleep(2*time.Second)
	fmt.Println("before unsubscribe sleeping...")
	eh.Unsubscribe(topic,sub2)
	eh.Publish(event)
	time.Sleep(2*time.Second)

	fmt.Println("random event...")
	randomevent := eventhub.Event{Publisher:publisher,Message:PubMessage{"hello fellows"},Topic:topic,Policy:eventhub.PublishPolicyRandom}
	for i:=0 ;i<10;i++{
		eh.Publish(randomevent)
	}

	time.Sleep(2*time.Second)

	fmt.Println("roundrobin event...")
	roundevent := eventhub.Event{Publisher:publisher,Message:PubMessage{"hello fellows"},Topic:topic,Policy:eventhub.PublishPolicyRoundRobin}
	for i:=0 ;i<10;i++{
		eh.Publish(roundevent)
	}

	time.Sleep(2*time.Second)
}


```

# Inter-node Communication
This project adopts two ways to implement inter-node communication, namely grpc and zeromq, corresponding to the remote and zmqremote packages in the project. During use, please select the package to be imported according to requirements (the interface is the same).

In order to use zero mq need to install libzmq [https://github.com/zeromq/libzmq]

## Performance
### Microsoft Cloud Intra Node

| mode | 256B Message Size | 512B Message Size | 1k Message Size | 10k Message Size | 100k  Message Size| 1M  Message Size| 4M Message Size | 8M   Message Size|
| :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: |
| grpc | 120000/s | 100000/s | 85000/s | 40000/s | 4600/s | 490/s | 123/s | Out of grpc default 4m limit |
| zeromq | 170000/s | 140000/s | 10000/s | 45000/s | 4900/s | 500/s | 123/s | 62/s |

## Serialization
The module uses protobuf for serialization and deserialization by default. During use, a series of protobuf message structures need to be defined. In order to integrate with the serialization methods used in the Ontology project, and reduce the system modification workload, personalized serialization and deserialization methods is currently supported:

The current approach is to define a generic system message. The message structure is: message type + message content (data after serialization). For Ontology's commonly used structures, there are currently six commonly used message types (address, Block, header, Transaction, TxAttribute, VMCode) as follows:


```text
enum MsgType {
  ADDRESS_MSG_TYPE = 0;
  BLOCK_MSG_TYPE = 1;
  HEADER_MSG_TYPE = 2;
  TX_MSG_TYPE    = 3;
  TX_ATT_MSG_TYPE = 4;
  VM_CODE_MSG_TYPE = 5;
}

message MsgData {
  MsgType msgType = 1;
  bytes data = 2;
}
```
While using, serialize the data to be transmitted in a custom manner, and then construct MsgData {msgType:xx, data:xx}, using the above enumeration definition to define msgType. After these processes ,data is custom serialized. The same is true for the received message. After receiving the message, execute the corresponding deserialization method according to msgType to deserialize the data. A simple case is as follows:

server.go
```go
func main() {
	log.Debug("test")
	runtime.GOMAXPROCS(runtime.NumCPU() * 1)
	runtime.GC()

	zmqremote.Start("127.0.0.1:8080")

	props := actor.
		FromFunc(
			func(context actor.Context) {
				switch context.Message().(type) {
				case *zmqremote.MsgData:
					switch MsgData.MsgType:
						case 0:  //Deserialization MsgData.Data
						case 1:  //Deserialization MsgData.Data
						case 2:  //Deserialization MMsgData.Data
						case 3:  //Deserialization MMsgData.Data
						case 4:  //Deserialization MMsgData.Data
						case 5:  //Deserialization MMsgData.Data
					context.Sender().Tell(&zmqremote.MsgData{MsgType: 1, Data: []byte("123")})
				}
			}).
		WithMailbox(mailbox.Bounded(1000000))

	pid, _ := actor.SpawnNamed(props, "remote")
	fmt.Println(pid)

	for {
		time.Sleep(1 * time.Second)
	}
}
```

client.go
```go
func main() {

	log.Debug("test")
	runtime.GOMAXPROCS(runtime.NumCPU() * 1)
	runtime.GC()

	var wg sync.WaitGroup

	messageCount := 500

	zmqremote.Start("127.0.0.1:8081")

	props := actor.
		FromProducer(newLocalActor(&wg, messageCount)).
		WithMailbox(mailbox.Bounded(1000000))

	pid := actor.Spawn(props)
	fmt.Println(pid)

	remotePid := actor.NewPID("127.0.0.1:8080", "remote")
	wg.Add(1)

	start := time.Now()
	fmt.Println("Starting to send")

	message := &zmqremote.MsgData{MsgType: 1, Data: []byte("123")}
	for i := 0; i < messageCount; i++ {
		remotePid.Request(message, pid)
	}

	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("Elapsed %s", elapsed)

	x := int(float32(messageCount*2) / (float32(elapsed) / float32(time.Second)))
	fmt.Printf("Msg per sec %v", x)
}
```

## Benchmark
node2/main.go
```go
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 1)
	runtime.GC()

	remote.Start("127.0.0.1:8080")
	var sender *actor.PID
	props := actor.
		FromFunc(
			func(context actor.Context) {
				switch msg := context.Message().(type) {
				case *messages.StartRemote:
					fmt.Println("Starting")
					sender = msg.Sender
					context.Respond(&messages.Start{})
				case *messages.Ping:
					sender.Tell(&messages.Pong{})
				}
			}).
		WithMailbox(mailbox.Bounded(1000000))
	actor.SpawnNamed(props, "remote")
	for{
		time.Sleep(1 * time.Second)
	}
}
```

node1/main.go
```go
type localActor struct {
	count        int
	wgStop       *sync.WaitGroup
	messageCount int
}

func (state *localActor) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *messages.Pong:
		state.count++
		if state.count%50000 == 0 {
			fmt.Println(state.count)
		}
		if state.count == state.messageCount {
			state.wgStop.Done()
		}
	}
}

func newLocalActor(stop *sync.WaitGroup, messageCount int) actor.Producer {
	return func() actor.Actor {
		return &localActor{
			wgStop:       stop,
			messageCount: messageCount,
		}
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 1)
	runtime.GC()

	var wg sync.WaitGroup

	messageCount := 50000
	//remote.DefaultSerializerID = 1
	remote.Start("127.0.0.1:8081")

	props := actor.
		FromProducer(newLocalActor(&wg, messageCount)).
		WithMailbox(mailbox.Bounded(1000000))

	pid := actor.Spawn(props)

	remotePid := actor.NewPID("127.0.0.1:8080", "remote")
	remotePid.
		RequestFuture(&messages.StartRemote{
			Sender: pid,
		}, 5*time.Second).
		Wait()

	wg.Add(1)

	start := time.Now()
	fmt.Println("Starting to send")

	bb := bytes.NewBuffer([]byte(""))
	for i := 0; i < 2000; i++ {
		bb.WriteString("1234567890")
	}
	message := &messages.Ping{Data: bb.Bytes()}
	for i := 0; i < messageCount; i++ {
		remotePid.Tell(message)
	}

	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("Elapsed %s", elapsed)

	x := int(float32(messageCount*2) / (float32(elapsed) / float32(time.Second)))
	fmt.Printf("Msg per sec %v", x)
}
```

messages/protos.proto

Protobuf file generation command:
protoc -I=$GOPATH/src -I=$GOPATH/src/github.com/gogo/protobuf/protobuf/ --gogoslick_out=plugins=grpc:. /path/to/protos.proto

```text
syntax = "proto3";
package messages;
import "github.com/ontio/ontology-eventbus/actor/protos.proto";

message Start {}
message StartRemote {
    actor.PID Sender = 1;
}
message Ping {
    bytes Data = 1;
}
message Pong {}
```

# Signature Verification Test

The code could be found at directories: example/testRemoteCrypto and example/testSyncCrypto

The test environment comes from Microsoft Azure.

## Asynchronous Signature Verification Test
| Mode                 | 256B Message Size | 512B Message Size | 1k Message Size | 10k Message Size |
| :-:                  | :-:               | :-:               | :-:             | :-:              |
| One Machine(zeromq)  | 3666/s            | 3590/s            | 3479/s          | 2848/s           |
| Two Machines(zeromq) | 7509/s            | 7431/s            | 7204/s          | 6976/s           |

## Synchronous Signature Verification Test
| Quota                           | 256B Message Size | 512B Message Size | 1k Message Size | 10k Message Size |
| :-:                             | :-:               | :-:               | :-:             | :-:              |
| Time for Signature Verification | 0.242ms           | 0.247ms           | 0.246ms         | 0.334ms          |
| latency                         | 1.36ms            | 1.31ms            | 1.39ms          | 1.94ms           |



This Module is based on AsynkronIT/protoactor-go project, more details goes to https://github.com/AsynkronIT/protoactor-go.
