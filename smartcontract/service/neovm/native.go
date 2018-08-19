

package neovm

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/common/serialization"
	"github.com/mixbee/mixbee/smartcontract/service/native"
	"github.com/mixbee/mixbee/smartcontract/states"
	vm "github.com/mixbee/mixbee/vm/neovm"
	"github.com/mixbee/mixbee/vm/neovm/types"
	"math/big"
)

func NativeInvoke(service *NeoVmService, engine *vm.ExecutionEngine) error {
	count := vm.EvaluationStackCount(engine)
	if count < 4 {
		return fmt.Errorf("invoke native contract invalid parameters %d < 4 ", count)
	}
	// parm1: version
	version, err := vm.PopInt(engine)
	if err != nil {
		return err
	}
	// parm2: contract address
	address, err := vm.PopByteArray(engine)
	if err != nil {
		return err
	}
	addr, err := common.AddressParseFromBytes(address)
	if err != nil {
		return fmt.Errorf("invoke native contract:%s, address invalid", address)
	}
	// parm3: method
	method, err := vm.PopByteArray(engine)
	if err != nil {
		return err
	}
	if len(method) > METHOD_LENGTH_LIMIT {
		return fmt.Errorf("invoke native contract:%s method:%s too long, over max length 1024 limit", address, method)
	}
	// parm4: args
	args := vm.PopStackItem(engine)

	buf := new(bytes.Buffer)
	if err := BuildParamToNative(buf, args); err != nil {
		return err
	}

	// build contract
	contract := &states.Contract{
		Version: byte(version),
		Address: addr,
		Method:  string(method),
		Args:    buf.Bytes(),
	}

	bf := new(bytes.Buffer)
	if err := contract.Serialize(bf); err != nil {
		return err
	}

	// build nativeservice
	native := &native.NativeService{
		CloneCache: service.CloneCache,
		Code:       bf.Bytes(),
		Tx:         service.Tx,
		Height:     service.Height,
		Time:       service.Time,
		ContextRef: service.ContextRef,
		ServiceMap: make(map[string]native.Handler),
	}

	// native invoke
	result, err := native.Invoke()
	if err != nil {
		return err
	}
	vm.PushData(engine, result)
	return nil
}

func BuildParamToNative(bf *bytes.Buffer, item types.StackItems) error {
	switch item.(type) {
	case *types.ByteArray:
		a, _ := item.GetByteArray()
		if err := serialization.WriteVarBytes(bf, a); err != nil {
			return err
		}
	case *types.Integer:
		i, _ := item.GetByteArray()
		if err := serialization.WriteVarBytes(bf, i); err != nil {
			return err
		}
	case *types.Boolean:
		b, _ := item.GetBoolean()
		if err := serialization.WriteBool(bf, b); err != nil {
			return err
		}
	case *types.Array:
		arr, _ := item.GetArray()
		if err := serialization.WriteVarBytes(bf, types.BigIntToBytes(big.NewInt(int64(len(arr))))); err != nil {
			return err
		}
		for _, v := range arr {
			if err := BuildParamToNative(bf, v); err != nil {
				return err
			}
		}
	case *types.Struct:
		st, _ := item.GetStruct()
		for _, v := range st {
			if err := BuildParamToNative(bf, v); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("convert neovm params to native invalid type support: %s", reflect.TypeOf(item))
	}
	return nil
}
