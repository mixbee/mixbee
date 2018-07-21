

package neovm

func opCat(e *ExecutionEngine) (VMState, error) {
	b2, err := PopByteArray(e)
	if err != nil {
		return FAULT, err
	}
	b1, err := PopByteArray(e)
	if err != nil {
		return FAULT, err
	}
	r := Concat(b1, b2)
	PushData(e, r)
	return NONE, nil
}

func opSubStr(e *ExecutionEngine) (VMState, error) {
	count, err := PopInt(e)
	if err != nil {
		return FAULT, err
	}
	index, err := PopInt(e)
	if err != nil {
		return FAULT, err
	}
	arr, err := PopByteArray(e)
	if err != nil {
		return FAULT, err
	}
	b := arr[index : index+count]
	PushData(e, b)
	return NONE, nil
}

func opLeft(e *ExecutionEngine) (VMState, error) {
	count, err := PopInt(e)
	if err != nil {
		return FAULT, err
	}
	s, err := PopByteArray(e)
	if err != nil {
		return FAULT, err
	}
	b := s[:count]
	PushData(e, b)
	return NONE, nil
}

func opRight(e *ExecutionEngine) (VMState, error) {
	count, err := PopInt(e)
	if err != nil {
		return FAULT, err
	}
	arr, err := PopByteArray(e)
	if err != nil {
		return FAULT, err
	}
	b := arr[len(arr)-count:]
	PushData(e, b)
	return NONE, nil
}

func opSize(e *ExecutionEngine) (VMState, error) {
	b, err := PopByteArray(e)
	if err != nil {
		return FAULT, err
	}
	PushData(e, len(b))
	return NONE, nil
}
