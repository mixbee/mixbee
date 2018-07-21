

package neovm

func opHash(e *ExecutionEngine) (VMState, error) {
	x, err := PopByteArray(e)
	if err != nil {
		return FAULT, err
	}
	PushData(e, Hash(x, e))
	return NONE, nil
}
