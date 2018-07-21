
package neovm

import (
	"github.com/mixbee/mixbee/core/types"
	vm "github.com/mixbee/mixbee/vm/neovm"
	vmtypes "github.com/mixbee/mixbee/vm/neovm/types"
)

// BlockGetTransactionCount put block's transactions count to vm stack
func BlockGetTransactionCount(service *NeoVmService, engine *vm.ExecutionEngine) error {
	i, err := vm.PopInteropInterface(engine)
	if err != nil {
		return err
	}
	vm.PushData(engine, len(i.(*types.Block).Transactions))
	return nil
}

// BlockGetTransactions put block's transactions to vm stack
func BlockGetTransactions(service *NeoVmService, engine *vm.ExecutionEngine) error {
	i, err := vm.PopInteropInterface(engine)
	if err != nil {
		return err
	}
	transactions := i.(*types.Block).Transactions
	transactionList := make([]vmtypes.StackItems, 0)
	for _, v := range transactions {
		transactionList = append(transactionList, vmtypes.NewInteropInterface(v))
	}
	vm.PushData(engine, transactionList)
	return nil
}

// BlockGetTransaction put block's transaction to vm stack
func BlockGetTransaction(service *NeoVmService, engine *vm.ExecutionEngine) error {
	i, err := vm.PopInteropInterface(engine)
	if err != nil {
		return err
	}
	index, err := vm.PopInt(engine)
	if err != nil {
		return err
	}
	vm.PushData(engine, i.(*types.Block).Transactions[index])
	return nil
}
