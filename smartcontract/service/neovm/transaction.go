

package neovm

import (
	"github.com/mixbee/mixbee/core/types"
	vm "github.com/mixbee/mixbee/vm/neovm"
	vmtypes "github.com/mixbee/mixbee/vm/neovm/types"
)

// GetExecutingAddress push transaction's hash to vm stack
// 获得当前交易的 Hash
func TransactionGetHash(service *NeoVmService, engine *vm.ExecutionEngine) error {
	txn, _ := vm.PopInteropInterface(engine)
	tx := txn.(*types.Transaction)
	txHash := tx.Hash()
	vm.PushData(engine, txHash.ToArray())
	return nil
}

// TransactionGetType push transaction's type to vm stack
// 获得当前交易的类型
func TransactionGetType(service *NeoVmService, engine *vm.ExecutionEngine) error {
	txn, _ := vm.PopInteropInterface(engine)
	tx := txn.(*types.Transaction)
	vm.PushData(engine, int(tx.TxType))
	return nil
}

// TransactionGetAttributes push transaction's attributes to vm stack
// 查询当前交易的所有属性
func TransactionGetAttributes(service *NeoVmService, engine *vm.ExecutionEngine) error {
	vm.PopInteropInterface(engine)
	attributList := make([]vmtypes.StackItems, 0)
	vm.PushData(engine, attributList)
	return nil
}
