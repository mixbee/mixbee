

package neovm

import (
	"bytes"
	"math/big"
	"testing"

	vtypes "github.com/mixbee/mixbee/vm/neovm/types"
)

func TestOpCat(t *testing.T) {
	var e ExecutionEngine
	stack := NewRandAccessStack()
	stack.Push(NewStackItem(vtypes.NewByteArray([]byte("aaa"))))
	stack.Push(NewStackItem(vtypes.NewByteArray([]byte("bbb"))))
	e.EvaluationStack = stack

	opCat(&e)
	v, err := PeekNByteArray(0, &e)
	if err != nil {
		t.Fatal("NeoVM OpCat test failed.")
	}
	if Count(&e) != 1 || !bytes.Equal(v, []byte("aaabbb")) {
		t.Fatalf("NeoVM OpCat test failed, expect aaabbb, got %s.", string(v))
	}
}

func TestOpSubStr(t *testing.T) {
	var e ExecutionEngine
	stack := NewRandAccessStack()
	stack.Push(NewStackItem(vtypes.NewByteArray([]byte("12345"))))
	stack.Push(NewStackItem(vtypes.NewInteger(big.NewInt(1))))
	stack.Push(NewStackItem(vtypes.NewInteger(big.NewInt(4))))
	e.EvaluationStack = stack

	opSubStr(&e)
	v, err := PeekNByteArray(0, &e)
	if err != nil {
		t.Fatal("NeoVM OpSubStr test failed.")
	}
	if !bytes.Equal(v, []byte("2345")) {
		t.Fatalf("NeoVM OpSubStr test failed, expect 234, got %s.", string(v))
	}
}

func TestOpLeft(t *testing.T) {
	var e ExecutionEngine
	stack := NewRandAccessStack()
	stack.Push(NewStackItem(vtypes.NewByteArray([]byte("12345"))))
	stack.Push(NewStackItem(vtypes.NewInteger(big.NewInt(4))))
	e.EvaluationStack = stack

	opLeft(&e)
	v, err := PeekNByteArray(0, &e)
	if err != nil {
		t.Fatal("NeoVM OpLeft test failed.")
	}
	if !bytes.Equal(v, []byte("1234")) {
		t.Fatalf("NeoVM OpLeft test failed, expect 1234, got %s.", string(v))
	}
}

func TestOpRight(t *testing.T) {
	var e ExecutionEngine
	stack := NewRandAccessStack()
	stack.Push(NewStackItem(vtypes.NewByteArray([]byte("12345"))))
	stack.Push(NewStackItem(vtypes.NewInteger(big.NewInt(3))))
	e.EvaluationStack = stack

	opRight(&e)
	v, err := PeekNByteArray(0, &e)
	if err != nil {
		t.Fatal("NeoVM OpRight test failed.")
	}
	if !bytes.Equal(v, []byte("345")) {
		t.Fatalf("NeoVM OpRight test failed, expect 345, got %s.", string(v))
	}
}

func TestOpSize(t *testing.T) {
	var e ExecutionEngine
	stack := NewRandAccessStack()
	stack.Push(NewStackItem(vtypes.NewByteArray([]byte("12345"))))
	e.EvaluationStack = stack

	opSize(&e)
	v, err := PeekInt(&e)
	if err != nil {
		t.Fatal("NeoVM OpSize test failed.")
	}
	if v != 5 {
		t.Fatalf("NeoVM OpSize test failed, expect 5, got %d.", v)
	}
}
