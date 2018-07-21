

package neovm

import (
	"github.com/mixbee/mixbee/errors"
	vm "github.com/mixbee/mixbee/vm/neovm"
)

func StoreGasCost(engine *vm.ExecutionEngine) (uint64, error) {
	key, err := vm.PeekNByteArray(1, engine)
	if err != nil {
		return 0, err
	}
	value, err := vm.PeekNByteArray(2, engine)
	if err != nil {
		return 0, err
	}
	if putCost, ok := GAS_TABLE.Load(STORAGE_PUT_NAME); ok {
		return uint64(((len(key)+len(value)-1)/1024 + 1)) * putCost.(uint64), nil
	} else {
		return uint64(0), errors.NewErr("[StoreGasCost] get STORAGE_PUT_NAME gas failed")
	}
}

func GasPrice(engine *vm.ExecutionEngine, name string) (uint64, error) {
	switch name {
	case STORAGE_PUT_NAME:
		return StoreGasCost(engine)
	default:
		if value, ok := GAS_TABLE.Load(name); ok {
			return value.(uint64), nil
		}
		return OPCODE_GAS, nil
	}
}
