

package neovm

import (
	"github.com/mixbee/mixbee/vm/neovm/errors"
)

func opNop(e *ExecutionEngine) (VMState, error) {
	return NONE, nil
}

func opJmp(e *ExecutionEngine) (VMState, error) {
	offset := int(e.Context.OpReader.ReadInt16())

	offset = e.Context.GetInstructionPointer() + offset - 3

	if offset < 0 || offset > len(e.Context.Code) {
		return FAULT, errors.ERR_FAULT
	}
	var fValue = true

	if e.OpCode > JMP {
		if EvaluationStackCount(e) < 1 {
			return FAULT, errors.ERR_UNDER_STACK_LEN
		}
		var err error
		fValue, err = PopBoolean(e)
		if err != nil {
			return FAULT, err
		}
		if e.OpCode == JMPIFNOT {
			fValue = !fValue
		}
	}

	if fValue {
		e.Context.SetInstructionPointer(int64(offset))
	}
	return NONE, nil
}

func opCall(e *ExecutionEngine) (VMState, error) {
	context := e.Context.Clone()
	e.Context.SetInstructionPointer(int64(e.Context.GetInstructionPointer() + 2))
	e.OpCode = JMP
	e.PushContext(context)
	return opJmp(e)
}

func opRet(e *ExecutionEngine) (VMState, error) {
	e.PopContext()
	return NONE, nil
}
