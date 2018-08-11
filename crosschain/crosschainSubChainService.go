package crosschain

import (
	"github.com/mixbee/mixbee/core/types"
	"github.com/mixbee/mixbee/common/log"
	"github.com/mixbee/mixbee/smartcontract/event"
	bcomn "github.com/mixbee/mixbee/http/base/common"
	"github.com/mixbee/mixbee/events/message"
	"github.com/mixbee/mixbee/mixbee-eventbus/actor"
	"github.com/mixbee/mixbee/events"
	"github.com/mixbee/mixbee/smartcontract/service/native/utils"
	"encoding/json"
	"github.com/mixbee/mixbee/smartcontract/service/native/crosschain"
	"github.com/mixbee/mixbee/common/config"
	"time"
)

var SubCrossChainServerInstant *SubCrossChainServer

type SubCrossChainServer struct {
	MainVerifyNodes *VerifyNodes
	SubCrossPid     *actor.PID
}

type SubChainTimeOut struct {
}

func StartSubChainServer() {
	subPID := actor.Spawn(actor.FromFunc(Receive))
	sub1 := events.NewActorSubscriber(subPID)
	sub1.Subscribe(message.TOPIC_SMART_CODE_EVENT)

	server := &SubCrossChainServer{}
	server.SubCrossPid = subPID
	server.MainVerifyNodes = NewVerifyNodes()
	SubCrossChainServerInstant = server

	go func() {
		ticker := time.NewTicker(config.DEFAULT_CROSS_CHAIN_VERIFY_PING_TIME * time.Second)
		for {
			select {
			case <-ticker.C:
				subPID.Tell(&SubChainTimeOut{})
			}
		}
	}()
}

func Receive(c actor.Context) {
	switch msg := c.Message().(type) {
	case *message.SmartCodeEventMsg:
		pushSmartCodeEvent(*msg.Event)
	case *SubChainTimeOut:
		timeout()
	default:
	}
}

func timeout() {
	addr := config.DefConfig.CrossChain.MainVerifyNode[0]
	params := []interface{}{}
	result, err := SendRpcRequestWithAddr(addr, "getAllVerifyNodeInfo", params)
	if err != nil {
		log.Errorf("SubCrossChainServerInstant getAllVerifyNodeInfo err %s", err.Error())
		return
	}
	log.Debugf("getAllVerifyNodeInfo error %s",result)
	SubCrossChainServerInstant.UpdateMainVerifyNodes(result)
}

func (s *SubCrossChainServer) UpdateMainVerifyNodes(str []byte) {
	var list []*CrossChainVerifyNode
	err := json.Unmarshal(str, &list)
	if err != nil {
		log.Errorf("cross chain SubCrossChainServerInstant json unmarshal err %s\n", err)
		return
	}

	if len(list) == 0 {
		return
	}

	s.MainVerifyNodes.Lock()
	defer s.MainVerifyNodes.Unlock()

	nodes := make(map[string]*CrossChainVerifyNode)
	for _, value := range list {
		nodes[value.PublicKey] = value
	}
	s.MainVerifyNodes.VerifyerNodes = nodes
}

func (s *SubCrossChainServer) GetVerifyNodeInfoByPublicKey(pbk string) *CrossChainVerifyNode {
	return s.MainVerifyNodes.VerifyerNodes[pbk]
}

func (s *SubCrossChainServer) IsExsitNode(pbk string) bool {
	_,ok := s.MainVerifyNodes.VerifyerNodes[pbk]
	return ok
}

func pushSmartCodeEvent(v interface{}) {

	rs, ok := v.(types.SmartCodeEvent)
	if !ok {
		log.Errorf("[PushSmartCodeEvent]", "SmartCodeEvent err")
		return
	}

	go func() {
		switch object := rs.Result.(type) {
		case *event.ExecuteNotify:
			contractAddrs, notify := bcomn.GetExecuteNotify(object)
			pushCrossChainTxToMainChain(contractAddrs, notify)
		default:
		}
	}()
}

func pushCrossChainTxToMainChain(bools map[string]bool, notify bcomn.ExecuteNotify) {
	if _, ok := bools[utils.CrossChainContractAddress.ToHexString()]; !ok {
		return
	}

	if len(notify.Notify) < 1 {
		return
	}

	stateInfos := notify.Notify[0].States.([]interface{})
	method := stateInfos[0].(string)
	infoStr := stateInfos[1].(string)
	log.Infof("subchain pushCrossChainTxToMainChain method=%s,info=%s", method, infoStr)
	if method == crosschain.CROSS_TRANSFER {
		info := crosschain.CrossChainStateResult{}
		err := json.Unmarshal([]byte(infoStr), &info)
		if err != nil {
			log.Errorf("subchain pushCrossChainTxToMainChain json unmarshal err",err)
			return
		}
		nodeInfo := SubCrossChainServerInstant.GetVerifyNodeInfoByPublicKey(info.VerifyPublicKey)
		if nodeInfo == nil {
			log.Errorf("pushCrossChainTxToMainChain no verify public %s node info",info.VerifyPublicKey)
			return
		}
		params := []interface{}{info.From, info.To, info.AValue, info.BValue, info.AChainId, info.BChainId, notify.TxHash, info.SeqId, info.Timestamp, info.Nonce, info.VerifyPublicKey}
		result, err := SendRpcRequestWithAddr(nodeInfo.Host, "pushCrossChainTxInfo", params)
		if err != nil {
			log.Errorf("pushCrossChainTxToMainChain error %s", err.Error())
		}
		log.Infof("pushCrossChainTxToMainChain result %s", result)
	}
}
