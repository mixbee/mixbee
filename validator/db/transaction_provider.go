

package db

import (
	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/core/types"
)

type TransactionProvider interface {
	BestStateProvider
	ContainTransaction(hash common.Uint256) bool
	GetTransactionBytes(hash common.Uint256) ([]byte, error)
	GetTransaction(hash common.Uint256) (*types.Transaction, error)
	PersistBlock(block *types.Block) error
}
