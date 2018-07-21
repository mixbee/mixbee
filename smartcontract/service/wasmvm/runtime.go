
package wasmvm

import (
	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/core/signature"
	"github.com/mixbee/mixbee/errors"
	"github.com/mixbee/mixbee/smartcontract/event"
	"github.com/mixbee/mixbee/vm/wasmvm/exec"
	"github.com/mixbee/mixbee/vm/wasmvm/util"
)

func (this *WasmVmService) runtimeGetTime(engine *exec.ExecutionEngine) (bool, error) {
	vm := engine.GetVM()
	vm.RestoreCtx()
	vm.PushResult(uint64(this.Time))
	return true, nil
}
func (this *WasmVmService) runtimeLog(engine *exec.ExecutionEngine) (bool, error) {

	vm := engine.GetVM()
	envCall := vm.GetEnvCall()
	params := envCall.GetParams()
	if len(params) != 1 {
		return false, errors.NewErr("[RuntimeLog]parameter count error ")
	}

	item, err := vm.GetPointerMemory(params[0])
	if err != nil {
		return false, err
	}

	context := this.ContextRef.CurrentContext()
	txHash := this.Tx.Hash()
	event.PushSmartCodeEvent(txHash, 0, event.EVENT_LOG, &event.LogEventArgs{TxHash: txHash, ContractAddress: context.ContractAddress, Message: string(item)})
	vm.RestoreCtx()

	return true, nil
}

func (this *WasmVmService) runtimeCheckSig(engine *exec.ExecutionEngine) (bool, error) {
	vm := engine.GetVM()
	envCall := vm.GetEnvCall()
	params := envCall.GetParams()
	if len(params) != 3 {
		return false, errors.NewErr("[RuntimeCheckSig]parameter count error ")
	}
	pubKey, err := vm.GetPointerMemory(params[0])
	if err != nil {
		return false, err
	}
	data, err := vm.GetPointerMemory(params[1])
	if err != nil {
		return false, err
	}
	sig, err := vm.GetPointerMemory(params[2])
	if err != nil {
		return false, err
	}
	res := 0
	err = signature.Verify(pubKey, data, sig)
	if err == nil {
		res = 1
	}

	vm.RestoreCtx()
	vm.PushResult(uint64(res))

	return true, nil
}

func (this *WasmVmService) runtimeNotify(engine *exec.ExecutionEngine) (bool, error) {
	vm := engine.GetVM()
	envCall := vm.GetEnvCall()
	params := envCall.GetParams()
	if len(params) != 1 {
		return false, errors.NewErr("[RuntimeNotify]parameter count error ")
	}
	item, err := vm.GetPointerMemory(params[0])
	if err != nil {
		return false, err
	}
	context := this.ContextRef.CurrentContext()

	this.Notifications = append(this.Notifications, &event.NotifyEventInfo{ContractAddress: context.ContractAddress, States: []string{string(item)}})
	vm.RestoreCtx()
	return true, nil
}

func (this *WasmVmService) runtimeCheckWitness(engine *exec.ExecutionEngine) (bool, error) {
	vm := engine.GetVM()

	envCall := vm.GetEnvCall()
	params := envCall.GetParams()
	if len(params) != 1 {
		return false, errors.NewErr("[CheckWitness]get parameter count error!")
	}
	data, err := vm.GetPointerMemory(params[0])
	if err != nil {
		return false, errors.NewErr("[CheckWitness]" + err.Error())
	}
	address, err := common.AddressFromBase58(util.TrimBuffToString(data))
	if err != nil {
		return false, errors.NewErr("[CheckWitness]" + err.Error())
	}
	chkRes := this.ContextRef.CheckWitness(address)
	res := 0
	if chkRes == true {
		res = 1
	}
	vm.RestoreCtx()
	if vm.GetEnvCall().GetReturns() {
		vm.PushResult(uint64(res))
	}
	return true, nil
}
