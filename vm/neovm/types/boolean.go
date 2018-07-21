

package types

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/mixbee/mixbee/vm/neovm/interfaces"
)

type Boolean struct {
	value bool
}

func NewBoolean(value bool) *Boolean {
	var this Boolean
	this.value = value
	return &this
}

func (this *Boolean) Equals(other StackItems) bool {
	if this == other {
		return true
	}
	b, err := other.GetByteArray()
	if err != nil {
		return false
	}

	tb, err := this.GetByteArray()
	if err != nil {
		return false
	}

	return bytes.Equal(tb, b)
}

func (this *Boolean) GetBigInteger() (*big.Int, error) {
	if this.value {
		return big.NewInt(1), nil
	}
	return big.NewInt(0), nil
}

func (this *Boolean) GetBoolean() (bool, error) {
	return this.value, nil
}

func (this *Boolean) GetByteArray() ([]byte, error) {
	if this.value {
		return []byte{1}, nil
	}
	return []byte{0}, nil
}

func (this *Boolean) GetInterface() (interfaces.Interop, error) {
	return nil, fmt.Errorf("%s", "Not support boolean to interface")
}

func (this *Boolean) GetArray() ([]StackItems, error) {
	return nil, fmt.Errorf("%s", "Not support boolean to array")
}

func (this *Boolean) GetStruct() ([]StackItems, error) {
	return nil, fmt.Errorf("%s", "Not support boolean to struct")
}

func (this *Boolean) GetMap() (map[StackItems]StackItems, error) {
	return nil, fmt.Errorf("%s", "Not support boolean to map")
}
