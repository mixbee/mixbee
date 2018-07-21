

package handlers

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"sort"

	"github.com/mixbee/mixbee-crypto/keypair"
	clisvrcom "github.com/mixbee/mixbee/cmd/sigsvr/common"
	cliutil "github.com/mixbee/mixbee/cmd/utils"
	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/common/constants"
	"github.com/mixbee/mixbee/common/log"
	"github.com/mixbee/mixbee/core/signature"
	"github.com/mixbee/mixbee/core/types"
)

type SigMutilRawTransactionReq struct {
	RawTx   string   `json:"raw_tx"`
	M       int      `json:"m"`
	PubKeys []string `json:"pub_keys"`
}

type SigMutilRawTransactionRsp struct {
	SignedTx string `json:"signed_tx"`
}

func SigMutilRawTransaction(req *clisvrcom.CliRpcRequest, resp *clisvrcom.CliRpcResponse) {
	rawReq := &SigMutilRawTransactionReq{}
	err := json.Unmarshal(req.Params, rawReq)
	if err != nil {
		resp.ErrorCode = clisvrcom.CLIERR_INVALID_PARAMS
		return
	}
	numkeys := len(rawReq.PubKeys)
	if rawReq.M <= 0 || numkeys < rawReq.M || numkeys <= 1 || numkeys > constants.MULTI_SIG_MAX_PUBKEY_SIZE {
		resp.ErrorCode = clisvrcom.CLIERR_INVALID_PARAMS
		return
	}
	rawTxData, err := hex.DecodeString(rawReq.RawTx)
	if err != nil {
		log.Infof("Cli Qid:%s SigMutilRawTransaction hex.DecodeString error:%s", req.Qid, err)
		resp.ErrorCode = clisvrcom.CLIERR_INVALID_PARAMS
		return
	}
	rawTx := &types.Transaction{}
	err = rawTx.Deserialize(bytes.NewBuffer(rawTxData))
	if err != nil {
		log.Infof("Cli Qid:%s SigMutilRawTransaction tx Deserialize error:%s", req.Qid, err)
		resp.ErrorCode = clisvrcom.CLIERR_INVALID_TX
		return
	}

	pubKeys := make([]keypair.PublicKey, 0, len(rawReq.PubKeys))
	for _, pkStr := range rawReq.PubKeys {
		pkData, err := hex.DecodeString(pkStr)
		if err != nil {
			log.Info("Cli Qid:%s SigMutilRawTransaction pk hex.DecodeString error:%s", req.Qid, err)
			resp.ErrorCode = clisvrcom.CLIERR_INVALID_PARAMS
			return
		}
		pk, err := keypair.DeserializePublicKey(pkData)
		if err != nil {
			log.Info("Cli Qid:%s SigMutilRawTransaction keypair.DeserializePublicKey error:%s", req.Qid, err)
			resp.ErrorCode = clisvrcom.CLIERR_INVALID_PARAMS
			return
		}
		pubKeys = append(pubKeys, pk)
	}

	var emptyAddress = common.Address{}
	if rawTx.Payer == emptyAddress {
		payer, err := types.AddressFromMultiPubKeys(pubKeys, rawReq.M)
		if err != nil {
			log.Infof("Cli Qid:%s SigMutilRawTransaction AddressFromMultiPubKeys error:%s", req.Qid, err)
			resp.ErrorCode = clisvrcom.CLIERR_INTERNAL_ERR
			return
		}
		rawTx.Payer = payer
	}
	if len(rawTx.Sigs) == 0 {
		rawTx.Sigs = make([]*types.Sig, 0)
	}

	signer := clisvrcom.DefAccount
	txHash := rawTx.Hash()
	sigData, err := cliutil.Sign(txHash.ToArray(), signer)
	if err != nil {
		log.Infof("Cli Qid:%s SigMutilRawTransaction Sign error:%s", req.Qid, err)
		resp.ErrorCode = clisvrcom.CLIERR_INTERNAL_ERR
		return
	}

	hasMutilSig := false
	for i, sigs := range rawTx.Sigs {
		if pubKeysEqual(sigs.PubKeys, pubKeys) {
			hasMutilSig = true
			if hasAlreadySig(txHash.ToArray(), signer.PublicKey, sigs.SigData) {
				break
			}
			sigs.SigData = append(sigs.SigData, sigData)
			rawTx.Sigs[i] = sigs
			break
		}
	}
	if !hasMutilSig {
		rawTx.Sigs = append(rawTx.Sigs, &types.Sig{
			PubKeys: pubKeys,
			M:       uint16(rawReq.M),
			SigData: [][]byte{sigData},
		})
	}

	buf := bytes.NewBuffer(nil)
	err = rawTx.Serialize(buf)
	if err != nil {
		log.Infof("Cli Qid:%s SigMutilRawTransaction tx Serialize error:%s", req.Qid, err)
		resp.ErrorCode = clisvrcom.CLIERR_INTERNAL_ERR
		return
	}
	resp.Result = &SigRawTransactionRsp{
		SignedTx: hex.EncodeToString(buf.Bytes()),
	}
}

func hasAlreadySig(data []byte, pk keypair.PublicKey, sigDatas [][]byte) bool {
	for _, sigData := range sigDatas {
		err := signature.Verify(pk, data, sigData)
		if err == nil {
			return true
		}
	}
	return false
}

func pubKeysEqual(pks1, pks2 []keypair.PublicKey) bool {
	if len(pks1) != len(pks2) {
		return false
	}
	size := len(pks1)
	if size == 0 {
		return true
	}
	pkstr1 := make([]string, 0, size)
	for _, pk := range pks1 {
		pkstr1 = append(pkstr1, hex.EncodeToString(keypair.SerializePublicKey(pk)))
	}
	pkstr2 := make([]string, 0, size)
	for _, pk := range pks2 {
		pkstr2 = append(pkstr2, hex.EncodeToString(keypair.SerializePublicKey(pk)))
	}
	sort.Strings(pkstr1)
	sort.Strings(pkstr2)
	for i := 0; i < size; i++ {
		if pkstr1[i] != pkstr2[i] {
			return false
		}
	}
	return true
}
