package native

import (
	"bytes"
	"fmt"
	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/core/types"
	"github.com/mixbee/mixbee/errors"
	"github.com/mixbee/mixbee/smartcontract/context"
	"github.com/mixbee/mixbee/smartcontract/event"
	"github.com/mixbee/mixbee/smartcontract/states"
	sstates "github.com/mixbee/mixbee/smartcontract/states"
	"github.com/mixbee/mixbee/smartcontract/storage"
)

type (
	Handler func(native *NativeService) ([]byte, error)
	RegisterService func(native *NativeService)
)

var (
	Contracts = make(map[common.Address]RegisterService)
)

// Native service struct
// Invoke a native smart contract, new a native service
type NativeService struct {
	CloneCache    *storage.CloneCache
	ServiceMap    map[string]Handler
	Notifications []*event.NotifyEventInfo
	Code          []byte
	Input         []byte
	Tx            *types.Transaction
	Height        uint32
	Time          uint32
	ContextRef    context.ContextRef
}

func (this *NativeService) Register(methodName string, handler Handler) {
	this.ServiceMap[methodName] = handler
}

func (this *NativeService) Invoke() (interface{}, error) {
	bf := bytes.NewBuffer(this.Code)
	contract := new(sstates.Contract)
	if err := contract.Deserialize(bf); err != nil {
		return false, err
	}
	services, ok := Contracts[contract.Address]
	if !ok {
		return false, fmt.Errorf("Native contract address %x haven't been registered.", contract.Address)
	}

	if this.Tx.SystemTx == 1 && !IsSystemTx(contract.Address, contract.Method) {
		return false, fmt.Errorf("Native contract address %x method %x not support systemTx.", contract.Address, contract.Method)
	}

	services(this)
	service, ok := this.ServiceMap[contract.Method]
	if !ok {
		return false, fmt.Errorf("Native contract %x doesn't support this function %s.",
			contract.Address, contract.Method)
	}
	args := this.Input
	this.Input = contract.Args
	this.ContextRef.PushContext(&context.Context{ContractAddress: contract.Address})
	notifications := this.Notifications
	this.Notifications = []*event.NotifyEventInfo{}
	result, err := service(this)
	if err != nil {
		return result, errors.NewDetailErr(err, errors.ErrNoCode, "[Invoke] Native serivce function execute error!")
	}
	this.ContextRef.PopContext()
	this.ContextRef.PushNotifications(this.Notifications)
	this.Notifications = notifications
	this.Input = args
	return result, nil
}

func (this *NativeService) NativeCall(address common.Address, method string, args []byte) (interface{}, error) {
	bf := new(bytes.Buffer)
	c := states.Contract{
		Address: address,
		Method:  method,
		Args:    args,
	}
	if err := c.Serialize(bf); err != nil {
		return nil, err
	}
	this.Code = bf.Bytes()
	return this.Invoke()
}

func IsSystemTx(address common.Address, method string) bool {

	CrossChainContractAddress, _ := common.AddressParseFromBytes([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x09})

	if address == CrossChainContractAddress && (method == "crossRelease" || method == "crossUnlock") {
		return true
	}
	CrossChainPairEvidenceContractAddress, _ := common.AddressParseFromBytes([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10})
	if address == CrossChainPairEvidenceContractAddress && method == "pushEvidence" {
		return true
	}

	CrossChainVerifynodeContractAddress, _ := common.AddressParseFromBytes([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x11})
	if address == CrossChainVerifynodeContractAddress && method == "registerVerifyNode" {
		return true
	}

	return false
}
