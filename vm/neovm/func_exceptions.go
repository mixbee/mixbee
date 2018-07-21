

package neovm

func opThrow(e *ExecutionEngine) (VMState, error) {
	return FAULT, nil
}

func opThrowIfNot(e *ExecutionEngine) (VMState, error) {
	b, err := PopBoolean(e)
	if err != nil {
		return FAULT, err
	}
	if !b {
		return FAULT, nil
	}
	return NONE, nil
}
