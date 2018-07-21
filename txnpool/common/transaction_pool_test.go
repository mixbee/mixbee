

package common

import (
	"bytes"
	"encoding/hex"
	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/common/log"
	"github.com/mixbee/mixbee/core/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	txn *types.Transaction
)

func init() {
	log.Init(log.PATH, log.Stdout)

	txn = &types.Transaction{
		Version: 0,
		TxType:  types.Bookkeeper,
		Payload: nil,
	}

	tempStr := "3369930accc1ddd067245e8edadcd9bea207ba5e1753ac18a51df77a343bfe92"
	hex, _ := hex.DecodeString(tempStr)
	var hash common.Uint256
	hash.Deserialize(bytes.NewReader(hex))
	txn.SetHash(hash)
}

func TestTxPool(t *testing.T) {
	txPool := &TXPool{}
	txPool.Init()

	txEntry := &TXEntry{
		Tx:    txn,
		Attrs: []*TXAttr{},
	}

	ret := txPool.AddTxList(txEntry)
	if ret == false {
		t.Error("Failed to add tx to the pool")
		return
	}

	ret = txPool.AddTxList(txEntry)
	if ret == true {
		t.Error("Failed to add tx to the pool")
		return
	}

	txList, oldTxList := txPool.GetTxPool(true, 0)
	for _, v := range txList {
		assert.NotNil(t, v)
	}

	for _, v := range oldTxList {
		assert.NotNil(t, v)
	}

	entry := txPool.GetTransaction(txn.Hash())
	if entry == nil {
		t.Error("Failed to get the transaction")
		return
	}

	assert.Equal(t, txn.Hash(), entry.Hash())

	status := txPool.GetTxStatus(txn.Hash())
	if status == nil {
		t.Error("failed to get the status")
		return
	}

	assert.Equal(t, txn.Hash(), status.Hash)

	count := txPool.GetTransactionCount()
	assert.Equal(t, count, 1)

	err := txPool.CleanTransactionList([]*types.Transaction{txn})
	if err != nil {
		t.Error("Failed to clean transaction list")
		return
	}
}
