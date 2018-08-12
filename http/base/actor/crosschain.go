package actor

import (
	"github.com/mixbee/mixbee/mixbee-eventbus/actor"
	"time"
	"github.com/mixbee/mixbee/common/log"
	"fmt"
)

var CrosschainSrvPid *actor.PID

func SetCrossChainPid(actr *actor.PID) {
	CrosschainSrvPid = actr
}

type GetAllVerifyNodeReq struct {
}

type CrossSubNetNodeReq struct {
	Host  string `json:"host"`
	NetId uint32 `json:"netId"`
}

type CheckSubNetId struct {
	NetIds []uint32
}

func CheckSubNetIdFunction(ids []uint32) bool {

	if CrosschainSrvPid == nil {
		log.Errorf("Crosschain service pid actor not init")
		return false
	}

	if len(ids) == 0 {
		return false
	}

	nets := &CheckSubNetId{
		NetIds: ids,
	}
	future := CrosschainSrvPid.RequestFuture(nets, REQ_TIMEOUT*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Errorf(ERR_ACTOR_COMM, err)
		return false
	}
	return result.(bool)
}

func GetAllCrossChainVerifyNodeInfo() (interface{}, error) {

	if CrosschainSrvPid == nil {
		return nil, fmt.Errorf("Crosschain service pid actor not init")
	}
	future := CrosschainSrvPid.RequestFuture(&GetAllVerifyNodeReq{}, REQ_TIMEOUT*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Errorf(ERR_ACTOR_COMM, err)
		return nil, err
	}
	return result, nil
}

type PushCrossChainTxRsq struct {
	From       string
	To         string
	FromValue  uint64
	ToValue    uint64
	TxHash     string
	ANetWorkId uint32
	BNetWorkId uint32
	State      uint32 //0 待确认  1 打包成功
	SeqId      string
	Type       uint   //跨链资产类型
	Sig        []byte //验证节点对结果的签名
	Pubk       string //验证节点公钥
	TimeStamp  uint32 //过期时间
	Nonce      uint32 //交易双方的nonce值,必须一样
}

func RegisterSubNetNode(netId uint32, host string) {
	CrosschainSrvPid.Tell(&CrossSubNetNodeReq{
		NetId: netId,
		Host:  host,
	})
}

func PushCrossChainTx(params []interface{}) error {
	if CrosschainSrvPid == nil {
		return fmt.Errorf("Crosschain service pid actor not init")
	}

	from, ok := params[0].(string)
	if !ok {
		log.Errorf("PushCrossChainTx || param-0 invalid %v", params[0])
		return fmt.Errorf("param-0 invalid %v", params[0])
	}
	to, ok := params[1].(string)
	if !ok {
		log.Errorf("PushCrossChainTx ||param-1 invalid %v", params[1])
		return fmt.Errorf("param-1 invalid %v", params[1])
	}
	fromV, ok := params[2].(float64)
	if !ok {
		log.Errorf("PushCrossChainTx || param-2 invalid %v", params[2])
		return fmt.Errorf("param-2 invalid %v", params[2])
	}
	toV, ok := params[3].(float64)
	if !ok {
		log.Errorf("PushCrossChainTx ||param-3 invalid %v", params[2])
		return fmt.Errorf("param-3 invalid %v", params[3])
	}
	aNId, ok := params[4].(float64)
	if !ok {
		return fmt.Errorf("PushCrossChainTx || param-4 invalid %v", params[4])
	}
	bNId, ok := params[5].(float64)
	if !ok {
		log.Infof("PushCrossChainTx || param-5 invalid %v", params[5])
		return fmt.Errorf("param-5 invalid %v", params[5])
	}
	txHash, ok := params[6].(string)
	if !ok {
		log.Infof("PushCrossChainTx || param-6 invalid %v", params[6])
		return fmt.Errorf("param-6 invalid %v", params[6])
	}
	seqId, ok := params[7].(string)
	if !ok {
		log.Infof("PushCrossChainTx || param-7 invalid %v", params[7])
		return fmt.Errorf("param-7 invalid %v", params[7])
	}
	timestamp, ok := params[8].(float64)
	if !ok {
		log.Infof("PushCrossChainTx || param-8 invalid %v", params[8])
		return fmt.Errorf("param-8 invalid %v", params[8])
	}
	nonce, ok := params[9].(float64)
	if !ok {
		log.Infof("PushCrossChainTx || param-9 invalid %v", params[9])
		return fmt.Errorf("param-9 invalid %v", params[9])
	}
	publickKey, ok := params[10].(string)
	if !ok {
		log.Infof("PushCrossChainTx || param-10 invalid %v", params[10])
		return fmt.Errorf("param-7 invalid %v", params[10])
	}

	ok = CheckSubNetIdFunction([]uint32{uint32(aNId), uint32(bNId)})
	if !ok {
		log.Warnf("PushCrossChainTx || nid=%d or %d invalid", uint32(aNId), uint32(bNId))
		return fmt.Errorf("invalid networkId %v or %v,subchain not register in mainchain", aNId, bNId)
	}

	entry := &PushCrossChainTxRsq{
		From:       from,
		To:         to,
		FromValue:  uint64(fromV),
		ToValue:    uint64(toV),
		ANetWorkId: uint32(aNId),
		BNetWorkId: uint32(bNId),
		TxHash:     txHash,
		SeqId:      seqId,
		TimeStamp:  uint32(timestamp),
		Nonce:      uint32(nonce),
		Pubk:       publickKey,
	}

	log.Infof("PushCrossChainTx info = %#v", entry)
	CrosschainSrvPid.Tell(entry)

	return nil
}
