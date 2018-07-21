

package neovm

import (
	"math/big"
	"testing"

	vtypes "github.com/mixbee/mixbee/vm/neovm/types"
)

func TestOpInvert(t *testing.T) {
	var e ExecutionEngine
	stack := NewRandAccessStack()
	stack.Push(NewStackItem(vtypes.NewInteger(big.NewInt(123456789))))
	e.EvaluationStack = stack

	opInvert(&e)
	i := big.NewInt(123456789)

	v, err := PeekBigInteger(&e)
	if err != nil {
		t.Fatal("NeoVM OpInvert test failed.")
	}
	if v.Cmp(i.Not(i)) != 0 {
		t.Fatal("NeoVM OpInvert test failed.")
	}
}

func TestOpEqual(t *testing.T) {
	var e ExecutionEngine
	stack := NewRandAccessStack()
	stack.Push(NewStackItem(vtypes.NewInteger(big.NewInt(123456789))))
	stack.Push(NewStackItem(vtypes.NewInteger(big.NewInt(123456789))))
	e.EvaluationStack = stack

	opEqual(&e)
	v, err := PopBoolean(&e)
	if err != nil {
		t.Fatal("NeoVM OpEqual test failed.")
	}
	if !v {
		t.Fatal("NeoVM OpEqual test failed.")
	}
}
