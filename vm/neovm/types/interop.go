

package types

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/mixbee/mixbee/vm/neovm/interfaces"
)

type Interop struct {
	_object interfaces.Interop
}

func NewInteropInterface(value interfaces.Interop) *Interop {
	var ii Interop
	ii._object = value
	return &ii
}

func (this *Interop) Equals(other StackItems) bool {
	v, err := other.GetInterface()
	if err != nil {
		return false
	}
	if this._object == nil || v == nil {
		return false
	}
	if !bytes.Equal(this._object.ToArray(), v.ToArray()) {
		return false
	}
	return true
}

func (this *Interop) GetBigInteger() (*big.Int, error) {
	return nil, fmt.Errorf("%s", "Not support interface to biginteger")
}

func (this *Interop) GetBoolean() (bool, error) {
	if this._object == nil {
		return false, nil
	}
	return true, nil
}

func (this *Interop) GetByteArray() ([]byte, error) {
	return nil, fmt.Errorf("%s", "Not support interface to bytearray")
}

func (this *Interop) GetInterface() (interfaces.Interop, error) {
	return this._object, nil
}

func (this *Interop) GetArray() ([]StackItems, error) {
	return nil, fmt.Errorf("%s", "Not support interface to array")
}

func (this *Interop) GetStruct() ([]StackItems, error) {
	return nil, fmt.Errorf("%s", "Not support interface to struct")
}

func (this *Interop) GetMap() (map[StackItems]StackItems, error) {
	return nil, fmt.Errorf("%s", "Not support interface to map")
}
