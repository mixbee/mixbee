

package types

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/vm/neovm/interfaces"
)

type Integer struct {
	value *big.Int
}

func NewInteger(value *big.Int) *Integer {
	var this Integer
	this.value = value
	return &this
}

func (this *Integer) Equals(other StackItems) bool {
	if this == other {
		return true
	}
	if other == nil {
		return false
	}

	v, err := other.GetBigInteger()
	if err == nil {
		if this.value.Cmp(v) == 0 {
			return true
		} else {
			return false
		}
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

func (this *Integer) GetBigInteger() (*big.Int, error) {
	return this.value, nil
}

func (this *Integer) GetBoolean() (bool, error) {
	if this.value.Cmp(big.NewInt(0)) == 0 {
		return false, nil
	}
	return true, nil
}

func (this *Integer) GetByteArray() ([]byte, error) {
	return common.BigIntToNeoBytes(this.value), nil
}

func (this *Integer) GetInterface() (interfaces.Interop, error) {
	return nil, fmt.Errorf("%s", "Not support integer to interface")
}

func (this *Integer) GetArray() ([]StackItems, error) {
	return nil, fmt.Errorf("%s", "Not support integer to array")
}

func (this *Integer) GetStruct() ([]StackItems, error) {
	return nil, fmt.Errorf("%s", "Not support integer to struct")
}

func (this *Integer) GetMap() (map[StackItems]StackItems, error) {
	return nil, fmt.Errorf("%s", "Not support integer to map")
}
