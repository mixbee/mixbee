

package handlers

import (
	"encoding/json"
	"github.com/mixbee/mixbee/account"
	clisvrcom "github.com/mixbee/mixbee/cmd/sigsvr/common"
	"testing"
)

func TestSigTransferTransaction(t *testing.T) {
	acc := account.NewAccount("")
	defAcc := clisvrcom.DefAccount
	sigReq := &SigTransferTransactionReq{
		GasLimit: 0,
		GasPrice: 0,
		Asset:    "mbc",
		From:     defAcc.Address.ToBase58(),
		To:       acc.Address.ToBase58(),
		Amount:   10,
	}
	data, err := json.Marshal(sigReq)
	if err != nil {
		t.Errorf("json.Marshal SigTransferTransactionReq error:%s", err)
	}
	req := &clisvrcom.CliRpcRequest{
		Qid:    "t",
		Method: "sigtransfertx",
		Params: data,
	}
	rsp := &clisvrcom.CliRpcResponse{}
	SigTransferTransaction(req, rsp)
	if rsp.ErrorCode != 0 {
		t.Errorf("SigTransferTransaction failed. ErrorCode:%d", rsp.ErrorCode)
		return
	}
}
