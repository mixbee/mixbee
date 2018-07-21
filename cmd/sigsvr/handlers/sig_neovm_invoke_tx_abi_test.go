
package handlers

import (
	"encoding/json"
	clisvrcom "github.com/mixbee/mixbee/cmd/sigsvr/common"
	"testing"
)

var testNeovmAbi = `{
  "hash": "e827bf96529b5780ad0702757b8bad315e2bb8ce",
  "entrypoint": "Main",
  "functions": [
    {
      "name": "Main",
      "parameters": [
        {
          "name": "operation",
          "type": "String"
        },
        {
          "name": "args",
          "type": "Array"
        }
      ],
      "returntype": "Any"
    },
    {
      "name": "Add",
      "parameters": [
        {
          "name": "a",
          "type": "Integer"
        },
        {
          "name": "b",
          "type": "Integer"
        }
      ],
      "returntype": "Integer"
    }
  ],
  "events": []
}`

func TestSigNeoVMInvokeAbiTx(t *testing.T) {
	invokeReq := &SigNeoVMInvokeTxAbiReq{
		GasPrice: 0,
		GasLimit: 0,
		Address:  "e827bf96529b5780ad0702757b8bad315e2bb8ce",
		Method:   "Add",
		Params: []string{
			"12",
			"13",
		},
		ContractAbi: []byte(testNeovmAbi),
	}
	data, err := json.Marshal(invokeReq)
	if err != nil {
		t.Errorf("json.Marshal SigNeoVMInvokeTxReq error:%s", err)
		return
	}
	req := &clisvrcom.CliRpcRequest{
		Qid:    "t",
		Method: "signeovminvokeabitx",
		Params: data,
	}
	rsp := &clisvrcom.CliRpcResponse{}
	SigNeoVMInvokeAbiTx(req, rsp)
	if rsp.ErrorCode != 0 {
		t.Errorf("SigNeoVMInvokeAbiTx failed. ErrorCode:%d ErrorInfo:%s", rsp.ErrorCode, rsp.ErrorInfo)
		return
	}
}
