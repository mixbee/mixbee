

package neovm

import (
	"github.com/mixbee/mixbee/core/types"
	vm "github.com/mixbee/mixbee/vm/neovm"
)

// AttributeGetUsage put attribute's usage to vm stack
func AttributeGetUsage(service *NeoVmService, engine *vm.ExecutionEngine) error {
	i, err := vm.PopInteropInterface(engine)
	if err != nil {
		return err
	}
	vm.PushData(engine, int(i.(*types.TxAttribute).Usage))
	return nil
}

// AttributeGetData put attribute's data to vm stack
func AttributeGetData(service *NeoVmService, engine *vm.ExecutionEngine) error {
	i, err := vm.PopInteropInterface(engine)
	if err != nil {
		return err
	}
	vm.PushData(engine, i.(*types.TxAttribute).Data)
	return nil
}
