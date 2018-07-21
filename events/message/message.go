

package message

import (
	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/core/types"
)

const (
	TOPIC_SAVE_BLOCK_COMPLETE       = "svblkcmp"
	TOPIC_NEW_INVENTORY             = "newinv"
	TOPIC_NODE_DISCONNECT           = "noddis"
	TOPIC_NODE_CONSENSUS_DISCONNECT = "nodcnsdis"
	TOPIC_SMART_CODE_EVENT          = "scevt"
)

type SaveBlockCompleteMsg struct {
	Block *types.Block
}

type NewInventoryMsg struct {
	Inventory *common.Inventory
}

type SmartCodeEventMsg struct {
	Event *types.SmartCodeEvent
}
