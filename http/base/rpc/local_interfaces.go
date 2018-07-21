

package rpc

import (
	"os"
	"path/filepath"

	"github.com/mixbee/mixbee/common/log"
	bactor "github.com/mixbee/mixbee/http/base/actor"
	"github.com/mixbee/mixbee/http/base/common"
	berr "github.com/mixbee/mixbee/http/base/error"
)

const (
	RANDBYTELEN = 4
)

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func GetNeighbor(params []interface{}) map[string]interface{} {
	addr := bactor.GetNeighborAddrs()
	return responseSuccess(addr)
}

func GetNodeState(params []interface{}) map[string]interface{} {
	state, err := bactor.GetConnectionState()
	if err != nil {
		return responsePack(berr.INTERNAL_ERROR, false)
	}
	t, err := bactor.GetNodeTime()
	if err != nil {
		return responsePack(berr.INTERNAL_ERROR, false)
	}
	port, err := bactor.GetNodePort()
	if err != nil {
		return responsePack(berr.INTERNAL_ERROR, false)
	}
	id, err := bactor.GetID()
	if err != nil {
		return responsePack(berr.INTERNAL_ERROR, false)
	}
	ver, err := bactor.GetVersion()
	if err != nil {
		return responsePack(berr.INTERNAL_ERROR, false)
	}
	tpe, err := bactor.GetNodeType()
	if err != nil {
		return responsePack(berr.INTERNAL_ERROR, false)
	}
	relay, err := bactor.GetRelayState()
	if err != nil {
		return responsePack(berr.INTERNAL_ERROR, false)
	}
	height := bactor.GetCurrentBlockHeight()
	txnCnt, err := bactor.GetTxnCount()
	if err != nil {
		return responsePack(berr.INTERNAL_ERROR, false)
	}
	n := common.NodeInfo{
		NodeState:   uint(state),
		NodeTime:    t,
		NodePort:    port,
		ID:          id,
		NodeVersion: ver,
		NodeType:    tpe,
		Relay:       relay,
		Height:      height,
		TxnCnt:      txnCnt,
	}
	return responseSuccess(n)
}

func StartConsensus(params []interface{}) map[string]interface{} {
	if err := bactor.ConsensusSrvStart(); err != nil {
		return responsePack(berr.INTERNAL_ERROR, false)
	}
	return responsePack(berr.SUCCESS, true)
}

func StopConsensus(params []interface{}) map[string]interface{} {
	if err := bactor.ConsensusSrvHalt(); err != nil {
		return responsePack(berr.INTERNAL_ERROR, false)
	}
	return responsePack(berr.SUCCESS, true)
}

func SetDebugInfo(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return responsePack(berr.INVALID_PARAMS, "")
	}
	switch params[0].(type) {
	case float64:
		level := params[0].(float64)
		if err := log.Log.SetDebugLevel(int(level)); err != nil {
			return responsePack(berr.INVALID_PARAMS, "")
		}
	default:
		return responsePack(berr.INVALID_PARAMS, "")
	}
	return responsePack(berr.SUCCESS, true)
}
