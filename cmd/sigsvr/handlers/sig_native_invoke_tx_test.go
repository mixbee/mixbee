

package handlers

import (
	"encoding/json"
	"github.com/mixbee/mixbee/cmd/abi"
	clisvrcom "github.com/mixbee/mixbee/cmd/sigsvr/common"
	"github.com/mixbee/mixbee/common"
	nutils "github.com/mixbee/mixbee/smartcontract/service/native/utils"
	"testing"
)

func TestSigNativeInvokeTx(t *testing.T) {
	addr1 := common.Address([20]byte{1})
	address1 := addr1.ToBase58()
	addr2 := common.Address([20]byte{2})
	address2 := addr2.ToBase58()
	invokeReq := &SigNativeInvokeTxReq{
		GasPrice: 0,
		GasLimit: 40000,
		Address:  nutils.MbcContractAddress.ToHexString(),
		Method:   "transfer",
		Version:  0,
		Params: []interface{}{
			[]interface{}{
				[]interface{}{
					address1,
					address2,
					"10000000000",
				},
			},
		},
	}
	data, err := json.Marshal(invokeReq)
	if err != nil {
		t.Errorf("json.Marshal SigNativeInvokeTxReq error:%s", err)
		return
	}
	req := &clisvrcom.CliRpcRequest{
		Qid:    "t",
		Method: "signativeinvoketx",
		Params: data,
	}
	rsp := &clisvrcom.CliRpcResponse{}
	abiPath := "../../abi"
	abi.DefAbiMgr.Init(abiPath)
	SigNativeInvokeTx(req, rsp)
	if rsp.ErrorCode != 0 {
		t.Errorf("SigNativeInvokeTx failed. ErrorCode:%d ErrorInfo:%s", rsp.ErrorCode, rsp.ErrorInfo)
		return
	}
}
