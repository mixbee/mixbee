

package neovm

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/mixbee/mixbee/vm/neovm/utils"
)

func TestGetPushData(t *testing.T) {
	var e ExecutionEngine
	var iRet int8
	var ret []byte
	e.Context = NewExecutionContext(&e, nil)

	e.OpCode = PUSH0
	iRet, ok := getPushData(&e).(int8)
	if !ok || iRet != 0 {
		t.Error("NeoVM getPushData PUSH0 execute failed.")
	}

	e.OpCode = PUSHDATA1
	e.Context.OpReader = utils.NewVmReader([]byte{4, 1, 1, 1, 1})
	ret, ok = getPushData(&e).([]byte)
	if !ok || !bytes.Equal(ret, []byte{1, 1, 1, 1}) {
		t.Fatal("NeoVM getPushData PUSHDATA1 execute failed.")
	}

	e.OpCode = PUSHDATA2
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, 4)
	b = append(b, []byte{1, 1, 1, 1}...)
	e.Context.OpReader = utils.NewVmReader(b)
	ret, ok = getPushData(&e).([]byte)
	if !ok || !bytes.Equal(ret, []byte{1, 1, 1, 1}) {
		t.Fatal("NeoVM getPushData PUSHDATA2 execute failed.")
	}

	e.OpCode = PUSHDATA4
	b = make([]byte, 4)
	binary.LittleEndian.PutUint32(b, 4)
	b = append(b, []byte{1, 1, 1, 1}...)
	e.Context.OpReader = utils.NewVmReader(b)
	ret, ok = getPushData(&e).([]byte)
	if !ok || !bytes.Equal(ret, []byte{1, 1, 1, 1}) {
		t.Fatal("NeoVM getPushData PUSHDATA4 execute failed.")
	}

	for _, opCode := range []OpCode{PUSHM1, PUSH1, PUSH2, PUSH3, PUSH4, PUSH5, PUSH6, PUSH7,
		PUSH8, PUSH9, PUSH10, PUSH11, PUSH12, PUSH13, PUSH14, PUSH15, PUSH16} {
		e.OpCode = opCode
		expect := int8(opCode - PUSH1 + 1)
		iRet, ok = getPushData(&e).(int8)
		if !ok || expect != iRet {
			t.Fatal("NeoVM getPushData execute failed.")
		}
	}
}
