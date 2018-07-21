

package global_params

import (
	"bytes"

	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/common/config"
	cstates "github.com/mixbee/mixbee/core/states"
	"github.com/mixbee/mixbee/smartcontract/event"
	"github.com/mixbee/mixbee/smartcontract/service/native"
	"github.com/mixbee/mixbee/smartcontract/service/native/utils"
)

const (
	PARAM    = "param"
	TRANSFER = "transfer"
	ADMIN    = "admin"
	OPERATOR = "operator"
)

func getRoleStorageItem(role common.Address) *cstates.StorageItem {
	bf := new(bytes.Buffer)
	utils.WriteAddress(bf, role)
	return &cstates.StorageItem{Value: bf.Bytes()}
}

func getParamStorageItem(params Params) *cstates.StorageItem {
	bf := new(bytes.Buffer)
	params.Serialize(bf)
	return &cstates.StorageItem{Value: bf.Bytes()}
}

func generateParamKey(contract common.Address, valueType paramType) []byte {
	key := append(contract[:], PARAM...)
	key = append(key[:], byte(valueType))
	return key
}

func generateAdminKey(contract common.Address, isTransferAdmin bool) []byte {
	if isTransferAdmin {
		return append(contract[:], TRANSFER...)
	} else {
		return append(contract[:], ADMIN...)
	}
}

func GenerateOperatorKey(contract common.Address) []byte {
	return append(contract[:], OPERATOR...)
}

func getStorageParam(native *native.NativeService, key []byte) (Params, error) {
	item, err := utils.GetStorageItem(native, key)
	params := Params{}
	if err != nil || item == nil {
		return params, err
	}
	bf := bytes.NewBuffer(item.Value)
	err = params.Deserialize(bf)
	return params, err
}

func GetStorageRole(native *native.NativeService, key []byte) (common.Address, error) {
	item, err := utils.GetStorageItem(native, key)
	var role common.Address
	if err != nil || item == nil {
		return role, err
	}
	bf := bytes.NewBuffer(item.Value)
	role, err = utils.ReadAddress(bf)
	return role, err
}

func NotifyRoleChange(native *native.NativeService, contract common.Address, functionName string,
	newAddr common.Address) {
	if !config.DefConfig.Common.EnableEventLog {
		return
	}
	native.Notifications = append(native.Notifications,
		&event.NotifyEventInfo{
			ContractAddress: contract,
			States:          []interface{}{functionName, newAddr.ToBase58()},
		})
}

func NotifyTransferAdmin(native *native.NativeService, contract common.Address, functionName string,
	originAdmin, newAdmin common.Address) {
	if !config.DefConfig.Common.EnableEventLog {
		return
	}
	native.Notifications = append(native.Notifications,
		&event.NotifyEventInfo{
			ContractAddress: contract,
			States:          []interface{}{functionName, originAdmin.ToBase58(), newAdmin.ToBase58()},
		})
}

func NotifyParamChange(native *native.NativeService, contract common.Address, functionName string, params Params) {
	if !config.DefConfig.Common.EnableEventLog {
		return
	}
	paramsString := ""
	for _, param := range params {
		paramsString += param.Key + "," + param.Value + ";"
	}
	paramsString = paramsString[:len(paramsString)-1]
	native.Notifications = append(native.Notifications,
		&event.NotifyEventInfo{
			ContractAddress: contract,
			States:          []interface{}{functionName, paramsString},
		})
}
