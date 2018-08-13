package common

import (
	"encoding/hex"
	"fmt"
	"github.com/mixbee/mixbee/smartcontract/service/native/utils"
	"github.com/mixbee/mixbee/common"
	bactor "github.com/mixbee/mixbee/http/base/actor"
)

type MixTestOfRsp struct {
	Value string `json:"value"`
}

func GetKey(key string) (*MixTestOfRsp, error) {
	value, err := GetContractKey(0, utils.MixTestContractAddress,"getkey", key)
	if err != nil {
		return nil, fmt.Errorf("get key balance error:%s", err)
	}
	return &MixTestOfRsp{
		Value: fmt.Sprintf("%s", value),
	}, nil
}

func CrossChainQuery(seqId string) (*MixTestOfRsp, error) {
	value, err := GetContractKey(0, utils.CrossChainContractAddress, "crossQuery",seqId)
	if err != nil {
		return nil, fmt.Errorf("cross chain query error:%s", err)
	}
	return &MixTestOfRsp{
		Value: fmt.Sprintf("%s", value),
	}, nil
}

func CrossChainPairEvidenceQuery(seqId string) (*MixTestOfRsp, error) {
	value, err := GetContractKey(0, utils.CrossChainPairEvidenceContractAddress, "getCrossPairTx",seqId)
	if err != nil {
		return nil, fmt.Errorf("cross chain query error:%s", err)
	}
	return &MixTestOfRsp{
		Value: fmt.Sprintf("%s", value),
	}, nil
}

func CrossChainHistory(from string) (*MixTestOfRsp, error) {
	value, err := GetContractKey(0, utils.CrossChainContractAddress, "crossHistory",from)
	if err != nil {
		return nil, fmt.Errorf("cross chain history error:%s", err)
	}
	return &MixTestOfRsp{
		Value: fmt.Sprintf("%s", value),
	}, nil
}

func GetContractKey(cVersion byte, contractAddr common.Address,method,key string) (string, error) {
	tx, err := NewNativeInvokeTransaction(0, 0, contractAddr, cVersion, method, []interface{}{key[:]})
	if err != nil {
		return "", fmt.Errorf("NewNativeInvokeTransaction error:%s", err)
	}
	result, err := bactor.PreExecuteContract(tx)
	if err != nil {
		return "", fmt.Errorf("PrepareInvokeContract error:%s", err)
	}
	if result.State == 0 {
		return "", fmt.Errorf("prepare invoke failed")
	}
	data, err := hex.DecodeString(result.Result.(string))
	if err != nil {
		return "", fmt.Errorf("hex.DecodeString error:%s", err)
	}

	return string(data), nil
}
