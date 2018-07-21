

package types

import (
	"math/big"

	"github.com/mixbee/mixbee/vm/neovm/interfaces"
)

type StackItems interface {
	Equals(other StackItems) bool
	GetBigInteger() (*big.Int, error)
	GetBoolean() (bool, error)
	GetByteArray() ([]byte, error)
	GetInterface() (interfaces.Interop, error)
	GetArray() ([]StackItems, error)
	GetStruct() ([]StackItems, error)
	GetMap() (map[StackItems]StackItems, error)
}

const (
	ByteArrayType byte = 0x00
	BooleanType   byte = 0x01
	IntegerType   byte = 0x02
	InterfaceType byte = 0x40
	ArrayType     byte = 0x80
	StructType    byte = 0x81
	MapType       byte = 0x82
)
