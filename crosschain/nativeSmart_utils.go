package crosschain

import (
	"bytes"
	"github.com/mixbee/mixbee/common/serialization"
	"github.com/mixbee/mixbee/smartcontract/service/native/utils"
	"github.com/mixbee/mixbee/smartcontract/service/native/crossverifynode"
	"github.com/mixbee/mixbee/core/ledger"
	httpcom "github.com/mixbee/mixbee/http/base/common"
	"github.com/mixbee/mixbee/smartcontract/event"
	"github.com/mixbee/mixbee/common/log"
	"fmt"
	ntypes "github.com/mixbee/mixbee/vm/neovm/types"
	"encoding/hex"
	"encoding/json"
)

func getVerifyNodeInfoFromNative(pbk string) (*crossverifynode.CrossVerifyNodeInfo, error) {
	param := make([]interface{}, 0)
	param = append(param, pbk)
	tx, err := httpcom.NewNativeInvokeTransaction(100, 20000, utils.CrossChainVerifynodeContractAddress, byte(0), crossverifynode.QUERY_VERIFY_NODE_INFO, param)
	if err != nil {
		return nil, err
	}
	preResult, err := ledger.DefLedger.GetStore().PreExecuteContract(tx)
	if err != nil {
		return nil, err
	}
	if preResult.State == event.CONTRACT_STATE_FAIL {
		log.Errorf("getVerifyNodeInfoFromNative error %#v", preResult.Result)
		return nil, fmt.Errorf("getVerifyNodeInfoFromNative error")
	}
	if preResult.Result.(string) == "" {
		return nil, nil
	}
	data, err := hex.DecodeString(preResult.Result.(string))
	if err != nil {
		return nil, err
	}
	nodeInfo := &crossverifynode.CrossVerifyNodeInfo{}
	err = json.Unmarshal(data, nodeInfo)
	if err != nil {
		return nil, err
	}
	fmt.Printf("getVerifyNodeInfoFromNative result %#v\n", nodeInfo)
	return nodeInfo, nil
}

func getVerifyNodeInfoMinDepositMbc(pbk string) (uint64, error) {
	param := make([]interface{}, 0)
	tx, err := httpcom.NewNativeInvokeTransaction(100, 20000, utils.CrossChainVerifynodeContractAddress, byte(0), crossverifynode.MIN_DEPOSIT_MBC_FUNCTION, param)
	if err != nil {
		return 0, err
	}
	preResult, err := ledger.DefLedger.GetStore().PreExecuteContract(tx)
	if err != nil {
		return 0, err
	}
	if preResult.State == event.CONTRACT_STATE_FAIL {
		log.Errorf("getVerifyNodeInfoFromNative error %s", preResult.Result)
		return 0, fmt.Errorf("getVerifyNodeInfoFromNative error")
	}
	data, err := hex.DecodeString(preResult.Result.(string))
	return ntypes.BigIntFromBytes(data).Uint64(), nil
}

func isExsitVerifyNodeInSmartMgr(pbk string) (bool, error) {
	bf := new(bytes.Buffer)
	if err := serialization.WriteString(bf, pbk); err != nil {
		return false, err
	}
	param := make([]interface{}, 0)
	param = append(param, bf.Bytes())
	tx, err := httpcom.NewNativeInvokeTransaction(100, 20000, utils.CrossChainVerifynodeContractAddress, byte(0), crossverifynode.IS_EXSIT_VERIFY_NODE, param)
	if err != nil {
		return false, err
	}
	preResult, err := ledger.DefLedger.GetStore().PreExecuteContract(tx)
	if err != nil {
		return false, err
	}
	if preResult.State == event.CONTRACT_STATE_FAIL {
		log.Errorf("isExsitVerifyNodeInSmartMgr error %s", preResult.Result)
		return false, fmt.Errorf("isExsitVerifyNodeInSmartMgr error")
	}
	data, err := hex.DecodeString(preResult.Result.(string))
	fmt.Printf("isExsitVerifyNodeInSmartMgr result %s\n", data)
	return string(data) == "true", nil
}
