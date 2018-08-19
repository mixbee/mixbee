

package neovm

import (
	"fmt"

	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/core/payload"
	scommon "github.com/mixbee/mixbee/core/store/common"
	"github.com/mixbee/mixbee/core/types"
	"github.com/mixbee/mixbee/errors"
	vm "github.com/mixbee/mixbee/vm/neovm"
)

// ContractCreate create a new smart contract on blockchain, and put it to vm stack
func ContractCreate(service *NeoVmService, engine *vm.ExecutionEngine) error {
	contract, err := isContractParamValid(engine)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[ContractCreate] contract parameters invalid!")
	}
	contractAddress := types.AddressFromVmCode(contract.Code)
	state, err := service.CloneCache.GetOrAdd(scommon.ST_CONTRACT, contractAddress[:], contract)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[ContractCreate] GetOrAdd error!")
	}
	vm.PushData(engine, state)
	return nil
}

// ContractMigrate migrate old smart contract to a new contract, and destroy old contract
func ContractMigrate(service *NeoVmService, engine *vm.ExecutionEngine) error {
	// get new contract
	contract, err := isContractParamValid(engine)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[ContractMigrate] contract parameters invalid!")
	}
	// calculate new contract address
	contractAddress := types.AddressFromVmCode(contract.Code)

	// Find out if the contract already exists based on the contract address, if not exists , then continue
	if err := isContractExist(service, contractAddress); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[ContractMigrate] contract invalid!")
	}
	// smart contract execute context struct , contains address and code
	context := service.ContextRef.CurrentContext()

	// new contract add to clonecache
	service.CloneCache.Add(scommon.ST_CONTRACT, contractAddress[:], contract)
	// Replace the old key with the new one, and return old key-value list.
	items, err := storeMigration(service, context.ContractAddress, contractAddress)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[ContractMigrate] contract store migration error!")
	}
	// delete from cache
	service.CloneCache.Delete(scommon.ST_CONTRACT, context.ContractAddress[:])
	for _, v := range items {
		service.CloneCache.Delete(scommon.ST_STORAGE, []byte(v.Key))
	}
	vm.PushData(engine, contract)
	return nil
}

// ContractDestory destroy a contract
func ContractDestory(service *NeoVmService, engine *vm.ExecutionEngine) error {
	context := service.ContextRef.CurrentContext()
	if context == nil {
		return errors.NewErr("[ContractDestory] current contract context invalid!")
	}
	item, err := service.CloneCache.Store.TryGet(scommon.ST_CONTRACT, context.ContractAddress[:])

	if err != nil || item == nil {
		return errors.NewErr("[ContractDestory] get current contract fail!")
	}

	service.CloneCache.Delete(scommon.ST_CONTRACT, context.ContractAddress[:])
	stateValues, err := service.CloneCache.Store.Find(scommon.ST_CONTRACT, context.ContractAddress[:])
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[ContractDestory] find error!")
	}
	for _, v := range stateValues {
		service.CloneCache.Delete(scommon.ST_STORAGE, []byte(v.Key))
	}
	return nil
}

// ContractGetStorageContext put contract storage context to vm stack
// 获得合约的存储上下文
func ContractGetStorageContext(service *NeoVmService, engine *vm.ExecutionEngine) error {
	if vm.EvaluationStackCount(engine) < 1 {
		return errors.NewErr("[GetStorageContext] Too few input parameter!")
	}
	opInterface, err := vm.PopInteropInterface(engine)
	if err != nil {
		return err
	}
	if opInterface == nil {
		return errors.NewErr("[GetStorageContext] Pop data nil!")
	}
	contractState, ok := opInterface.(*payload.DeployCode)
	if !ok {
		return errors.NewErr("[GetStorageContext] Pop data not contract!")
	}
	address := types.AddressFromVmCode(contractState.Code)
	item, err := service.CloneCache.Store.TryGet(scommon.ST_CONTRACT, address[:])
	if err != nil || item == nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[GetStorageContext] Get StorageContext nil")
	}
	if address != service.ContextRef.CurrentContext().ContractAddress {
		return errors.NewErr("[GetStorageContext] CodeHash not equal!")
	}
	vm.PushData(engine, NewStorageContext(address))
	return nil
}

// ContractGetCode put contract to vm stack
func ContractGetCode(service *NeoVmService, engine *vm.ExecutionEngine) error {
	i, err := vm.PopInteropInterface(engine)
	if err != nil {
		return err
	}
	vm.PushData(engine, i.(*payload.DeployCode).Code)
	return nil
}

func isContractParamValid(engine *vm.ExecutionEngine) (*payload.DeployCode, error) {
	if vm.EvaluationStackCount(engine) < 7 {
		return nil, errors.NewErr("[Contract] Too few input parameters")
	}
	code, err := vm.PopByteArray(engine)
	if err != nil {
		return nil, err
	}
	if len(code) > 1024*1024 {
		return nil, errors.NewErr("[Contract] Code too long!")
	}
	needStorage, err := vm.PopBoolean(engine)
	if err != nil {
		return nil, err
	}
	name, err := vm.PopByteArray(engine)
	if err != nil {
		return nil, err
	}
	if len(name) > 252 {
		return nil, errors.NewErr("[Contract] Name too long!")
	}
	version, err := vm.PopByteArray(engine)
	if err != nil {
		return nil, err
	}
	if len(version) > 252 {
		return nil, errors.NewErr("[Contract] Version too long!")
	}
	author, err := vm.PopByteArray(engine)
	if err != nil {
		return nil, err
	}
	if len(author) > 252 {
		return nil, errors.NewErr("[Contract] Author too long!")
	}
	email, err := vm.PopByteArray(engine)
	if err != nil {
		return nil, err
	}
	if len(email) > 252 {
		return nil, errors.NewErr("[Contract] Email too long!")
	}
	desc, err := vm.PopByteArray(engine)
	if err != nil {
		return nil, err
	}
	if len(desc) > 65536 {
		return nil, errors.NewErr("[Contract] Desc too long!")
	}
	contract := &payload.DeployCode{
		Code:        code,
		NeedStorage: needStorage,
		Name:        string(name),
		Version:     string(version),
		Author:      string(author),
		Email:       string(email),
		Description: string(desc),
	}
	return contract, nil
}

func isContractExist(service *NeoVmService, contractAddress common.Address) error {
	item, err := service.CloneCache.Get(scommon.ST_CONTRACT, contractAddress[:])

	if err != nil || item != nil {
		return fmt.Errorf("[Contract] Get contract %x error or contract exist!", contractAddress)
	}
	return nil
}

func storeMigration(service *NeoVmService, oldAddr common.Address, newAddr common.Address) ([]*scommon.StateItem, error) {
	stateValues, err := service.CloneCache.Store.Find(scommon.ST_STORAGE, oldAddr[:])
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "[Contract] Find error!")
	}
	for _, v := range stateValues {
		service.CloneCache.Add(scommon.ST_STORAGE, getStorageKey(newAddr, []byte(v.Key)[20:]), v.Value)
	}
	return stateValues, nil
}
