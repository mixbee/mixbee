

package mbc

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/common/constants"
	"github.com/mixbee/mixbee/common/log"
	"github.com/mixbee/mixbee/common/serialization"
	scommon "github.com/mixbee/mixbee/core/store/common"
	"github.com/mixbee/mixbee/errors"
	"github.com/mixbee/mixbee/smartcontract/service/native"
	"github.com/mixbee/mixbee/smartcontract/service/native/utils"
	"github.com/mixbee/mixbee/vm/neovm/types"
)

const (
	TRANSFER_FLAG byte = 1
	APPROVE_FLAG  byte = 2
)

func InitMbc() {
	native.Contracts[utils.MbcContractAddress] = RegisterMbcContract
}

func RegisterMbcContract(native *native.NativeService) {
	native.Register(INIT_NAME, MbcInit)
	native.Register(TRANSFER_NAME, MbcTransfer)
	native.Register(APPROVE_NAME, MbcApprove)
	native.Register(TRANSFERFROM_NAME, MbcTransferFrom)
	native.Register(NAME_NAME, MbcName)
	native.Register(SYMBOL_NAME, MbcSymbol)
	native.Register(DECIMALS_NAME, MbcDecimals)
	native.Register(TOTALSUPPLY_NAME, MbcTotalSupply)
	native.Register(BALANCEOF_NAME, MbcBalanceOf)
	native.Register(ALLOWANCE_NAME, MbcAllowance)
}

func MbcInit(native *native.NativeService) ([]byte, error) {
	contract := native.ContextRef.CurrentContext().ContractAddress
	amount, err := utils.GetStorageUInt64(native, GenTotalSupplyKey(contract))
	if err != nil {
		return utils.BYTE_FALSE, err
	}

	if amount > 0 {
		return utils.BYTE_FALSE, errors.NewErr("Init mbc has been completed!")
	}

	distribute := make(map[common.Address]uint64)
	buf, err := serialization.ReadVarBytes(bytes.NewBuffer(native.Input))
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "serialization.ReadVarBytes, contract params deserialize error!")
	}
	input := bytes.NewBuffer(buf)
	num, err := utils.ReadVarUint(input)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("read number error:%v", err)
	}
	sum := uint64(0)
	overflow := false
	for i := uint64(0); i < num; i++ {
		addr, err := utils.ReadAddress(input)
		if err != nil {
			return utils.BYTE_FALSE, fmt.Errorf("read address error:%v", err)
		}
		value, err := utils.ReadVarUint(input)
		if err != nil {
			return utils.BYTE_FALSE, fmt.Errorf("read value error:%v", err)
		}
		sum, overflow = common.SafeAdd(sum, value)
		if overflow {
			return utils.BYTE_FALSE, errors.NewErr("wrong config. overflow detected")
		}
		distribute[addr] += value
	}
	if sum != constants.MBC_TOTAL_SUPPLY {
		return utils.BYTE_FALSE, fmt.Errorf("wrong config. total supply %d != %d", sum, constants.MBC_TOTAL_SUPPLY)
	}

	for addr, val := range distribute {
		balanceKey := GenBalanceKey(contract, addr)
		item := utils.GenUInt64StorageItem(val)
		native.CloneCache.Add(scommon.ST_STORAGE, balanceKey, item)
		AddNotifications(native, contract, &State{To: addr, Value: val})
	}
	native.CloneCache.Add(scommon.ST_STORAGE, GenTotalSupplyKey(contract), utils.GenUInt64StorageItem(constants.MBC_TOTAL_SUPPLY))

	return utils.BYTE_TRUE, nil
}

func MbcTransfer(native *native.NativeService) ([]byte, error) {
	transfers := new(Transfers)
	if err := transfers.Deserialize(bytes.NewBuffer(native.Input)); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[Transfer] Transfers deserialize error!")
	}
	contract := native.ContextRef.CurrentContext().ContractAddress
	for _, v := range transfers.States {
		if v.Value == 0 {
			continue
		}
		if v.Value > constants.MBC_TOTAL_SUPPLY {
			return utils.BYTE_FALSE, fmt.Errorf("transfer mbc amount:%d over totalSupply:%d", v.Value, constants.MBC_TOTAL_SUPPLY)
		}
		fromBalance, toBalance, err := Transfer(native, contract, v)
		if err != nil {
			return utils.BYTE_FALSE, err
		}

		if err := grantMbg(native, contract, v.From, fromBalance); err != nil {
			return utils.BYTE_FALSE, err
		}

		if err := grantMbg(native, contract, v.To, toBalance); err != nil {
			return utils.BYTE_FALSE, err
		}

		AddNotifications(native, contract, v)
	}
	return utils.BYTE_TRUE, nil
}

func MbcTransferFrom(native *native.NativeService) ([]byte, error) {
	state := new(TransferFrom)
	if err := state.Deserialize(bytes.NewBuffer(native.Input)); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[MbcTransferFrom] State deserialize error!")
	}
	if state.Value == 0 {
		return utils.BYTE_FALSE, nil
	}
	if state.Value > constants.MBC_TOTAL_SUPPLY {
		return utils.BYTE_FALSE, fmt.Errorf("transferFrom mbc amount:%d over totalSupply:%d", state.Value, constants.MBC_TOTAL_SUPPLY)
	}
	contract := native.ContextRef.CurrentContext().ContractAddress
	fromBalance, toBalance, err := TransferedFrom(native, contract, state)
	if err != nil {
		return utils.BYTE_FALSE, err
	}
	if err := grantMbg(native, contract, state.From, fromBalance); err != nil {
		return utils.BYTE_FALSE, err
	}
	if err := grantMbg(native, contract, state.To, toBalance); err != nil {
		return utils.BYTE_FALSE, err
	}
	AddNotifications(native, contract, &State{From: state.From, To: state.To, Value: state.Value})
	return utils.BYTE_TRUE, nil
}

func MbcApprove(native *native.NativeService) ([]byte, error) {
	state := new(State)
	if err := state.Deserialize(bytes.NewBuffer(native.Input)); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[MbgApprove] state deserialize error!")
	}
	if state.Value == 0 {
		return utils.BYTE_FALSE, nil
	}
	if state.Value > constants.MBC_TOTAL_SUPPLY {
		return utils.BYTE_FALSE, fmt.Errorf("approve mbc amount:%d over totalSupply:%d", state.Value, constants.MBC_TOTAL_SUPPLY)
	}
	if native.ContextRef.CheckWitness(state.From) == false {
		return utils.BYTE_FALSE, errors.NewErr("authentication failed!")
	}
	contract := native.ContextRef.CurrentContext().ContractAddress
	native.CloneCache.Add(scommon.ST_STORAGE, GenApproveKey(contract, state.From, state.To), utils.GenUInt64StorageItem(state.Value))
	return utils.BYTE_TRUE, nil
}

func MbcName(native *native.NativeService) ([]byte, error) {
	return []byte(constants.MBC_NAME), nil
}

func MbcDecimals(native *native.NativeService) ([]byte, error) {
	return types.BigIntToBytes(big.NewInt(int64(constants.MBC_DECIMALS))), nil
}

func MbcSymbol(native *native.NativeService) ([]byte, error) {
	return []byte(constants.MBC_SYMBOL), nil
}

func MbcTotalSupply(native *native.NativeService) ([]byte, error) {
	contract := native.ContextRef.CurrentContext().ContractAddress
	amount, err := utils.GetStorageUInt64(native, GenTotalSupplyKey(contract))
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[MbcTotalSupply] get totalSupply error!")
	}
	return types.BigIntToBytes(big.NewInt(int64(amount))), nil
}

func MbcBalanceOf(native *native.NativeService) ([]byte, error) {
	return GetBalanceValue(native, TRANSFER_FLAG)
}

func MbcAllowance(native *native.NativeService) ([]byte, error) {
	return GetBalanceValue(native, APPROVE_FLAG)
}

func GetBalanceValue(native *native.NativeService, flag byte) ([]byte, error) {
	var key []byte
	buf := bytes.NewBuffer(native.Input)
	from, err := utils.ReadAddress(buf)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[GetBalanceValue] get from address error!")
	}
	contract := native.ContextRef.CurrentContext().ContractAddress
	if flag == APPROVE_FLAG {
		to, err := utils.ReadAddress(buf)
		if err != nil {
			return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[GetBalanceValue] get from address error!")
		}
		key = GenApproveKey(contract, from, to)
	} else if flag == TRANSFER_FLAG {
		key = GenBalanceKey(contract, from)
	}
	amount, err := utils.GetStorageUInt64(native, key)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[GetBalanceValue] address parse error!")
	}
	return types.BigIntToBytes(big.NewInt(int64(amount))), nil
}

func grantMbg(native *native.NativeService, contract, address common.Address, balance uint64) error {
	startOffset, err := getUnboundOffset(native, contract, address)
	if err != nil {
		return err
	}
	if native.Time <= constants.GENESIS_BLOCK_TIMESTAMP {
		return nil
	}
	endOffset := native.Time - constants.GENESIS_BLOCK_TIMESTAMP
	if endOffset < startOffset {
		errstr := fmt.Sprintf("grant Mbg error: wrong timestamp endOffset: %d < startOffset: %d", endOffset, startOffset)
		log.Error(errstr)
		return errors.NewErr(errstr)
	} else if endOffset == startOffset {
		return nil
	}

	if balance != 0 {
		value := utils.CalcUnbindMbg(balance, startOffset, endOffset)

		args, err := getApproveArgs(native, contract, utils.MbgContractAddress, address, value)
		if err != nil {
			return err
		}

		if _, err := native.NativeCall(utils.MbgContractAddress, "approve", args); err != nil {
			return err
		}
	}

	native.CloneCache.Add(scommon.ST_STORAGE, genAddressUnboundOffsetKey(contract, address), utils.GenUInt32StorageItem(endOffset))
	return nil
}

func getApproveArgs(native *native.NativeService, contract, ongContract, address common.Address, value uint64) ([]byte, error) {
	bf := new(bytes.Buffer)
	approve := &State{
		From:  contract,
		To:    address,
		Value: value,
	}

	stateValue, err := utils.GetStorageUInt64(native, GenApproveKey(ongContract, approve.From, approve.To))
	if err != nil {
		return nil, err
	}

	approve.Value += stateValue

	if err := approve.Serialize(bf); err != nil {
		return nil, err
	}
	return bf.Bytes(), nil
}
