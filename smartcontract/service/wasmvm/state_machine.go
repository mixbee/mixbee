

package wasmvm

import (
	"fmt"

	"github.com/mixbee/mixbee/common/log"
	"github.com/mixbee/mixbee/errors"
	"github.com/mixbee/mixbee/vm/wasmvm/exec"
	"github.com/mixbee/mixbee/vm/wasmvm/util"
)

type LogLevel byte

const (
	Debug LogLevel = iota
	Info
	Error
)

type ParamType byte

const (
	Json ParamType = iota
	Raw
)

type WasmStateMachine struct {
	*WasmStateReader
}

func NewWasmStateMachine() *WasmStateMachine {

	stateMachine := WasmStateMachine{WasmStateReader: NewWasmStateReader()}

	//only for debug test
	stateMachine.Register("ContractLogDebug", stateMachine.contractLogDebug)
	stateMachine.Register("ContractLogInfo", stateMachine.contractLogInfo)
	stateMachine.Register("ContractLogError", stateMachine.contractLogError)

	return &stateMachine
}

func (s *WasmStateMachine) contractLogDebug(engine *exec.ExecutionEngine) (bool, error) {
	_, err := contractLog(Debug, engine)
	if err != nil {
		return false, err
	}

	engine.GetVM().RestoreCtx()
	return true, nil
}

func (s *WasmStateMachine) contractLogInfo(engine *exec.ExecutionEngine) (bool, error) {
	_, err := contractLog(Info, engine)
	if err != nil {
		return false, err

	}
	engine.GetVM().RestoreCtx()
	return true, nil
}

func (s *WasmStateMachine) contractLogError(engine *exec.ExecutionEngine) (bool, error) {
	_, err := contractLog(Error, engine)
	if err != nil {
		return false, err
	}

	engine.GetVM().RestoreCtx()
	return true, nil
}

func contractLog(lv LogLevel, engine *exec.ExecutionEngine) (bool, error) {
	vm := engine.GetVM()
	envCall := vm.GetEnvCall()
	params := envCall.GetParams()
	if len(params) != 1 {
		return false, errors.NewErr("parameter count error while call contractLong")
	}

	Idx := params[0]
	addr, err := vm.GetPointerMemory(Idx)
	if err != nil {
		return false, errors.NewErr("get Contract address failed")
	}

	msg := fmt.Sprintf("[WASM Contract] Address:%s message:%s", vm.ContractAddress.ToHexString(), util.TrimBuffToString(addr))

	switch lv {
	case Debug:
		log.Debug(msg)
	case Info:
		log.Info(msg)
	case Error:
		log.Error(msg)
	}
	return true, nil

}
