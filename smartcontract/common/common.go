

package common

import (
	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/common/log"
	"github.com/mixbee/mixbee/vm/neovm/types"
)

// ConvertReturnTypes return neovm stack element value
// According item types convert to hex string value
// Now neovm support type contain: ByteArray/Integer/Boolean/Array/Struct/Interop/StackItems
func ConvertNeoVmTypeHexString(item interface{}) interface{} {
	if item == nil {
		return nil
	}
	switch v := item.(type) {
	case *types.ByteArray:
		arr, _ := v.GetByteArray()
		return common.ToHexString(arr)
	case *types.Integer:
		i, _ := v.GetBigInteger()
		if i.Sign() == 0 {
			return common.ToHexString([]byte{0})
		} else {
			return common.ToHexString(common.BigIntToNeoBytes(i))
		}
	case *types.Boolean:
		b, _ := v.GetBoolean()
		if b {
			return common.ToHexString([]byte{1})
		} else {
			return common.ToHexString([]byte{0})
		}
	case *types.Array:
		var arr []interface{}
		ar, _ := v.GetArray()
		for _, val := range ar {
			arr = append(arr, ConvertNeoVmTypeHexString(val))
		}
		return arr
	case *types.Struct:
		var arr []interface{}
		ar, _ := v.GetStruct()
		for _, val := range ar {
			arr = append(arr, ConvertNeoVmTypeHexString(val))
		}
		return arr
	case *types.Interop:
		it, _ := v.GetInterface()
		return common.ToHexString(it.ToArray())
	default:
		log.Error("[ConvertTypes] Invalid Types!")
		return nil
	}
}
