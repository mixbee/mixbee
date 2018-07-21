

package event

import (
	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/vm/neovm/types"
)

const (
	CONTRACT_STATE_FAIL    byte = 0
	CONTRACT_STATE_SUCCESS byte = 1
)

// NotifyEventArgs describe smart contract event notify arguments struct
type NotifyEventArgs struct {
	ContractAddress common.Address
	States          types.StackItems
}

// NotifyEventInfo describe smart contract event notify info struct
type NotifyEventInfo struct {
	ContractAddress common.Address
	States          interface{}
}

type ExecuteNotify struct {
	TxHash      common.Uint256
	State       byte
	GasConsumed uint64
	Notify      []*NotifyEventInfo
}
