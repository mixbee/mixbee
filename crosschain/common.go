package crosschain

import (
	"sync"
	"time"
	tx "github.com/mixbee/mixbee/core/types"
	"github.com/mixbee/mixbee/validator/types"
)

// txPoolWorker handles the tasks scheduled by server
type ctxPoolWorker struct {
	mu      sync.RWMutex
	workId  uint8                     // Worker ID
	rcvTXCh chan *tx.Transaction      // The channel of receive transaction
	stfTxCh chan *tx.Transaction      // The channel of txs to be re-verified stateful
	rspCh   chan *types.CheckResponse // The channel of verified response
	server  *CTXPoolServer            // The txn pool server pointer
	timer   *time.Timer               // The timer of reverifying
	stopCh  chan bool                 // stop routine
}

