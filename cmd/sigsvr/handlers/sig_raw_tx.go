

package handlers

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"github.com/mixbee/mixbee-crypto/keypair"
	clisvrcom "github.com/mixbee/mixbee/cmd/sigsvr/common"
	cliutil "github.com/mixbee/mixbee/cmd/utils"
	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/common/log"
	"github.com/mixbee/mixbee/core/types"
)

type SigRawTransactionReq struct {
	RawTx string `json:"raw_tx"`
}

type SigRawTransactionRsp struct {
	SignedTx string `json:"signed_tx"`
}

func SigRawTransaction(req *clisvrcom.CliRpcRequest, resp *clisvrcom.CliRpcResponse) {
	rawReq := &SigRawTransactionReq{}
	err := json.Unmarshal(req.Params, rawReq)
	if err != nil {
		resp.ErrorCode = clisvrcom.CLIERR_INVALID_PARAMS
		return
	}
	rawTxData, err := hex.DecodeString(rawReq.RawTx)
	if err != nil {
		log.Infof("Cli Qid:%s SigRawTransaction hex.DecodeString error:%s", req.Qid, err)
		resp.ErrorCode = clisvrcom.CLIERR_INVALID_PARAMS
		return
	}
	rawTx := &types.Transaction{}
	err = rawTx.Deserialize(bytes.NewBuffer(rawTxData))
	if err != nil {
		log.Infof("Cli Qid:%s SigRawTransaction tx Deserialize error:%s", req.Qid, err)
		resp.ErrorCode = clisvrcom.CLIERR_INVALID_TX
		return
	}
	signer := clisvrcom.DefAccount
	var emptyAddress = common.Address{}
	if rawTx.Payer == emptyAddress {
		rawTx.Payer = signer.Address
	}

	txHash := rawTx.Hash()
	sigData, err := cliutil.Sign(txHash.ToArray(), signer)
	if err != nil {
		log.Infof("Cli Qid:%s SigRawTransaction Sign error:%s", req.Qid, err)
		resp.ErrorCode = clisvrcom.CLIERR_INTERNAL_ERR
		return
	}
	if len(rawTx.Sigs) == 0 {
		rawTx.Sigs = make([]*types.Sig, 0)
	}
	rawTx.Sigs = append(rawTx.Sigs, &types.Sig{
		PubKeys: []keypair.PublicKey{signer.PublicKey},
		M:       1,
		SigData: [][]byte{sigData},
	})
	buf := bytes.NewBuffer(nil)
	err = rawTx.Serialize(buf)
	if err != nil {
		log.Infof("Cli Qid:%s SigRawTransaction tx Serialize error:%s", req.Qid, err)
		resp.ErrorCode = clisvrcom.CLIERR_INTERNAL_ERR
		return
	}
	resp.Result = &SigRawTransactionRsp{
		SignedTx: hex.EncodeToString(buf.Bytes()),
	}
}
