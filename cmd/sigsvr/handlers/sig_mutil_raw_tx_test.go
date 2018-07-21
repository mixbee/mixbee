
package handlers

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"github.com/mixbee/mixbee-crypto/keypair"
	"github.com/mixbee/mixbee/account"
	clisvrcom "github.com/mixbee/mixbee/cmd/sigsvr/common"
	"github.com/mixbee/mixbee/cmd/utils"
	"github.com/mixbee/mixbee/core/types"
	"testing"
)

func TestSigMutilRawTransaction(t *testing.T) {
	acc1 := account.NewAccount("")
	acc2 := account.NewAccount("")
	pubKeys := []keypair.PublicKey{acc1.PublicKey, acc2.PublicKey}
	m := 2
	fromAddr, err := types.AddressFromMultiPubKeys(pubKeys, m)
	if err != nil {
		t.Errorf("TestSigMutilRawTransaction AddressFromMultiPubKeys error:%s", err)
		return
	}
	defAcc := clisvrcom.DefAccount
	tx, err := utils.TransferTx(0, 0, "ont", fromAddr.ToBase58(), defAcc.Address.ToBase58(), 10)
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

	rawReq := &SigMutilRawTransactionReq{
		RawTx:   hex.EncodeToString(buf.Bytes()),
		M:       m,
		PubKeys: []string{hex.EncodeToString(keypair.SerializePublicKey(acc1.PublicKey)), hex.EncodeToString(keypair.SerializePublicKey(acc2.PublicKey))},
	}
	data, err := json.Marshal(rawReq)
	if err != nil {
		t.Errorf("json.Marshal SigRawTransactionReq error:%s", err)
		return
	}
	req := &clisvrcom.CliRpcRequest{
		Qid:    "t",
		Method: "sigmutilrawtx",
		Params: data,
	}
	resp := &clisvrcom.CliRpcResponse{}
	clisvrcom.DefAccount = acc1
	SigMutilRawTransaction(req, resp)
	if resp.ErrorCode != clisvrcom.CLIERR_OK {
		t.Errorf("SigMutilRawTransaction failed,ErrorCode:%d ErrorString:%s", resp.ErrorCode, resp.ErrorInfo)
		return
	}

	clisvrcom.DefAccount = acc2
	SigMutilRawTransaction(req, resp)
	if resp.ErrorCode != clisvrcom.CLIERR_OK {
		t.Errorf("SigMutilRawTransaction failed,ErrorCode:%d ErrorString:%s", resp.ErrorCode, resp.ErrorInfo)
		return
	}
}
