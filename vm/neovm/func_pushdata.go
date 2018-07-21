

package neovm

func opPushData(e *ExecutionEngine) (VMState, error) {
	data := getPushData(e)
	PushData(e, data)
	return NONE, nil
}

func getPushData(e *ExecutionEngine) interface{} {
	var data interface{}
	if e.OpCode >= PUSHBYTES1 && e.OpCode <= PUSHBYTES75 {
		data = e.Context.OpReader.ReadBytes(int(e.OpCode))
	}
	switch e.OpCode {
	case PUSH0:
		data = int8(0)
	case PUSHDATA1:
		d, _ := e.Context.OpReader.ReadByte()
		data = e.Context.OpReader.ReadBytes(int(d))
	case PUSHDATA2:
		data = e.Context.OpReader.ReadBytes(int(e.Context.OpReader.ReadUint16()))
	case PUSHDATA4:
		i := int(e.Context.OpReader.ReadInt32())
		data = e.Context.OpReader.ReadBytes(i)
	case PUSHM1, PUSH1, PUSH2, PUSH3, PUSH4, PUSH5, PUSH6, PUSH7, PUSH8, PUSH9, PUSH10, PUSH11, PUSH12, PUSH13, PUSH14, PUSH15, PUSH16:
		data = int8(e.OpCode - PUSH1 + 1)
	}

	return data
}
