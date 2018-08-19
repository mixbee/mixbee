

package neovm

import (
	"bytes"
	"io"
	"math/big"
	"reflect"
	"sort"

	"github.com/mixbee/mixbee-crypto/keypair"
	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/common/serialization"
	"github.com/mixbee/mixbee/core/types"
	"github.com/mixbee/mixbee/errors"
	scommon "github.com/mixbee/mixbee/smartcontract/common"
	"github.com/mixbee/mixbee/smartcontract/event"
	vm "github.com/mixbee/mixbee/vm/neovm"
	vmtypes "github.com/mixbee/mixbee/vm/neovm/types"
)

// HeaderGetNextConsensus put current block time to vm stack
func RuntimeGetTime(service *NeoVmService, engine *vm.ExecutionEngine) error {
	vm.PushData(engine, int(service.Time))
	return nil
}

// RuntimeCheckWitness provide check permissions service
// If param address isn't exist in authorization list, check fail
func RuntimeCheckWitness(service *NeoVmService, engine *vm.ExecutionEngine) error {
	data, err := vm.PopByteArray(engine)
	if err != nil {
		return err
	}
	var result bool
	if len(data) == 20 {
		address, err := common.AddressParseFromBytes(data)
		if err != nil {
			return err
		}
		result = service.ContextRef.CheckWitness(address)
	} else {
		// parse the byte sequencce to a public key
		pk, err := keypair.DeserializePublicKey(data)
		if err != nil {
			return errors.NewDetailErr(err, errors.ErrNoCode, "[RuntimeCheckWitness] data invalid.")
		}
		result = service.ContextRef.CheckWitness(types.AddressFromPubKey(pk))
	}

	vm.PushData(engine, result)
	return nil
}

func RuntimeSerialize(service *NeoVmService, engine *vm.ExecutionEngine) error {
	item := vm.PopStackItem(engine)

	buf, err := SerializeStackItem(item)
	if err != nil {
		return err
	}
	vm.PushData(engine, buf)
	return nil
}

func RuntimeDeserialize(service *NeoVmService, engine *vm.ExecutionEngine) error {
	data, err := vm.PopByteArray(engine)
	if err != nil {
		return err
	}
	bf := bytes.NewBuffer(data)
	item, err := DeserializeStackItem(bf)
	if err != nil {
		return err
	}

	if item == nil {
		return nil
	}
	vm.PushData(engine, item)
	return nil
}

// RuntimeNotify put smart contract execute event notify to notifications
func RuntimeNotify(service *NeoVmService, engine *vm.ExecutionEngine) error {
	item := vm.PopStackItem(engine)
	context := service.ContextRef.CurrentContext()
	service.Notifications = append(service.Notifications, &event.NotifyEventInfo{ContractAddress: context.ContractAddress, States: scommon.ConvertNeoVmTypeHexString(item)})
	return nil
}

// RuntimeLog push smart contract execute event log to client
func RuntimeLog(service *NeoVmService, engine *vm.ExecutionEngine) error {
	item, err := vm.PopByteArray(engine)
	if err != nil {
		return err
	}
	context := service.ContextRef.CurrentContext()
	txHash := service.Tx.Hash()
	event.PushSmartCodeEvent(txHash, 0, event.EVENT_LOG, &event.LogEventArgs{TxHash: txHash, ContractAddress: context.ContractAddress, Message: string(item)})
	return nil
}

func RuntimeGetTrigger(service *NeoVmService, engine *vm.ExecutionEngine) error {
	vm.PushData(engine, 0)
	return nil
}

func SerializeStackItem(item vmtypes.StackItems) ([]byte, error) {
	if CircularRefAndDepthDetection(item) {
		return nil, errors.NewErr("runtime serialize: can not serialize circular reference data")
	}

	bf := new(bytes.Buffer)
	err := serializeStackItem(item, common.NewLimitedWriter(bf, uint64(vm.MAX_BYTEARRAY_SIZE)))
	if err != nil {
		return nil, err
	}

	return bf.Bytes(), nil
}

func serializeStackItem(item vmtypes.StackItems, w io.Writer) error {
	switch item.(type) {
	case *vmtypes.ByteArray:
		if err := serialization.WriteByte(w, vmtypes.ByteArrayType); err != nil {
			return errors.NewErr("Serialize ByteArray stackItems error: " + err.Error())
		}
		ba, _ := item.GetByteArray()
		if err := serialization.WriteVarBytes(w, ba); err != nil {
			return errors.NewErr("Serialize ByteArray stackItems error: " + err.Error())
		}

	case *vmtypes.Boolean:
		if err := serialization.WriteByte(w, vmtypes.BooleanType); err != nil {
			return errors.NewErr("Serialize Boolean StackItems error: " + err.Error())
		}
		b, _ := item.GetBoolean()
		if err := serialization.WriteBool(w, b); err != nil {
			return errors.NewErr("Serialize Boolean stackItems error: " + err.Error())
		}

	case *vmtypes.Integer:
		if err := serialization.WriteByte(w, vmtypes.IntegerType); err != nil {
			return errors.NewErr("Serialize Integer stackItems error: " + err.Error())
		}
		i, _ := item.GetByteArray()
		if err := serialization.WriteVarBytes(w, i); err != nil {
			return errors.NewErr("Serialize Integer stackItems error: " + err.Error())
		}

	case *vmtypes.Array:
		if err := serialization.WriteByte(w, vmtypes.ArrayType); err != nil {
			return errors.NewErr("Serialize Array stackItems error: " + err.Error())
		}
		a, _ := item.GetArray()
		if err := serialization.WriteVarUint(w, uint64(len(a))); err != nil {
			return errors.NewErr("Serialize Array stackItems error: " + err.Error())
		}

		for _, v := range a {
			err := serializeStackItem(v, w)
			if err != nil {
				return err
			}
		}

	case *vmtypes.Struct:
		if err := serialization.WriteByte(w, vmtypes.StructType); err != nil {
			return errors.NewErr("Serialize Struct stackItems error: " + err.Error())
		}
		s, _ := item.GetStruct()
		if err := serialization.WriteVarUint(w, uint64(len(s))); err != nil {
			return errors.NewErr("Serialize Struct stackItems error: " + err.Error())
		}

		for _, v := range s {
			err := serializeStackItem(v, w)
			if err != nil {
				return err
			}
		}

	case *vmtypes.Map:
		var unsortKey []string
		mp, _ := item.GetMap()
		keyMap := make(map[string]vmtypes.StackItems, 0)

		if err := serialization.WriteByte(w, vmtypes.MapType); err != nil {
			return errors.NewErr("Serialize Map stackItems error: " + err.Error())
		}
		if err := serialization.WriteVarUint(w, uint64(len(mp))); err != nil {
			return errors.NewErr("Serialize Map stackItems error: " + err.Error())
		}

		for k := range mp {
			switch k.(type) {
			case *vmtypes.ByteArray, *vmtypes.Integer:
				ba, _ := k.GetByteArray()
				key := string(ba)
				if key == "" {
					return errors.NewErr("Serialize Map error: invalid key type")
				}
				unsortKey = append(unsortKey, key)
				keyMap[key] = k

			default:
				return errors.NewErr("Unsupport map key type.")
			}
		}

		sort.Strings(unsortKey)
		for _, v := range unsortKey {
			key := keyMap[v]
			err := serializeStackItem(key, w)
			if err != nil {
				return err
			}
			err = serializeStackItem(mp[key], w)
			if err != nil {
				return err
			}
		}

	default:
		return errors.NewErr("unknown type")
	}

	return nil
}

func DeserializeStackItem(r io.Reader) (items vmtypes.StackItems, err error) {
	t, err := serialization.ReadByte(r)
	if err != nil {
		return nil, errors.NewErr("Deserialize error: " + err.Error())
	}

	switch t {
	case vmtypes.ByteArrayType:
		b, err := serialization.ReadVarBytes(r)
		if err != nil {
			return nil, errors.NewErr("Deserialize stackItems ByteArray error: " + err.Error())
		}
		return vmtypes.NewByteArray(b), nil

	case vmtypes.BooleanType:
		b, err := serialization.ReadBool(r)
		if err != nil {
			return nil, errors.NewErr("Deserialize stackItems Boolean error: " + err.Error())
		}
		return vmtypes.NewBoolean(b), nil

	case vmtypes.IntegerType:
		b, err := serialization.ReadVarBytes(r)
		if err != nil {
			return nil, errors.NewErr("Deserialize stackItems Integer error: " + err.Error())
		}
		return vmtypes.NewInteger(new(big.Int).SetBytes(b)), nil

	case vmtypes.ArrayType, vmtypes.StructType:
		count, err := serialization.ReadVarUint(r, 0)
		if err != nil {
			return nil, errors.NewErr("Deserialize stackItems error: " + err.Error())
		}

		var arr []vmtypes.StackItems
		for count > 0 {
			item, err := DeserializeStackItem(r)
			if err != nil {
				return nil, err
			}
			arr = append(arr, item)
			count--
		}

		if t == vmtypes.StructType {
			return vmtypes.NewStruct(arr), nil
		}

		return vmtypes.NewArray(arr), nil

	case vmtypes.MapType:
		count, err := serialization.ReadVarUint(r, 0)
		if err != nil {
			return nil, errors.NewErr("Deserialize stackItems map error: " + err.Error())
		}

		mp := vmtypes.NewMap()
		for count > 0 {
			key, err := DeserializeStackItem(r)
			if err != nil {
				return nil, err
			}

			value, err := DeserializeStackItem(r)
			if err != nil {
				return nil, err
			}
			m, _ := mp.GetMap()
			m[key] = value
			count--
		}
		return mp, nil

	default:
		return nil, errors.NewErr("unknown type")
	}

	return nil, nil
}

func CircularRefAndDepthDetection(value vmtypes.StackItems) bool {
	return circularRefAndDepthDetection(value, make(map[uintptr]bool), 0)
}

func circularRefAndDepthDetection(value vmtypes.StackItems, visited map[uintptr]bool, depth int) bool {
	if depth > vmtypes.MAX_STRCUT_DEPTH {
		return true
	}
	switch value.(type) {
	case *vmtypes.Array:
		a, _ := value.GetArray()
		if len(a) == 0 {
			return false
		}

		p := reflect.ValueOf(a).Pointer()
		if visited[p] {
			return true
		}
		visited[p] = true

		for _, v := range a {
			if circularRefAndDepthDetection(v, visited, depth+1) {
				return true
			}
		}

		delete(visited, p)
		return false
	case *vmtypes.Struct:
		s, _ := value.GetStruct()
		if len(s) == 0 {
			return false
		}

		p := reflect.ValueOf(s).Pointer()
		if visited[p] {
			return true
		}
		visited[p] = true

		for _, v := range s {
			if circularRefAndDepthDetection(v, visited, depth+1) {
				return true
			}
		}

		delete(visited, p)
		return false
	case *vmtypes.Map:
		mp, _ := value.GetMap()

		p := reflect.ValueOf(mp).Pointer()
		if visited[p] {
			return true
		}
		visited[p] = true

		for k, v := range mp {
			if circularRefAndDepthDetection(k, visited, depth+1) {
				return true
			}
			if circularRefAndDepthDetection(v, visited, depth+1) {
				return true
			}
		}

		delete(visited, p)
		return false
	default:
		return false
	}

	return false
}
