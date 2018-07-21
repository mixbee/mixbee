

package neovm

func opBigInt(e *ExecutionEngine) (VMState, error) {
	x, err := PopBigInt(e)
	if err != nil {
		return FAULT, err
	}
	PushData(e, BigIntOp(x, e.OpCode))
	return NONE, nil
}

func opSign(e *ExecutionEngine) (VMState, error) {
	x, err := PopBigInt(e)
	if err != nil {
		return FAULT, err
	}
	PushData(e, x.Sign())
	return NONE, nil
}

func opNot(e *ExecutionEngine) (VMState, error) {
	x, err := PopBoolean(e)
	if err != nil {
		return FAULT, err
	}
	PushData(e, !x)
	return NONE, nil
}

func opNz(e *ExecutionEngine) (VMState, error) {
	x, err := PopBigInt(e)
	if err != nil {
		return FAULT, err
	}
	PushData(e, BigIntComp(x, e.OpCode))
	return NONE, nil
}

func opBigIntZip(e *ExecutionEngine) (VMState, error) {
	x2, err := PopBigInt(e)
	if err != nil {
		return FAULT, err
	}
	x1, err := PopBigInt(e)
	if err != nil {
		return FAULT, err
	}
	b := BigIntZip(x1, x2, e.OpCode)
	PushData(e, b)
	return NONE, nil
}

func opBoolZip(e *ExecutionEngine) (VMState, error) {
	x2, err := PopBoolean(e)
	if err != nil {
		return FAULT, err
	}
	x1, err := PopBoolean(e)
	if err != nil {
		return FAULT, err
	}
	PushData(e, BoolZip(x1, x2, e.OpCode))
	return NONE, nil
}

func opBigIntComp(e *ExecutionEngine) (VMState, error) {
	x2, err := PopBigInt(e)
	if err != nil {
		return FAULT, err
	}
	x1, err := PopBigInt(e)
	if err != nil {
		return FAULT, err
	}
	PushData(e, BigIntMultiComp(x1, x2, e.OpCode))
	return NONE, nil
}

func opWithIn(e *ExecutionEngine) (VMState, error) {
	b, err := PopBigInt(e)
	if err != nil {
		return FAULT, err
	}
	a, err := PopBigInt(e)
	if err != nil {
		return FAULT, err
	}
	c, err := PopBigInt(e)
	if err != nil {
		return FAULT, err
	}
	PushData(e, WithInOp(c, a, b))
	return NONE, nil
}
