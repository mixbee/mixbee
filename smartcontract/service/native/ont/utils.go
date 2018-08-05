

package ont

import (
	"bytes"
	"fmt"
	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/common/config"
	"github.com/mixbee/mixbee/common/serialization"
	cstates "github.com/mixbee/mixbee/core/states"
	scommon "github.com/mixbee/mixbee/core/store/common"
	"github.com/mixbee/mixbee/errors"
	"github.com/mixbee/mixbee/smartcontract/event"
	"github.com/mixbee/mixbee/smartcontract/service/native"
	"github.com/mixbee/mixbee/smartcontract/service/native/utils"
)

const (
	UNBOUND_TIME_OFFSET = "unboundTimeOffset"
	TOTAL_SUPPLY_NAME   = "totalSupply"
	INIT_NAME           = "init"
	TRANSFER_NAME       = "transfer"
	APPROVE_NAME        = "approve"
	TRANSFERFROM_NAME   = "transferFrom"
	NAME_NAME           = "name"
	SYMBOL_NAME         = "symbol"
	DECIMALS_NAME       = "decimals"
	TOTALSUPPLY_NAME    = "totalSupply"
	BALANCEOF_NAME      = "balanceOf"
	ALLOWANCE_NAME      = "allowance"

	SET_KEY             = "setkey"
	GET_KEY             = "getkey"
)

func AddNotifications(native *native.NativeService, contract common.Address, state *State) {
	if !config.DefConfig.Common.EnableEventLog {
		return
	}
	native.Notifications = append(native.Notifications,
		&event.NotifyEventInfo{
			ContractAddress: contract,
			States:          []interface{}{TRANSFER_NAME, state.From.ToBase58(), state.To.ToBase58(), state.Value},
		})
}

func GetToUInt64StorageItem(toBalance, value uint64) *cstates.StorageItem {
	bf := new(bytes.Buffer)
	serialization.WriteUint64(bf, toBalance+value)
	return &cstates.StorageItem{Value: bf.Bytes()}
}

func GenTotalSupplyKey(contract common.Address) []byte {
	return append(contract[:], TOTAL_SUPPLY_NAME...)
}

func GenBalanceKey(contract, addr common.Address) []byte {
	return append(contract[:], addr[:]...)
}

func Transfer(native *native.NativeService, contract common.Address, state *State) (uint64, uint64, error) {
	if !native.ContextRef.CheckWitness(state.From) {
		return 0, 0, errors.NewErr("authentication failed!")
	}

	fromBalance, err := FromTransfer(native, GenBalanceKey(contract, state.From), state.Value)
	if err != nil {
		return 0, 0, err
	}

	toBalance, err := ToTransfer(native, GenBalanceKey(contract, state.To), state.Value)
	if err != nil {
		return 0, 0, err
	}
	return fromBalance, toBalance, nil
}

func TransferForCrossChainContract(native *native.NativeService, contract common.Address, state *State) (uint64, uint64, error) {

	fromBalance, err := FromTransfer(native, GenBalanceKey(contract, state.From), state.Value)
	if err != nil {
		return 0, 0, err
	}

	toBalance, err := ToTransfer(native, GenBalanceKey(contract, state.To), state.Value)
	if err != nil {
		return 0, 0, err
	}
	return fromBalance, toBalance, nil
}

func GenApproveKey(contract, from, to common.Address) []byte {
	temp := append(contract[:], from[:]...)
	return append(temp, to[:]...)
}

func TransferedFrom(native *native.NativeService, currentContract common.Address, state *TransferFrom) (uint64, uint64, error) {
	if native.ContextRef.CheckWitness(state.Sender) == false {
		return 0, 0, errors.NewErr("authentication failed!")
	}

	if err := fromApprove(native, genTransferFromKey(currentContract, state), state.Value); err != nil {
		return 0, 0, err
	}

	fromBalance, err := FromTransfer(native, GenBalanceKey(currentContract, state.From), state.Value)
	if err != nil {
		return 0, 0, err
	}

	toBalance, err := ToTransfer(native, GenBalanceKey(currentContract, state.To), state.Value)
	if err != nil {
		return 0, 0, err
	}
	return fromBalance, toBalance, nil
}

func getUnboundOffset(native *native.NativeService, contract, address common.Address) (uint32, error) {
	offset, err := utils.GetStorageUInt32(native, genAddressUnboundOffsetKey(contract, address))
	if err != nil {
		return 0, err
	}
	return offset, nil
}

func genTransferFromKey(contract common.Address, state *TransferFrom) []byte {
	temp := append(contract[:], state.From[:]...)
	return append(temp, state.Sender[:]...)
}

func fromApprove(native *native.NativeService, fromApproveKey []byte, value uint64) error {
	approveValue, err := utils.GetStorageUInt64(native, fromApproveKey)
	if err != nil {
		return err
	}
	if approveValue < value {
		return fmt.Errorf("[TransferFrom] approve balance insufficient! have %d, got %d", approveValue, value)
	} else if approveValue == value {
		native.CloneCache.Delete(scommon.ST_STORAGE, fromApproveKey)
	} else {
		native.CloneCache.Add(scommon.ST_STORAGE, fromApproveKey, utils.GenUInt64StorageItem(approveValue-value))
	}
	return nil
}

func FromTransfer(native *native.NativeService, fromKey []byte, value uint64) (uint64, error) {
	fromBalance, err := utils.GetStorageUInt64(native, fromKey)
	if err != nil {
		return 0, err
	}
	if fromBalance < value {
		addr, _ := common.AddressParseFromBytes(fromKey[20:])
		return 0, fmt.Errorf("[Transfer] balance insufficient. contract:%s, account:%s,balance:%d, transfer amount:%d",
			native.ContextRef.CurrentContext().ContractAddress.ToHexString(), addr.ToBase58(), fromBalance, value)
	} else if fromBalance == value {
		native.CloneCache.Delete(scommon.ST_STORAGE, fromKey)
	} else {
		native.CloneCache.Add(scommon.ST_STORAGE, fromKey, utils.GenUInt64StorageItem(fromBalance-value))
	}
	return fromBalance, nil
}

func ToTransfer(native *native.NativeService, toKey []byte, value uint64) (uint64, error) {
	toBalance, err := utils.GetStorageUInt64(native, toKey)
	if err != nil {
		return 0, err
	}
	native.CloneCache.Add(scommon.ST_STORAGE, toKey, GetToUInt64StorageItem(toBalance, value))
	return toBalance, nil
}

func genAddressUnboundOffsetKey(contract, address common.Address) []byte {
	temp := append(contract[:], UNBOUND_TIME_OFFSET...)
	return append(temp, address[:]...)
}
