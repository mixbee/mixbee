

package handlers

import (
	"encoding/json"
	clisvrcom "github.com/mixbee/mixbee/cmd/sigsvr/common"
	"github.com/mixbee/mixbee/cmd/utils"
	"github.com/mixbee/mixbee/common"
	"testing"
)

func TestSigNeoVMInvokeTx(t *testing.T) {
	addr1 := common.Address([20]byte{1})
	address1 := addr1.ToHexString()
	invokeReq := &SigNeoVMInvokeTxReq{
		GasPrice: 0,
		GasLimit: 0,
		Address:  address1,
		Params: []interface{}{
			&utils.NeoVMInvokeParam{
				Type:  "string",
				Value: "foo",
			},
			&utils.NeoVMInvokeParam{
				Type: "array",
				Value: []interface{}{
					&utils.NeoVMInvokeParam{
						Type:  "int",
						Value: "0",
					},
					&utils.NeoVMInvokeParam{
						Type:  "bool",
						Value: "true",
					},
				},
			},
		},
	}
	data, err := json.Marshal(invokeReq)
	if err != nil {
		t.Errorf("json.Marshal SigNeoVMInvokeTxReq error:%s", err)
		return
	}
	req := &clisvrcom.CliRpcRequest{
		Qid:    "t",
		Method: "signeovminvoketx",
		Params: data,
	}
	rsp := &clisvrcom.CliRpcResponse{}
	SigNeoVMInvokeTx(req, rsp)
	if rsp.ErrorCode != 0 {
		t.Errorf("SigNeoVMInvokeTx failed. ErrorCode:%d ErrorInfo:%s", rsp.ErrorCode, rsp.ErrorInfo)
		return
	}
}
