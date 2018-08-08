

package mbg

import (
	"bytes"
	"math/big"

	"fmt"
	"github.com/mixbee/mixbee/common/constants"
	scommon "github.com/mixbee/mixbee/core/store/common"
	"github.com/mixbee/mixbee/errors"
	"github.com/mixbee/mixbee/smartcontract/service/native"
	"github.com/mixbee/mixbee/smartcontract/service/native/mbc"
	"github.com/mixbee/mixbee/smartcontract/service/native/utils"
	"github.com/mixbee/mixbee/vm/neovm/types"
)

func InitMbg() {
	native.Contracts[utils.MbgContractAddress] = RegisterMbgContract
}

func RegisterMbgContract(native *native.NativeService) {
	native.Register(mbc.INIT_NAME, MbgInit)
	native.Register(mbc.TRANSFER_NAME, MbgTransfer)
	native.Register(mbc.APPROVE_NAME, MbgApprove)
	native.Register(mbc.TRANSFERFROM_NAME, MbgTransferFrom)
	native.Register(mbc.NAME_NAME, MbgName)
	native.Register(mbc.SYMBOL_NAME, MbgSymbol)
	native.Register(mbc.DECIMALS_NAME, MbgDecimals)
	native.Register(mbc.TOTALSUPPLY_NAME, MbgTotalSupply)
	native.Register(mbc.BALANCEOF_NAME, MbgBalanceOf)
	native.Register(mbc.ALLOWANCE_NAME, MbgAllowance)
}

func MbgInit(native *native.NativeService) ([]byte, error) {
	contract := native.ContextRef.CurrentContext().ContractAddress
	amount, err := utils.GetStorageUInt64(native, mbc.GenTotalSupplyKey(contract))
	if err != nil {
		return utils.BYTE_FALSE, err
	}

	if amount > 0 {
		return utils.BYTE_FALSE, errors.NewErr("Init mbg has been completed!")
	}

	item := utils.GenUInt64StorageItem(constants.MBG_TOTAL_SUPPLY)
	native.CloneCache.Add(scommon.ST_STORAGE, mbc.GenTotalSupplyKey(contract), item)
	native.CloneCache.Add(scommon.ST_STORAGE, append(contract[:], utils.MbcContractAddress[:]...), item)
	mbc.AddNotifications(native, contract, &mbc.State{To: utils.MbcContractAddress, Value: constants.MBG_TOTAL_SUPPLY})
	return utils.BYTE_TRUE, nil
}

func MbgTransfer(native *native.NativeService) ([]byte, error) {
	transfers := new(mbc.Transfers)
	if err := transfers.Deserialize(bytes.NewBuffer(native.Input)); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[MbgTransfer] Transfers deserialize error!")
	}
	contract := native.ContextRef.CurrentContext().ContractAddress
	for _, v := range transfers.States {
		if v.Value == 0 {
			continue
		}
		if v.Value > constants.MBG_TOTAL_SUPPLY {
			return utils.BYTE_FALSE, fmt.Errorf("transfer mbg amount:%d over totalSupply:%d", v.Value, constants.MBG_TOTAL_SUPPLY)
		}
		if _, _, err := mbc.Transfer(native, contract, v); err != nil {
			return utils.BYTE_FALSE, err
		}
		mbc.AddNotifications(native, contract, v)
	}
	return utils.BYTE_TRUE, nil
}

func MbgApprove(native *native.NativeService) ([]byte, error) {
	state := new(mbc.State)
	if err := state.Deserialize(bytes.NewBuffer(native.Input)); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[MbgApprove] state deserialize error!")
	}
	if state.Value == 0 {
		return utils.BYTE_FALSE, nil
	}
	if state.Value > constants.MBG_TOTAL_SUPPLY {
		return utils.BYTE_FALSE, fmt.Errorf("approve mbg amount:%d over totalSupply:%d", state.Value, constants.MBG_TOTAL_SUPPLY)
	}
	if native.ContextRef.CheckWitness(state.From) == false {
		return utils.BYTE_FALSE, errors.NewErr("authentication failed!")
	}
	contract := native.ContextRef.CurrentContext().ContractAddress
	native.CloneCache.Add(scommon.ST_STORAGE, mbc.GenApproveKey(contract, state.From, state.To), utils.GenUInt64StorageItem(state.Value))
	return utils.BYTE_TRUE, nil
}

func MbgTransferFrom(native *native.NativeService) ([]byte, error) {
	state := new(mbc.TransferFrom)
	if err := state.Deserialize(bytes.NewBuffer(native.Input)); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[MbcTransferFrom] State deserialize error!")
	}
	if state.Value == 0 {
		return utils.BYTE_FALSE, nil
	}
	if state.Value > constants.MBG_TOTAL_SUPPLY {
		return utils.BYTE_FALSE, fmt.Errorf("approve mbg amount:%d over totalSupply:%d", state.Value, constants.MBG_TOTAL_SUPPLY)
	}
	contract := native.ContextRef.CurrentContext().ContractAddress
	if _, _, err := mbc.TransferedFrom(native, contract, state); err != nil {
		return utils.BYTE_FALSE, err
	}
	mbc.AddNotifications(native, contract, &mbc.State{From: state.From, To: state.To, Value: state.Value})
	return utils.BYTE_TRUE, nil
}

func MbgName(native *native.NativeService) ([]byte, error) {
	return []byte(constants.MBG_NAME), nil
}

func MbgDecimals(native *native.NativeService) ([]byte, error) {
	return big.NewInt(int64(constants.MBG_DECIMALS)).Bytes(), nil
}

func MbgSymbol(native *native.NativeService) ([]byte, error) {
	return []byte(constants.MBG_SYMBOL), nil
}

func MbgTotalSupply(native *native.NativeService) ([]byte, error) {
	contract := native.ContextRef.CurrentContext().ContractAddress
	amount, err := utils.GetStorageUInt64(native, mbc.GenTotalSupplyKey(contract))
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[MbcTotalSupply] get totalSupply error!")
	}
	return types.BigIntToBytes(big.NewInt(int64(amount))), nil
}

func MbgBalanceOf(native *native.NativeService) ([]byte, error) {
	return mbc.GetBalanceValue(native, mbc.TRANSFER_FLAG)
}

func MbgAllowance(native *native.NativeService) ([]byte, error) {
	return mbc.GetBalanceValue(native, mbc.APPROVE_FLAG)
}
