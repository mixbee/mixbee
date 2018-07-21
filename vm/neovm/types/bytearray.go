

package types

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/vm/neovm/interfaces"
)

type ByteArray struct {
	value []byte
}

func NewByteArray(value []byte) *ByteArray {
	var this ByteArray
	this.value = value
	return &this
}

func (this *ByteArray) Equals(other StackItems) bool {
	if this == other {
		return true
	}

	a1 := this.value
	a2, err := other.GetByteArray()
	if err != nil {
		return false
	}

	return bytes.Equal(a1, a2)
}

func (this *ByteArray) GetBigInteger() (*big.Int, error) {
	return common.BigIntFromNeoBytes(this.value), nil
}

func (this *ByteArray) GetBoolean() (bool, error) {
	for _, b := range this.value {
		if b != 0 {
			return true, nil
		}
	}
	return false, nil
}

func (this *ByteArray) GetByteArray() ([]byte, error) {
	return this.value, nil
}

func (this *ByteArray) GetInterface() (interfaces.Interop, error) {
	return nil, fmt.Errorf("%s", "Not support byte array to interface")
}

func (this *ByteArray) GetArray() ([]StackItems, error) {
	return nil, fmt.Errorf("%s", "Not support byte array to array")
}

func (this *ByteArray) GetStruct() ([]StackItems, error) {
	return nil, fmt.Errorf("%s", "Not support byte array to struct")
}

func (this *ByteArray) GetMap() (map[StackItems]StackItems, error) {
	return nil, fmt.Errorf("%s", "Not support byte array to map")
}
