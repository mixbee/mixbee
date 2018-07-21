

package neovm

func opToDupFromAltStack(e *ExecutionEngine) (VMState, error) {
	Push(e, e.AltStack.Peek(0))
	return NONE, nil
}

func opToAltStack(e *ExecutionEngine) (VMState, error) {
	e.AltStack.Push(PopStackItem(e))
	return NONE, nil
}

func opFromAltStack(e *ExecutionEngine) (VMState, error) {
	Push(e, e.AltStack.Pop())
	return NONE, nil
}

func opXDrop(e *ExecutionEngine) (VMState, error) {
	n, err := PopInt(e)
	if err != nil {
		return FAULT, err
	}
	e.EvaluationStack.Remove(n)
	return NONE, nil
}

func opXSwap(e *ExecutionEngine) (VMState, error) {
	n, err := PopInt(e)
	if err != nil {
		return FAULT, err
	}
	if n == 0 {
		return NONE, nil
	}
	e.EvaluationStack.Swap(0, n)
	return NONE, nil
}

func opXTuck(e *ExecutionEngine) (VMState, error) {
	n, err := PopInt(e)
	if err != nil {
		return FAULT, err
	}
	e.EvaluationStack.Insert(n, PeekStackItem(e))
	return NONE, nil
}

func opDepth(e *ExecutionEngine) (VMState, error) {
	PushData(e, Count(e))
	return NONE, nil
}

func opDrop(e *ExecutionEngine) (VMState, error) {
	PopStackItem(e)
	return NONE, nil
}

func opDup(e *ExecutionEngine) (VMState, error) {
	Push(e, PeekStackItem(e))
	return NONE, nil
}

func opNip(e *ExecutionEngine) (VMState, error) {
	x2 := PopStackItem(e)
	PopStackItem(e)
	Push(e, x2)
	return NONE, nil
}

func opOver(e *ExecutionEngine) (VMState, error) {
	x2 := PopStackItem(e)
	x1 := PeekStackItem(e)

	Push(e, x2)
	Push(e, x1)
	return NONE, nil
}

func opPick(e *ExecutionEngine) (VMState, error) {
	n, err := PopInt(e)
	if err != nil {
		return FAULT, err
	}
	Push(e, e.EvaluationStack.Peek(n))
	return NONE, nil
}

func opRoll(e *ExecutionEngine) (VMState, error) {
	n, err := PopInt(e)
	if err != nil {
		return FAULT, err
	}
	if n == 0 {
		return NONE, nil
	}
	Push(e, e.EvaluationStack.Remove(n))
	return NONE, nil
}

func opRot(e *ExecutionEngine) (VMState, error) {
	x3 := PopStackItem(e)
	x2 := PopStackItem(e)
	x1 := PopStackItem(e)
	Push(e, x2)
	Push(e, x3)
	Push(e, x1)
	return NONE, nil
}

func opSwap(e *ExecutionEngine) (VMState, error) {
	x2 := PopStackItem(e)
	x1 := PopStackItem(e)
	Push(e, x2)
	Push(e, x1)
	return NONE, nil
}

func opTuck(e *ExecutionEngine) (VMState, error) {
	x2 := PopStackItem(e)
	x1 := PopStackItem(e)
	Push(e, x2)
	Push(e, x1)
	Push(e, x2)
	return NONE, nil
}
