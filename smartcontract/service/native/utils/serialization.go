

package utils

import (
	"fmt"
	"io"
	"math/big"
	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/common/serialization"
	"github.com/mixbee/mixbee/vm/neovm/types"
)

func WriteVarUint(w io.Writer, value uint64) error {
	if err := serialization.WriteVarBytes(w, types.BigIntToBytes(big.NewInt(int64(value)))); err != nil {
		return fmt.Errorf("serialize value error:%v", err)
	}
	return nil
}

func ReadVarUint(r io.Reader) (uint64, error) {
	value, err := serialization.ReadVarBytes(r)
	if err != nil {
		return 0, fmt.Errorf("deserialize value error:%v", err)
	}
	v := types.BigIntFromBytes(value)
	if v.Cmp(big.NewInt(0)) < 0 {
		return 0, fmt.Errorf("%s", "value should not be a negative number.")
	}
	return v.Uint64(), nil
}

func ReadVarUint32(r io.Reader) (uint32, error) {
	value,err := ReadVarUint(r)
	if err != nil {
		return 0,err
	}
	return uint32(value),nil
}

func WriteAddress(w io.Writer, address common.Address) error {
	if err := serialization.WriteVarBytes(w, address[:]); err != nil {
		return fmt.Errorf("serialize value error:%v", err)
	}
	return nil
}

func WriteString(w io.Writer, str string) error {
	if err := serialization.WriteVarBytes(w, []byte(str)); err != nil {
		return fmt.Errorf("serialize value error:%v", err)
	}
	return nil
}

func ReadString(r io.Reader) (string, error) {
	from, err := serialization.ReadVarBytes(r)
	if err != nil {
		return "", fmt.Errorf("[State] deserialize from error:%v", err)
	}
	return string(from),nil
}

func ReadAddress(r io.Reader) (common.Address, error) {
	from, err := serialization.ReadVarBytes(r)
	if err != nil {
		return common.Address{}, fmt.Errorf("[State] deserialize from error:%v", err)
	}
	return common.AddressParseFromBytes(from)
}
