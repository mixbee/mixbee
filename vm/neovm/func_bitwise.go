

package neovm

func opInvert(e *ExecutionEngine) (VMState, error) {
	i, err := PopBigInt(e)
	if err != nil {
		return FAULT, err
	}
	PushData(e, i.Not(i))
	return NONE, nil
}

func opEqual(e *ExecutionEngine) (VMState, error) {
	b1 := PopStackItem(e)
	b2 := PopStackItem(e)
	PushData(e, b1.Equals(b2))
	return NONE, nil
}
