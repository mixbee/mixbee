

package proc

import (
	"bytes"
	"encoding/hex"
	"testing"
	"time"

	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/core/types"
	"github.com/mixbee/mixbee/errors"
	tc "github.com/mixbee/mixbee/txnpool/common"
	vt "github.com/mixbee/mixbee/validator/types"
	"github.com/stretchr/testify/assert"
)

func TestWorker(t *testing.T) {
	t.Log("Starting worker test")
	s := NewTxPoolServer(tc.MAX_WORKER_NUM, false)
	if s == nil {
		t.Error("Test case: new tx pool server failed")
		return
	}

	worker := &txPoolWorker{}
	worker.init(tc.MAX_WORKER_NUM, s)

	s.wg.Add(1)

	go worker.start()

	/* Case 1: For the given tx, validators reply
	 * with ErrNoError code. And tx pool should
	 * return the tx with the specified hash
	 */
	worker.rcvTXCh <- txn
	time.Sleep(1 * time.Second)

	statelessRsp := &vt.CheckResponse{
		WorkerId: tc.MAX_WORKER_NUM,
		ErrCode:  errors.ErrNoError,
		Hash:     txn.Hash(),
		Type:     vt.Stateless,
		Height:   0,
	}

	statefulRsp := &vt.CheckResponse{
		WorkerId: tc.MAX_WORKER_NUM,
		ErrCode:  errors.ErrNoError,
		Hash:     txn.Hash(),
		Type:     vt.Stateful,
		Height:   0,
	}
	worker.rspCh <- statelessRsp
	worker.rspCh <- statefulRsp

	time.Sleep(1 * time.Second)
	ret := worker.server.getTransaction(txn.Hash())
	assert.NotNil(t, ret)
	assert.Equal(t, ret.Hash(), txn.Hash())
	/* Case 2: Duplicate input tx, worker should reject
	 * it with the log
	 */
	worker.rcvTXCh <- txn

	/* Case 3: For the given tx, validators reply with
	 * Error code, worker should remove the tx from the
	 * pending list with the log
	 */
	time.Sleep(1 * time.Second)
	worker.server.cleanTransactionList([]*types.Transaction{txn}, 0)

	worker.rcvTXCh <- txn
	time.Sleep(1 * time.Second)

	statelessRsp = &vt.CheckResponse{
		WorkerId: tc.MAX_WORKER_NUM,
		ErrCode:  errors.ErrUnknown,
		Hash:     txn.Hash(),
		Type:     vt.Stateless,
		Height:   0,
	}

	statefulRsp = &vt.CheckResponse{
		WorkerId: tc.MAX_WORKER_NUM,
		ErrCode:  errors.ErrUnknown,
		Hash:     txn.Hash(),
		Type:     vt.Stateful,
		Height:   0,
	}
	worker.rspCh <- statelessRsp
	worker.rspCh <- statefulRsp

	/* Case 4: valdiators reply with invalid tx hash or invalid work id,
	 * worker should reject it
	 */
	time.Sleep(2 * time.Second)
	statelessRsp = &vt.CheckResponse{
		WorkerId: tc.MAX_WORKER_NUM,
		ErrCode:  errors.ErrUnknown,
		Hash:     txn.Hash(),
		Type:     vt.Stateless,
		Height:   0,
	}

	statefulRsp = &vt.CheckResponse{
		WorkerId: tc.MAX_WORKER_NUM + 1,
		ErrCode:  errors.ErrUnknown,
		Hash:     txn.Hash(),
		Type:     vt.Stateful,
		Height:   0,
	}
	worker.rspCh <- statelessRsp
	worker.rspCh <- statefulRsp

	/* Case 5: For the given tx, response time out, worker should
	 * retry verifying it till retries exhausted, and then remove
	 * it from the pending list
	 */
	time.Sleep(2 * time.Second)
	worker.rcvTXCh <- txn

	time.Sleep(15 * time.Second)
	txStatus := worker.GetTxStatus(txn.Hash())
	t.Log(txStatus)
	assert.Nil(t, txStatus)

	/* Case 6: For the given tx, worker handle it once, if
	 * duplicate input the tx, worker should reject it with
	 * the log.
	 */
	worker.rcvTXCh <- txn
	worker.rcvTXCh <- txn

	/* Case 7: For the pending tx, worker should get the entry
	 * with the valid hash
	 */
	time.Sleep(1 * time.Second)
	txStatus = worker.GetTxStatus(txn.Hash())
	t.Log(txStatus)
	assert.NotNil(t, txStatus)
	assert.Equal(t, txStatus.Hash, txn.Hash())
	/* Case 8: Given the invalid hash, worker should return nil
	 */
	tempStr := "3369930accc1ddd067245e8edadcd9bea207ba5e1753ac18a51df77a343bfe83"
	hex, _ := hex.DecodeString(tempStr)
	var hash common.Uint256
	hash.Deserialize(bytes.NewReader(hex))
	txStatus = worker.GetTxStatus(hash)
	assert.Nil(t, txStatus)

	worker.stop()
	s.Stop()
	t.Log("Ending worker test")
}
