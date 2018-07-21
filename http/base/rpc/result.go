

package rpc

import (
	Err "github.com/mixbee/mixbee/http/base/error"
)

func responseSuccess(result interface{}) map[string]interface{} {
	return responsePack(Err.SUCCESS, result)
}
func responsePack(errcode int64, result interface{}) map[string]interface{} {
	resp := map[string]interface{}{
		"error":  errcode,
		"desc":   Err.ErrMap[errcode],
		"result": result,
	}
	return resp
}
