

package mixtest

import (
	"github.com/mixbee/mixbee/common/constants"
	scommon "github.com/mixbee/mixbee/core/store/common"
	"github.com/mixbee/mixbee/errors"
	cstates "github.com/mixbee/mixbee/core/states"
	"github.com/mixbee/mixbee/smartcontract/service/native"
	"github.com/mixbee/mixbee/smartcontract/service/native/mbc"
	"github.com/mixbee/mixbee/smartcontract/service/native/utils"
	"bytes"
	"github.com/mixbee/mixbee/common/log"
	"github.com/mixbee/mixbee/common/serialization"
	"fmt"
	"github.com/mixbee/mixbee/common"
)

func InitMixTest() {
	native.Contracts[utils.MixTestContractAddress] = RegisterMixTestContract
}

func RegisterMixTestContract(native *native.NativeService) {
	native.Register(mbc.INIT_NAME, MixTestInit)
	native.Register(mbc.SET_KEY, MixTestSetKey)
	native.Register(mbc.NAME_NAME, MixTestName)
	native.Register(mbc.SYMBOL_NAME, MixTestSymbol)
	native.Register(mbc.GET_KEY, MixTestGetKey)
}

func MixTestInit(native *native.NativeService) ([]byte, error) {
	contract := native.ContextRef.CurrentContext().ContractAddress
	_ = contract.ToHexString
	item := &cstates.StorageItem{Value: []byte("MIXTEST")}
	native.CloneCache.Add(scommon.ST_STORAGE,[]byte("111"), &cstates.StorageItem{Value: []byte("11111111")})
	native.CloneCache.Add(scommon.ST_STORAGE, append(utils.MixTestContractAddress[:]), item)

	log.Infof("mix test contract address:%s",utils.MixTestContractAddress.ToBase58())
	mbc.ToTransfer(native, mbc.GenBalanceKey(utils.MbcContractAddress, utils.MixTestContractAddress),1000)

	mbc.FromTransfer(native, mbc.GenBalanceKey(utils.MbcContractAddress, utils.MixTestContractAddress),100)

	return utils.BYTE_TRUE, nil
}

func MixTestSetKey(native *native.NativeService) ([]byte, error) {

	setkeys := new(SetKeys)
	if err := setkeys.Deserialize(bytes.NewBuffer(native.Input)); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[MixTest] setkeys deserialize error!")
	}
	contract := native.ContextRef.CurrentContext().ContractAddress
	for _, v := range setkeys.States {
		fmt.Printf("MixTestSetKey %v\n",[]byte(v.Key))
		native.CloneCache.Add(scommon.ST_STORAGE,GenMixTestKey2(contract,v.From.ToBase58() + v.Key), &cstates.StorageItem{Value: []byte(v.Value)})
	}

	mbc.ToTransfer(native, mbc.GenBalanceKey(utils.MbcContractAddress, utils.MixTestContractAddress),1000)

	return utils.BYTE_TRUE, nil
}

func MixTestName(native *native.NativeService) ([]byte, error) {
	return []byte(constants.MIXT_NAME), nil
}

func MixTestSymbol(native *native.NativeService) ([]byte, error) {
	return []byte(constants.MIXT_SYMBOL), nil
}

func GenMixTestKey(contract common.Address,key []byte) []byte {
	return append(contract[:],key...)
}

func GenMixTestKey2(contract common.Address,key string) []byte {
	return append(contract[:],key...)
}

func MixTestGetKey(native *native.NativeService) ([]byte, error) {

	args := bytes.NewBuffer(native.Input)
	log.Debugf("mixTest args:%s",args)

	arg0, err := serialization.ReadVarBytes(args)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewErr("add key failed: argument 0 error, " + err.Error())
	}

	contract := native.ContextRef.CurrentContext().ContractAddress
	fmt.Printf("mix contact address %v\n",contract.ToBase58())

	state,err := utils.GetStorageItem(native,GenMixTestKey(contract,arg0))
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[MixTestGetKey] get from key error!")
	}
	if state == nil {
		return []byte("not exist"), nil
	}
	return state.Value,nil
}
