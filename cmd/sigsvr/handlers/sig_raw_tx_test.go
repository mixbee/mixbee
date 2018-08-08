

package handlers

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"github.com/mixbee/mixbee/account"
	clisvrcom "github.com/mixbee/mixbee/cmd/sigsvr/common"
	"github.com/mixbee/mixbee/cmd/utils"
	"github.com/mixbee/mixbee/common/log"
	"os"
	"testing"
)

var (
	wallet *account.ClientImpl
	passwd = []byte("123456")
)

func TestMain(m *testing.M) {
	log.InitLog(0, os.Stdout)
	clisvrcom.DefAccount = account.NewAccount("")
	m.Run()
	os.RemoveAll("./ActorLog")
	os.RemoveAll("./Log")
}

func TestSigRawTx(t *testing.T) {
	acc := account.NewAccount("")
	defAcc := clisvrcom.DefAccount
	tx, err := utils.TransferTx(0, 0, "mbc", defAcc.Address.ToBase58(), acc.Address.ToBase58(), 10)
	if err != nil {
		t.Errorf("TransferTx error:%s", err)
		return
	}
	buf := bytes.NewBuffer(nil)
	err = tx.Serialize(buf)
	if err != nil {
		t.Errorf("tx.Serialize error:%s", err)
		return
	}
	rawReq := &SigRawTransactionReq{
		RawTx: hex.EncodeToString(buf.Bytes()),
	}
	data, err := json.Marshal(rawReq)
	if err != nil {
		t.Errorf("json.Marshal SigRawTransactionReq error:%s", err)
		return
	}
	req := &clisvrcom.CliRpcRequest{
		Qid:    "t",
		Method: "sigrawtx",
		Params: data,
	}
	resp := &clisvrcom.CliRpcResponse{}
	SigRawTransaction(req, resp)
	if resp.ErrorCode != 0 {
		t.Errorf("SigRawTransaction failed. ErrorCode:%d", resp.ErrorCode)
		return
	}
}
