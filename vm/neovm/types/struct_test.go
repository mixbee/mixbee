
package types

import (
	"math/big"
	"testing"
)

func TestStruct_Clone(t *testing.T) {
	s := NewStruct(nil)
	//d := NewStruct([]StackItems{s})
	k := NewStruct([]StackItems{NewInteger(big.NewInt(1))})
	for i := 0; i < MAX_STRCUT_DEPTH-2; i++ {
		k = NewStruct([]StackItems{k})
	}
	//k.Add(d)
	s.Add(k)

	if _, err := s.Clone(); err != nil {
		t.Fatal(err)
	}

}
