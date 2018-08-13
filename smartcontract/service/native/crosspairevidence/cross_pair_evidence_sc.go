package crosspairevidence

import (
	scommon "github.com/mixbee/mixbee/core/store/common"
	"github.com/mixbee/mixbee/errors"
	cstates "github.com/mixbee/mixbee/core/states"
	"github.com/mixbee/mixbee/smartcontract/service/native"
	"github.com/mixbee/mixbee/smartcontract/service/native/utils"
	"bytes"
	"github.com/mixbee/mixbee/common/log"
	"github.com/mixbee/mixbee/common/serialization"
	"github.com/mixbee/mixbee/common"
	"strings"
	"encoding/hex"
)

const (
	PUSH_EVIDENCE = "pushEvidence"
	GET_EVIDENCE  = "getCrossPairTx"
)

func InitCrossPairEvidence() {
	native.Contracts[utils.CrossChainPairEvidenceContractAddress] = RegisterCrossPairEvidenceContract
}

func RegisterCrossPairEvidenceContract(native *native.NativeService) {
	native.Register(PUSH_EVIDENCE, PUSH_PAIR_EVIDENCE)
	native.Register(GET_EVIDENCE, GET_CROSS_PAIR_TX)
}

func PUSH_PAIR_EVIDENCE(native *native.NativeService) ([]byte, error) {

	args := bytes.NewBuffer(native.Input)
	arg0, err := serialization.ReadString(args)
	if err != nil {
		log.Warnf("pushEvidence failed: argument %s error %s ",arg0,err.Error())
		return utils.BYTE_FALSE, errors.NewErr("pushEvidence failed: argument 0 error, " + err.Error())
	}

	infoStrs := strings.Split(arg0, ";")
	if len(infoStrs) == 0 {
		log.Warnf("pushEvidence failed: argument %s error %s ",arg0,err.Error())
		return utils.BYTE_FALSE, errors.NewErr("pushEvidence failed: argument invalid")
	}

	infoMap := make(map[string][]byte)
	for _, value := range infoStrs {
		keyValue := strings.Split(value, ":")
		if len(keyValue) != 2 {
			log.Warnf("pushEvidence failed: argument %s error %s ",value,err.Error())
			return utils.BYTE_FALSE, errors.NewErr("pushEvidence failed: argument invalid")
		}
		vaBuf,err := hex.DecodeString(keyValue[1])
		if err != nil {
			log.Warnf("pushEvidence failed: argument %s error %s ",keyValue[1],err.Error())
			return utils.BYTE_FALSE, errors.NewErr("pushEvidence failed: argument invalid")
		}
		infoMap[keyValue[0]] = vaBuf
	}

	contract := native.ContextRef.CurrentContext().ContractAddress
	for k, v := range infoMap {
		native.CloneCache.Add(scommon.ST_STORAGE, GenCrossPairEvidenceKey(contract, []byte(k)), &cstates.StorageItem{Value: []byte(v)})
	}

	return utils.BYTE_TRUE, nil
}

func GenCrossPairEvidenceKey(contract common.Address, key []byte) []byte {
	return append(contract[:], key...)
}

func GET_CROSS_PAIR_TX(native *native.NativeService) ([]byte, error) {

	args := bytes.NewBuffer(native.Input)
	log.Debugf("mixTest args:%s", args)

	arg0, err := serialization.ReadVarBytes(args)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewErr("add key failed: argument 0 error, " + err.Error())
	}

	contract := native.ContextRef.CurrentContext().ContractAddress
	state, err := utils.GetStorageItem(native, GenCrossPairEvidenceKey(contract, arg0))
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[GET_CROSS_PAIR_TX] get from key error!")
	}
	if state == nil {
		return []byte("not exist"), nil
	}
	return state.Value, nil
}
