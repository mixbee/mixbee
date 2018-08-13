package rpc

import (
	berr "github.com/mixbee/mixbee/http/base/error"
	bcomn "github.com/mixbee/mixbee/http/base/common"
)
func GetKey(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return responsePack(berr.INVALID_PARAMS, "")
	}
	key, ok := params[0].(string)
	if !ok {
		return responsePack(berr.INVALID_PARAMS, "")
	}
	rsp, err := bcomn.GetKey(key)
	if err != nil {
		return responsePack(berr.INVALID_PARAMS, "")
	}
	return responseSuccess(rsp)
}

func CrossChainQuery(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return responsePack(berr.INVALID_PARAMS, "")
	}
	key, ok := params[0].(string)
	if !ok {
		return responsePack(berr.INVALID_PARAMS, "")
	}
	rsp, err := bcomn.CrossChainQuery(key)
	if err != nil {
		return responsePack(berr.INVALID_PARAMS, "")
	}
	return responseSuccess(rsp)
}

func CrossChainPairEvidenceQuery(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return responsePack(berr.INVALID_PARAMS, "")
	}
	key, ok := params[0].(string)
	if !ok {
		return responsePack(berr.INVALID_PARAMS, "")
	}
	rsp, err := bcomn.CrossChainPairEvidenceQuery(key)
	if err != nil {
		return responsePack(berr.INVALID_PARAMS, "")
	}
	return responseSuccess(rsp)
}

func CrossChainHistory(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return responsePack(berr.INVALID_PARAMS, "")
	}
	key, ok := params[0].(string)
	if !ok {
		return responsePack(berr.INVALID_PARAMS, "")
	}
	rsp, err := bcomn.CrossChainHistory(key)
	if err != nil {
		return responsePack(berr.INVALID_PARAMS, "")
	}
	return responseSuccess(rsp)
}
