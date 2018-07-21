

package neovm

import (
	"testing"

	"math/big"

	"github.com/mixbee/mixbee/vm/neovm/types"
)

func TestOpBigInt(t *testing.T) {
	var e ExecutionEngine
	e.EvaluationStack = NewRandAccessStack()

	for _, code := range []OpCode{INC, DEC, NEGATE, ABS, PUSH0} {
		e.EvaluationStack.Push(NewStackItem(types.NewInteger(big.NewInt(-10))))
		e.OpCode = code
		opBigInt(&e)
		v, err := PopBigInt(&e)
		if err != nil {
			t.Fatal("NeoVM OpBigInt test failed.")
		}
		if code == INC && !(v.Cmp(big.NewInt(-9)) == 0) {
			t.Fatal("NeoVM OpBigInt test failed.")
		} else if code == DEC && !(v.Cmp(big.NewInt(-11)) == 0) {
			t.Fatal("NeoVM OpBigInt test failed.")
		} else if code == NEGATE && !(v.Cmp(big.NewInt(10)) == 0) {
			t.Fatal("NeoVM OpBigInt test failed.")
		} else if code == ABS && !(v.Cmp(big.NewInt(10)) == 0) {
			t.Fatal("NeoVM OpBigInt test failed.")
		} else if code == PUSH0 && !(v.Cmp(big.NewInt(-10)) == 0) {
			t.Fatal("NeoVM OpBigInt test failed.")
		}
	}
}

func TestOpSign(t *testing.T) {
	var e ExecutionEngine
	e.EvaluationStack = NewRandAccessStack()
	i := big.NewInt(10)
	e.EvaluationStack.Push(NewStackItem(types.NewInteger(i)))

	opSign(&e)
	v, err := PopInt(&e)
	if err != nil {
		t.Fatal("NeoVM OpSign test failed.")
	}
	if !(v == i.Sign()) {
		t.Fatal("NeoVM OpSign test failed.")
	}
}

func TestOpNot(t *testing.T) {
	var e ExecutionEngine
	e.EvaluationStack = NewRandAccessStack()
	e.EvaluationStack.Push(NewStackItem(types.NewBoolean(true)))

	opNot(&e)
	v, err := PopBoolean(&e)
	if err != nil {
		t.Fatal("NeoVM OpNot test failed.")
	}
	if !(v == false) {
		t.Fatal("NeoVM OpNot test failed.")
	}
}

func TestOpNz(t *testing.T) {
	var e ExecutionEngine
	e.EvaluationStack = NewRandAccessStack()

	e.EvaluationStack.Push(NewStackItem(types.NewInteger(big.NewInt(0))))
	e.OpCode = NZ
	opNz(&e)
	v, err := PopBoolean(&e)
	if err != nil {
		t.Fatal("NeoVM OpNz test failed.")
	}
	if v == true {
		t.Fatal("NeoVM OpNz test failed.")
	}
	e.EvaluationStack.Push(NewStackItem(types.NewInteger(big.NewInt(10))))
	opNz(&e)

	v, err = PopBoolean(&e)
	if err != nil {
		t.Fatal("NeoVM OpNz test failed.")
	}
	if v == false {
		t.Fatal("NeoVM OpNz test failed.")
	}
	e.EvaluationStack.Push(NewStackItem(types.NewInteger(big.NewInt(0))))
	e.OpCode = PUSH0
	opNz(&e)

	v, err = PopBoolean(&e)
	if err != nil {
		t.Fatal("NeoVM OpNz test failed.")
	}
	if v == true {
		t.Fatal("NeoVM OpNz test failed.")
	}
}
