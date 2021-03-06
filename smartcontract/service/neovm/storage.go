

package neovm

import (
	"bytes"

	"fmt"
	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/core/states"
	scommon "github.com/mixbee/mixbee/core/store/common"
	"github.com/mixbee/mixbee/errors"
	vm "github.com/mixbee/mixbee/vm/neovm"
)

// StoragePut put smart contract storage item to cache
func StoragePut(service *NeoVmService, engine *vm.ExecutionEngine) error {
	if vm.EvaluationStackCount(engine) < 3 {
		return errors.NewErr("[Context] Too few input parameters ")
	}
	context, err := getContext(engine)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[StoragePut] get pop context error!")
	}
	if context.IsReadOnly {
		return fmt.Errorf("%s", "[StoragePut] storage read only!")
	}
	if err := checkStorageContext(service, context); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[StoragePut] check context error!")
	}

	key, err := vm.PopByteArray(engine)
	if err != nil {
		return err
	}
	if len(key) > 1024 {
		return errors.NewErr("[StoragePut] Storage key to long")
	}

	value, err := vm.PopByteArray(engine)
	if err != nil {
		return err
	}
	service.CloneCache.Add(scommon.ST_STORAGE, getStorageKey(context.Address, key), &states.StorageItem{Value: value})
	return nil
}

// StorageDelete delete smart contract storage item from cache
func StorageDelete(service *NeoVmService, engine *vm.ExecutionEngine) error {
	if vm.EvaluationStackCount(engine) < 2 {
		return errors.NewErr("[Context] Too few input parameters ")
	}
	context, err := getContext(engine)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[StorageDelete] get pop context error!")
	}
	if context.IsReadOnly {
		return fmt.Errorf("%s", "[StorageDelete] storage read only!")
	}
	if err := checkStorageContext(service, context); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[StorageDelete] check context error!")
	}
	ba, err := vm.PopByteArray(engine)
	if err != nil {
		return err
	}
	service.CloneCache.Delete(scommon.ST_STORAGE, getStorageKey(context.Address, ba))

	return nil
}

// StorageGet push smart contract storage item from cache to vm stack
func StorageGet(service *NeoVmService, engine *vm.ExecutionEngine) error {
	if vm.EvaluationStackCount(engine) < 2 {
		return errors.NewErr("[Context] Too few input parameters ")
	}
	context, err := getContext(engine)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[StorageGet] get pop context error!")
	}
	ba, err := vm.PopByteArray(engine)
	if err != nil {
		return err
	}
	item, err := service.CloneCache.Get(scommon.ST_STORAGE, getStorageKey(context.Address, ba))
	if err != nil {
		return err
	}

	if item == nil {
		vm.PushData(engine, []byte{})
	} else {
		vm.PushData(engine, item.(*states.StorageItem).Value)
	}
	return nil
}

// StorageGetContext push smart contract storage context to vm stack
func StorageGetContext(service *NeoVmService, engine *vm.ExecutionEngine) error {
	vm.PushData(engine, NewStorageContext(service.ContextRef.CurrentContext().ContractAddress))
	return nil
}

func StorageGetReadOnlyContext(service *NeoVmService, engine *vm.ExecutionEngine) error {
	context := NewStorageContext(service.ContextRef.CurrentContext().ContractAddress)
	context.IsReadOnly = true
	vm.PushData(engine, context)
	return nil
}

func checkStorageContext(service *NeoVmService, context *StorageContext) error {
	item, err := service.CloneCache.Get(scommon.ST_CONTRACT, context.Address[:])
	if err != nil || item == nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[CheckStorageContext] get context fail!")
	}
	return nil
}

func getContext(engine *vm.ExecutionEngine) (*StorageContext, error) {
	opInterface, err := vm.PopInteropInterface(engine)
	if err != nil {
		return nil, err
	}
	if opInterface == nil {
		return nil, errors.NewErr("[Context] Get storageContext nil")
	}
	context, ok := opInterface.(*StorageContext)
	if !ok {
		return nil, errors.NewErr("[Context] Get storageContext invalid")
	}
	return context, nil
}

func getStorageKey(address common.Address, key []byte) []byte {
	buf := bytes.NewBuffer(nil)
	buf.Write(address[:])
	buf.Write(key)
	return buf.Bytes()
}
