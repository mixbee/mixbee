
package wasmvm

import (
	"bytes"
	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/core/types"
	"github.com/mixbee/mixbee/errors"
	"github.com/mixbee/mixbee/vm/wasmvm/exec"
)

func (this *WasmVmService) blockGetCurrentHeaderHash(engine *exec.ExecutionEngine) (bool, error) {
	vm := engine.GetVM()
	vm.RestoreCtx()

	headerHash := this.Store.GetCurrentHeaderHash()
	idx, err := vm.SetPointerMemory(headerHash.ToArray())
	if err != nil {
		return false, err
	}
	vm.RestoreCtx()
	vm.PushResult(uint64(idx))
	return true, nil
}

func (this *WasmVmService) blockGetCurrentHeaderHeight(engine *exec.ExecutionEngine) (bool, error) {
	vm := engine.GetVM()
	vm.RestoreCtx()
	headerHight := this.Store.GetCurrentHeaderHeight()
	vm.RestoreCtx()
	vm.PushResult(uint64(headerHight))
	return true, nil
}

func (this *WasmVmService) blockGetCurrentBlockHash(engine *exec.ExecutionEngine) (bool, error) {
	vm := engine.GetVM()
	vm.RestoreCtx()

	bHash := this.Store.GetCurrentBlockHash()
	idx, err := vm.SetPointerMemory(bHash.ToArray())
	if err != nil {
		return false, err
	}
	vm.RestoreCtx()
	vm.PushResult(uint64(idx))
	return true, nil
}

func (this *WasmVmService) blockGetCurrentBlockHeight(engine *exec.ExecutionEngine) (bool, error) {
	vm := engine.GetVM()
	vm.RestoreCtx()
	bHight := this.Store.GetCurrentBlockHeight()
	vm.RestoreCtx()
	vm.PushResult(uint64(bHight))
	return true, nil
}

func (this *WasmVmService) blockGetTransactionByHash(engine *exec.ExecutionEngine) (bool, error) {
	vm := engine.GetVM()
	envCall := vm.GetEnvCall()
	params := envCall.GetParams()
	if len(params) != 1 {
		return false, errors.NewErr("[RuntimeLog]parameter count error ")
	}

	hashbytes, err := vm.GetPointerMemory(params[0])
	if err != nil {
		return false, err
	}

	thash, err := common.Uint256ParseFromBytes(hashbytes)
	if err != nil {
		return false, err
	}
	tx, _, err := this.Store.GetTransaction(thash)
	txbytes := tx.ToArray()
	idx, err := vm.SetPointerMemory(txbytes)
	if err != nil {
		return false, err
	}
	vm.RestoreCtx()
	vm.PushResult(uint64(idx))
	return true, nil

}

// BlockGetTransactionCount put block's transactions count to vm stack
func (this *WasmVmService) blockGetTransactionCount(engine *exec.ExecutionEngine) (bool, error) {
	vm := engine.GetVM()
	envCall := vm.GetEnvCall()
	params := envCall.GetParams()
	if len(params) != 1 {
		return false, errors.NewErr("[RuntimeLog]parameter count error ")
	}

	blockbytes, err := vm.GetPointerMemory(params[0])
	if err != nil {
		return false, err
	}
	block := &types.Block{}
	err = block.Deserialize(bytes.NewBuffer(blockbytes))
	if err != nil {
		return false, err
	}

	length := len(block.Transactions)

	vm.RestoreCtx()
	vm.PushResult(uint64(length))
	return true, nil
}

// BlockGetTransactions put block's transactions to vm stack
func (this *WasmVmService) blockGetTransactions(engine *exec.ExecutionEngine) (bool, error) {
	vm := engine.GetVM()
	envCall := vm.GetEnvCall()
	params := envCall.GetParams()
	if len(params) != 1 {
		return false, errors.NewErr("[BlockGetTransactions]parameter count error ")
	}

	blockbytes, err := vm.GetPointerMemory(params[0])
	if err != nil {
		return false, err
	}
	block := &types.Block{}
	err = block.Deserialize(bytes.NewBuffer(blockbytes))
	if err != nil {
		return false, err
	}
	transactionList := make([][]byte, len(block.Transactions))
	for i, tx := range block.Transactions {
		transactionList[i] = tx.ToArray()
	}

	idx, err := vm.SetPointerMemory(transactionList)
	if err != nil {
		return false, err
	}
	vm.RestoreCtx()
	vm.PushResult(uint64(idx))

	return true, nil
}
