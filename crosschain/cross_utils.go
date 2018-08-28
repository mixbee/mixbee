package crosschain

import (
	"fmt"
	"github.com/mixbee/mixbee/account"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"github.com/mixbee/mixbee/core/types"
	"time"
	"github.com/mixbee/mixbee/core/payload"
	"github.com/mixbee/mixbee/smartcontract/service/native/utils"
	"github.com/mixbee/mixbee-crypto/keypair"
	httpcom "github.com/mixbee/mixbee/http/base/common"
	sig "github.com/mixbee/mixbee-crypto/signature"
	"github.com/mixbee/mixbee/common/log"
	p2ptypes "github.com/mixbee/mixbee/p2pserver/message/types"
	"github.com/mixbee/mixbee/http/base/rpc"
	"github.com/mixbee/mixbee/smartcontract/service/native/crosspairevidence"
	"github.com/mixbee/mixbee/smartcontract/service/native/crosschaintx"
	"github.com/mixbee/mixbee/smartcontract/service/native/crossverifynode"
)

const (
	VERSION_CONTRACT_CROSS_CHAIN = byte(0)
)

type CrossChainStateResult struct {
	From            string `json:"from"`
	To              string `json:"to"`
	AValue          uint64 `json:"aValue"`
	BValue          uint64 `json:"bValue"`
	AChainId        uint32 `json:"achainId"`
	BChainId        uint32 `json:"bchainId"`
	Type            uint32 `json:"type"`
	Timestamp       uint32 `json:"timestamp"` //有效时间
	SeqId           string `json:"seqId"`
	Statue          uint32 `json:"status"`
	Nonce           uint32 `json:"nonce"`           //跨链双方的匹配随机数
	VerifyPublicKey string `json:"verifyPublicKey"` //主链验证节点公钥
	Sig             string `json:"sig"`
}

func pushCrossChainResult(signer *account.Account, addr, seqId string, sig []byte) (string, error) {
	log.Infof("cross chain || releaseLockToken||pushCrossChainResult addr=%v,seqid=%v", addr, seqId)
	result, err := CrossChainReleaseAssetByMainChain(signer, addr, seqId, sig)
	if err != nil {
		return "", err
	}
	log.Infof("cross chain || pushCrossChainReleaseResult %s", result)
	return result, nil
}

func CrossChainReleaseAssetByMainChain(signer *account.Account, addr, seqId string, sign []byte) (string, error) {
	transferTx, err := BuildCrossReleaseTx(seqId, sign)
	if err != nil {
		return "", err
	}
	err = SignTransaction(signer, transferTx)
	if err != nil {
		return "", fmt.Errorf("SignTransaction error:%s", err)
	}
	txHash, err := SendRawTransactionWithAddr(transferTx, addr)
	if err != nil {
		return "", fmt.Errorf("SendTransaction error:%s", err)
	}
	return txHash, nil
}

func CrossChainVerifyNodeRegister(signer *account.Account,info *CrossChainVerifyNode,addr string) (string, error) {
	transferTx, err := BuildCrossVerifyRegisterTx(info)
	if err != nil {
		return "", err
	}
	err = SignTransaction(signer, transferTx)
	if err != nil {
		return "", fmt.Errorf("SignTransaction error:%s", err)
	}
	txHash, err := SendRawTransactionWithAddr(transferTx,addr)
	if err != nil {
		return "", fmt.Errorf("SendTransaction error:%s", err)
	}
	return txHash, nil
}

func SignTransaction(signer *account.Account, tx *types.Transaction) error {
	tx.Payer = signer.Address
	txHash := tx.Hash()
	sigData, err := Sign(txHash.ToArray(), signer)
	if err != nil {
		return fmt.Errorf("sign error:%s", err)
	}
	sigInfo := &types.Sig{
		PubKeys: []keypair.PublicKey{signer.PublicKey},
		M:       1,
		SigData: [][]byte{sigData},
	}
	tx.Sigs = []*types.Sig{sigInfo}
	return nil
}

//Sign sign return the signature to the data of private key
func Sign(data []byte, signer *account.Account) ([]byte, error) {
	s, err := sig.Sign(signer.SigScheme, signer.PrivateKey, data, nil)
	if err != nil {
		return nil, err
	}
	sigData, err := sig.Serialize(s)
	if err != nil {
		return nil, fmt.Errorf("sig.Serialize error:%s", err)
	}
	return sigData, nil
}

func BuildCrossVerifyRegisterTx(info *CrossChainVerifyNode) (*types.Transaction, error) {
	contractAddr := utils.CrossChainVerifynodeContractAddress
	version := VERSION_CONTRACT_CROSS_CHAIN
	nodeInfo := &crossverifynode.CrossVerifyNodeInfo{
		Pbk:info.PublicKey,
	}
	invokeCode, err := httpcom.BuildNativeInvokeCode(contractAddr, version, crossverifynode.REGISTER_VERIFY_NODE, []interface{}{nodeInfo})
	if err != nil {
		return nil, fmt.Errorf("build invoke code error:%s", err)
	}
	invokePayload := &payload.InvokeCode{
		Code: invokeCode,
	}
	tx := &types.Transaction{
		GasPrice: 0,
		GasLimit: 20000,
		TxType:   types.Invoke,
		Nonce:    uint64(time.Now().UnixNano()/1e6),
		Payload:  invokePayload,
		SystemTx: true,
		Sigs:     make([]*types.Sig, 0, 0),
	}
	return tx, nil
}

func BuildCrossReleaseTx(seqId string, sig []byte) (*types.Transaction, error) {
	sigHex := hex.EncodeToString(sig)
	param := seqId + ":" + sigHex
	contractAddr := utils.CrossChainContractAddress
	version := VERSION_CONTRACT_CROSS_CHAIN
	invokeCode, err := httpcom.BuildNativeInvokeCode(contractAddr, version, crosschaintx.CROSS_RELEASE, []interface{}{param})
	if err != nil {
		return nil, fmt.Errorf("build invoke code error:%s", err)
	}
	invokePayload := &payload.InvokeCode{
		Code: invokeCode,
	}
	tx := &types.Transaction{
		GasPrice: 0,
		GasLimit: 20000,
		TxType:   types.Invoke,
		Nonce:    uint64(time.Now().UnixNano()/1e6),
		Payload:  invokePayload,
		SystemTx: true,
		Sigs:     make([]*types.Sig, 0, 0),
	}
	return tx, nil
}

func BuildCrossPairEvidenceTx(param string) (*types.Transaction, error) {
	contractAddr := utils.CrossChainPairEvidenceContractAddress
	version := VERSION_CONTRACT_CROSS_CHAIN
	invokeCode, err := httpcom.BuildNativeInvokeCode(contractAddr, version, crosspairevidence.PUSH_EVIDENCE, []interface{}{param})
	if err != nil {
		return nil, fmt.Errorf("build invoke code error:%s", err)
	}
	invokePayload := &payload.InvokeCode{
		Code: invokeCode,
	}
	tx := &types.Transaction{
		GasPrice: 0,
		GasLimit: 20000,
		TxType:   types.Invoke,
		Nonce:    uint64(time.Now().UnixNano()/1e6),
		Payload:  invokePayload,
		SystemTx: true,
		Sigs:     make([]*types.Sig, 0, 0),
	}
	return tx, nil
}

func SendRawTransactionWithAddr(tx *types.Transaction, addr string) (string, error) {
	var buffer bytes.Buffer
	err := tx.Serialize(&buffer)
	if err != nil {
		return "", fmt.Errorf("Serialize error:%s", err)
	}
	txData := hex.EncodeToString(buffer.Bytes())
	data, err := SendRpcRequestWithAddr(addr, "sendrawtransaction", []interface{}{txData})
	if err != nil {
		return "", err
	}
	hexHash := ""
	err = json.Unmarshal(data, &hexHash)
	if err != nil {
		return "", fmt.Errorf("json.Unmarshal hash:%s error:%s", data, err)
	}
	return hexHash, nil
}

func checkCrossChainTxBySeqId(tx *CTXEntry) (bool, bool) {

	firstPath := CtxServer.SubNetNodesMgr.GetSubNetNode(tx.ANetWorkId)
	tx.VerifyCheckCount = tx.VerifyCheckCount + 1
	firstInfo, err := GetCrossChainTxInfoBySeqId(firstPath, tx.SeqId)
	if err != nil {
		return false, false
	}
	secondPath := CtxServer.SubNetNodesMgr.GetSubNetNode(tx.BNetWorkId)
	secondInfo, err := GetCrossChainTxInfoBySeqId(secondPath, tx.SeqId)
	if err != nil {
		return false, false
	}

	//检查cross tx 的状态和时间是否
	if firstInfo.Statue != 0 {
		return false, true
	}

	if firstInfo.Timestamp < uint32(time.Now().Unix()) {
		return false, true
	}

	if secondInfo.Statue != 0 {
		return false, true
	}

	if secondInfo.Timestamp < uint32(time.Now().Unix()) {
		return false, true
	}

	return true, false
}

func GetCrossChainTxInfoBySeqId(addr, seqId string) (*CrossChainStateResult, error) {
	log.Debugf("GetCrossChainTxInfoBySeqId addr=%v,seqId=%v", addr, seqId)
	result, err := SendRpcRequestWithAddr(addr, "crossQuery", []interface{}{seqId})
	if err != nil {
		log.Errorf("GetCrossChainTxInfoBySeqId seqId=%s error %s", seqId, err)
		return nil, err
	}
	log.Debugf("GetCrossChainTxInfoBySeqId result = %s", string(result))

	rmap, err := Json2map(result)
	if err != nil {
		log.Errorf("GetCrossChainTxInfoBySeqId seqId=%s error %s", seqId, err)
		return nil, err
	}
	infoStr := rmap["value"].(string)
	info := &CrossChainStateResult{}
	err = json.Unmarshal([]byte(infoStr), info)
	if err != nil {
		log.Errorf("GetCrossChainTxInfoBySeqId seqId=%s error %s", seqId, err)
		return nil, err
	}
	return info, nil
}

func GetTxStateByHash(addr, hash string) (uint32, error) {

	log.Debugf("GetTxStateByHash param addr=%s,hash=%s", addr, hash)
	result, err := SendRpcRequestWithAddr(addr, "getsmartcodeevent", []interface{}{hash})
	if err != nil {
		return 0, err
	}
	log.Debugf("GetTxStateByHash result %s", string(result))
	re, err := Json2map(result)
	if err != nil {
		log.Errorf("crossChainService||GetTxStateByHash err %s", result)
		return 0, err
	}
	if re["State"] == nil {
		return 0, nil
	}

	bb, ok := re["State"].(float64)
	if !ok {
		log.Errorf("crossChainService||GetTxStateByHash err %s", result)
		return 0, nil
	}
	return uint32(bb), nil
}

func Json2map(param []byte) (s map[string]interface{}, err error) {
	var result map[string]interface{}
	if err := json.Unmarshal(param, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func pushSigedCrossTx2OtherNode(pair *CTXPairEntry, s *CTXPoolServer) {
	var info *CTXEntry
	if pair.First != nil {
		info = pair.First
	} else {
		info = pair.Second
	}
	s.pairTxPending.push(info)
	//to broadcast
	if s.P2pPid == nil {
		return
	}
	s.P2pPid.Tell(&p2ptypes.CrossChainTxInfoPayload{
		From:       info.From,
		To:         info.To,
		FromValue:  info.FromValue,
		ToValue:    info.ToValue,
		TxHash:     info.TxHash,
		ANetWorkId: info.ANetWorkId,
		BNetWorkId: info.BNetWorkId,
		SeqId:      info.SeqId,
		Type:       info.Type,
		Sig:        info.Sig,
		Pubk:       info.Pubk,
		TimeStamp:  info.TimeStamp,
		Nonce:      info.Nonce,
	})
}

func pushCrossTxEvidence2SmartContract(pair *CTXPairEntry, signer *account.Account) {

	seqId := pair.First.SeqId
	by, err := json.Marshal(pair)
	if err != nil {
		log.Errorf("cross chain || pushCrossTxEvidence2SmartContract ||json marshal err", err)
		return
	}
	hexStr := hex.EncodeToString(by)
	param := seqId + ":" + hexStr

	transferTx, err := BuildCrossPairEvidenceTx(param)
	if err != nil {
		return
	}
	err = SignTransaction(signer, transferTx)
	if err != nil {
		return
	}
	var buffer bytes.Buffer
	err = transferTx.Serialize(&buffer)
	if err != nil {
		return
	}
	txData := hex.EncodeToString(buffer.Bytes())
	result := rpc.SendRawTransaction([]interface{}{txData})
	log.Infof("cross chain || pushCrossTxEvidence2SmartContract || result %+v", result)
}
