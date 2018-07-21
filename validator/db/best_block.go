

package db

import (
	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/core/types"
)

type BestBlock struct {
	Height uint32
	Hash   common.Uint256
}

type BestStateProvider interface {
	GetBestBlock() (BestBlock, error)
	GetBestHeader() (*types.Header, error)
}
