

package context

import (
	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/smartcontract/event"
)

// ContextRef is a interface of smart context
// when need call a contract, push current context to smart contract contexts
// when execute smart contract finish, pop current context from smart contract contexts
// when need to check authorization, use CheckWitness
// when smart contract execute trigger event, use PushNotifications push it to smart contract notifications
// when need to invoke a smart contract, use AppCall to invoke it
type ContextRef interface {
	PushContext(context *Context)
	CurrentContext() *Context
	CallingContext() *Context
	EntryContext() *Context
	PopContext()
	CheckWitness(address common.Address) bool
	PushNotifications(notifications []*event.NotifyEventInfo)
	NewExecuteEngine(code []byte) (Engine, error)
	CheckUseGas(gas uint64) bool
	CheckExecStep() bool
}

type Engine interface {
	Invoke() (interface{}, error)
}

// Context describe smart contract execute context struct
type Context struct {
	ContractAddress common.Address
	Code            []byte
}
