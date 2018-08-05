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
	cmdutils "github.com/mixbee/mixbee/cmd/utils"
	"github.com/mixbee/mixbee/common/config"
	"fmt"
)

func StartSubChainServer() {
	subPID := actor.Spawn(actor.FromFunc(Receive))
	sub1 := events.NewActorSubscriber(subPID)
	sub1.Subscribe(message.TOPIC_SMART_CODE_EVENT)
}

func Receive(c actor.Context) {
	switch msg := c.Message().(type) {
	case *message.SmartCodeEventMsg:
		pushSmartCodeEvent(*msg.Event)
	default:
	}
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
	if _,ok := bools[utils.CrossChainContractAddress.ToHexString()];!ok {
		return
	}

	if len(notify.Notify) < 1 {
		return
	}

	stateInfos := notify.Notify[0].States.([]interface{})
	method := stateInfos[0].(string)
	infoStr := stateInfos[1].(string)
	log.Infof("subchain pushCrossChainTxToMainChain method=%s,info=%s",method,infoStr)
	if method == crosschain.CROSS_TRANSFER {
		info := crosschain.CrossChainStateResult{}
		json.Unmarshal([]byte(infoStr),&info)
		addr := config.DefConfig.CrossChain.MainVerifyNode[0]
		//["addrA","addrB","aAmount","bAmount","aNetId","bNetId","txHash","seqId","timestamp"]
		params := []interface{}{info.From,info.To,info.AValue,info.BValue,info.AChainId,info.BChainId,notify.TxHash,info.SeqId,info.Timestamp,info.Nonce}
		result,err := cmdutils.SendRpcRequestWithAddr(addr,"pushCrossChainTxInfo",params)
		if err != nil {
			log.Errorf("pushCrossChainTxToMainChain error %s",err.Error())
		}
		fmt.Printf("pushCrossChainTxToMainChain result %s\n",result)
	}
}
