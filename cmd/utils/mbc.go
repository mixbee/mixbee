package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/mixbee/mixbee-crypto/keypair"
	sig "github.com/mixbee/mixbee-crypto/signature"
	"github.com/mixbee/mixbee/account"
	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/common/serialization"
	"github.com/mixbee/mixbee/core/payload"
	"github.com/mixbee/mixbee/core/types"
	httpcom "github.com/mixbee/mixbee/http/base/common"
	rpccommon "github.com/mixbee/mixbee/http/base/common"
	"github.com/mixbee/mixbee/smartcontract/service/native/mbc"
	"github.com/mixbee/mixbee/smartcontract/service/native/utils"
	"github.com/mixbee/mixbee/smartcontract/service/wasmvm"
	cstates "github.com/mixbee/mixbee/smartcontract/states"
	"github.com/mixbee/mixbee/vm/wasmvm/exec"
	"strconv"
	"strings"
	"time"
	"github.com/mixbee/mixbee/smartcontract/service/native/mixtest"
	"github.com/mixbee/mixbee/smartcontract/service/native/crosschain"
)

const (
	VERSION_TRANSACTION          = byte(0)
	VERSION_CONTRACT_MBC         = byte(0)
	VERSION_CONTRACT_MBG         = byte(0)
	VERSION_CONTRACT_MIXT        = byte(0)
	VERSION_CONTRACT_CROSS_CHAIN = byte(0)
	CONTRACT_TRANSFER            = "transfer"
	CONTRACT_TRANSFER_FROM       = "transferFrom"
	CONTRACT_APPROVE             = "approve"

	CONTRACT_SETKEY = "setkey"

	CONTRACT_CROSS_TRANSFER = "crossTranfer"
	CONTRACT_CROSS_UNLOCK   = "crossUnlock"
	CONTRACT_CROSS_RELEASE  = "crossRelease"

	ASSET_MBC = "mbc"
	ASSET_MBG = "mbg"
)

//Return balance of address in base58 code
func GetBalance(address string) (*httpcom.BalanceOfRsp, error) {
	result, err := sendRpcRequest("getbalance", []interface{}{address})
	if err != nil {
		return nil, fmt.Errorf("sendRpcRequest error:%s", err)
	}
	balance := &httpcom.BalanceOfRsp{}
	err = json.Unmarshal(result, balance)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal error:%s", err)
	}
	return balance, nil
}

func CrossQuery(seqId string) (*httpcom.MixTestOfRsp, error) {
	result, err := sendRpcRequest("crossQuery", []interface{}{seqId})
	if err != nil {
		return nil, fmt.Errorf("sendRpcRequest error:%s", err)
	}

	rsp := &httpcom.MixTestOfRsp{}
	err = json.Unmarshal(result, rsp)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal error:%s", err)
	}
	return rsp, nil
}

func CrossPairEvidenceQuery(seqId string) (*httpcom.MixTestOfRsp, error) {
	result, err := sendRpcRequest("crossPairEvidenceQuery", []interface{}{seqId})
	if err != nil {
		return nil, fmt.Errorf("sendRpcRequest error:%s", err)
	}

	rsp := &httpcom.MixTestOfRsp{}
	err = json.Unmarshal(result, rsp)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal error:%s", err)
	}
	return rsp, nil
}

func CrossHistory(from string) (*httpcom.MixTestOfRsp, error) {
	result, err := sendRpcRequest("crossHistory", []interface{}{from})
	if err != nil {
		return nil, fmt.Errorf("sendRpcRequest error:%s", err)
	}

	rsp := &httpcom.MixTestOfRsp{}
	err = json.Unmarshal(result, rsp)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal error:%s", err)
	}
	return rsp, nil
}

func GetKey(key string) (*httpcom.MixTestOfRsp, error) {
	result, err := sendRpcRequest("getkey", []interface{}{key})
	if err != nil {
		return nil, fmt.Errorf("sendRpcRequest error:%s", err)
	}

	rsp := &httpcom.MixTestOfRsp{}
	err = json.Unmarshal(result, rsp)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal error:%s", err)
	}
	return rsp, nil
}

func SetKey(gasPrice, gasLimit uint64, signer *account.Account, key, value string) (string, error) {
	transferTx, err := SetKeyTx(gasPrice, gasLimit, signer.Address.ToBase58(), key, value)
	if err != nil {
		return "", err
	}
	err = SignTransaction(signer, transferTx)
	if err != nil {
		return "", fmt.Errorf("SignTransaction error:%s", err)
	}
	txHash, err := SendRawTransaction(transferTx)
	if err != nil {
		return "", fmt.Errorf("SendTransaction error:%s", err)
	}
	return txHash, nil
}

func CrossUnlockOrRelease(gasPrice, gasLimit uint64, signer *account.Account, seqId string, method string) (string, error) {
	transferTx, err := CrossUnlockTx(gasPrice, gasLimit, seqId, method)
	if err != nil {
		return "", err
	}
	err = SignTransaction(signer, transferTx)
	if err != nil {
		return "", fmt.Errorf("SignTransaction error:%s", err)
	}
	txHash, err := SendRawTransaction(transferTx)
	if err != nil {
		return "", fmt.Errorf("SendTransaction error:%s", err)
	}
	return txHash, nil
}

func GetAllowance(asset, from, to string) (string, error) {
	result, err := sendRpcRequest("getallowance", []interface{}{asset, from, to})
	if err != nil {
		return "", fmt.Errorf("sendRpcRequest error:%s", err)
	}
	balance := ""
	err = json.Unmarshal(result, &balance)
	if err != nil {
		return "", fmt.Errorf("json.Unmarshal error:%s", err)
	}
	return balance, nil
}

//tranfer asset cross chain
func CrossTransfer(gasPrice, gasLimit uint64, signer *account.Account, asset, to string, aAmount, bAmount, aChainId, bChainID, delayTime, nonce uint64, verifyPublicKey string) (string, string, error) {
	transferTx, seqId, err := CrossChainTransferTx(gasPrice, gasLimit, asset, signer.Address.ToBase58(), to, aAmount, bAmount, aChainId, bChainID, delayTime, nonce, verifyPublicKey)
	if err != nil {
		return "", "", err
	}
	err = SignTransaction(signer, transferTx)
	if err != nil {
		return "", "", fmt.Errorf("SignTransaction error:%s", err)
	}
	txHash, err := SendRawTransaction(transferTx)
	if err != nil {
		return "", "", fmt.Errorf("SendTransaction error:%s", err)
	}
	return txHash, seqId, nil
}

//Transfer mbc|mbg from account to another account
func Transfer(gasPrice, gasLimit uint64, signer *account.Account, asset, from, to string, amount uint64) (string, error) {
	transferTx, err := TransferTx(gasPrice, gasLimit, asset, signer.Address.ToBase58(), to, amount)
	if err != nil {
		return "", err
	}
	err = SignTransaction(signer, transferTx)
	if err != nil {
		return "", fmt.Errorf("SignTransaction error:%s", err)
	}
	txHash, err := SendRawTransaction(transferTx)
	if err != nil {
		return "", fmt.Errorf("SendTransaction error:%s", err)
	}
	return txHash, nil
}

func TransferFrom(gasPrice, gasLimit uint64, signer *account.Account, asset, sender, from, to string, amount uint64) (string, error) {
	transferFromTx, err := TransferFromTx(gasPrice, gasLimit, asset, sender, from, to, amount)
	if err != nil {
		return "", err
	}
	err = SignTransaction(signer, transferFromTx)
	if err != nil {
		return "", fmt.Errorf("SignTransaction error:%s", err)
	}
	txHash, err := SendRawTransaction(transferFromTx)
	if err != nil {
		return "", fmt.Errorf("SendTransaction error:%s", err)
	}
	return txHash, nil
}

func Approve(gasPrice, gasLimit uint64, signer *account.Account, asset, from, to string, amount uint64) (string, error) {
	approveTx, err := ApproveTx(gasPrice, gasLimit, asset, from, to, amount)
	if err != nil {
		return "", err
	}
	err = SignTransaction(signer, approveTx)
	if err != nil {
		return "", fmt.Errorf("SignTransaction error:%s", err)
	}
	txHash, err := SendRawTransaction(approveTx)
	if err != nil {
		return "", fmt.Errorf("SendTransaction error:%s", err)
	}
	return txHash, nil
}

func ApproveTx(gasPrice, gasLimit uint64, asset string, from, to string, amount uint64) (*types.Transaction, error) {
	fromAddr, err := common.AddressFromBase58(from)
	if err != nil {
		return nil, fmt.Errorf("from address:%s invalid:%s", from, err)
	}
	toAddr, err := common.AddressFromBase58(to)
	if err != nil {
		return nil, fmt.Errorf("To address:%s invalid:%s", to, err)
	}
	var state = &mbc.State{
		From:  fromAddr,
		To:    toAddr,
		Value: amount,
	}
	var version byte
	var contractAddr common.Address
	switch strings.ToLower(asset) {
	case ASSET_MBC:
		version = VERSION_CONTRACT_MBC
		contractAddr = utils.MbcContractAddress
	case ASSET_MBG:
		version = VERSION_CONTRACT_MBG
		contractAddr = utils.MbgContractAddress
	default:
		return nil, fmt.Errorf("Unsupport asset:%s", asset)
	}
	invokeCode, err := httpcom.BuildNativeInvokeCode(contractAddr, version, CONTRACT_APPROVE, []interface{}{state})
	if err != nil {
		return nil, fmt.Errorf("build invoke code error:%s", err)
	}
	invokePayload := &payload.InvokeCode{
		Code: invokeCode,
	}
	tx := &types.Transaction{
		GasPrice: gasPrice,
		GasLimit: gasLimit,
		TxType:   types.Invoke,
		Nonce:    uint64(time.Now().UnixNano()/1e6),
		Payload:  invokePayload,
		Sigs:     make([]*types.Sig, 0, 0),
	}
	return tx, nil
}

func SetKeyTx(gasPrice, gasLimit uint64, from, key, value string) (*types.Transaction, error) {
	fromAddr, err := common.AddressFromBase58(from)
	if err != nil {
		return nil, fmt.Errorf("from address:%s invalid:%s", from, err)
	}

	var sts []*mixtest.State
	sts = append(sts, &mixtest.State{
		From:  fromAddr,
		Key:   key,
		Value: value,
	})
	contractAddr := utils.MixTestContractAddress
	version := VERSION_CONTRACT_MIXT
	invokeCode, err := httpcom.BuildNativeInvokeCode(contractAddr, version, CONTRACT_SETKEY, []interface{}{sts})
	if err != nil {
		return nil, fmt.Errorf("build invoke code error:%s", err)
	}
	invokePayload := &payload.InvokeCode{
		Code: invokeCode,
	}
	tx := &types.Transaction{
		GasPrice: gasPrice,
		GasLimit: gasLimit,
		TxType:   types.Invoke,
		Nonce:    uint64(time.Now().UnixNano()/1e6),
		Payload:  invokePayload,
		Sigs:     make([]*types.Sig, 0, 0),
	}
	return tx, nil
}

func CrossUnlockTx(gasPrice, gasLimit uint64, seqId string, method string) (*types.Transaction, error) {

	contractAddr := utils.CrossChainContractAddress
	version := VERSION_CONTRACT_CROSS_CHAIN
	invokeCode, err := httpcom.BuildNativeInvokeCode(contractAddr, version, method, []interface{}{seqId})
	if err != nil {
		return nil, fmt.Errorf("build invoke code error:%s", err)
	}
	invokePayload := &payload.InvokeCode{
		Code: invokeCode,
	}
	tx := &types.Transaction{
		GasPrice: gasPrice,
		GasLimit: gasLimit,
		TxType:   types.Invoke,
		Nonce:    uint64(time.Now().UnixNano()/1e6),
		Payload:  invokePayload,
		SystemTx:true,
		Sigs:     make([]*types.Sig, 0, 0),
	}
	return tx, nil
}

func CrossChainTransferTx(gasPrice, gasLimit uint64, asset, from, to string, aAmount, bAmount, aChainId, bChainId, delayTime, nonce uint64, verifyPublicKey string) (*types.Transaction, string, error) {
	fromAddr, err := common.AddressFromBase58(from)
	if err != nil {
		return nil, "", fmt.Errorf("from address:%s invalid:%s", from, err)
	}
	toAddr, err := common.AddressFromBase58(to)
	if err != nil {
		return nil, "", fmt.Errorf("To address:%s invalid:%s", to, err)
	}

	crossState := &crosschain.CrossChainState{
		From:            fromAddr,
		To:              toAddr,
		AValue:          aAmount,
		BValue:          bAmount,
		AChainId:        uint32(aChainId),
		BChainId:        uint32(bChainId),
		Type:            0,
		Timestamp:       uint32(time.Now().Unix() + int64(delayTime)),
		Nonce:           uint32(nonce),
		VerifyPublicKey: verifyPublicKey,
	}
	crossState.SeqId = crosschain.GetSeqId(crossState)

	invokeCode, err := httpcom.BuildNativeInvokeCode(utils.CrossChainContractAddress, VERSION_CONTRACT_CROSS_CHAIN, CONTRACT_CROSS_TRANSFER, []interface{}{crossState})
	if err != nil {
		return nil, "", fmt.Errorf("build invoke code error:%s", err)
	}
	invokePayload := &payload.InvokeCode{
		Code: invokeCode,
	}
	tx := &types.Transaction{
		GasPrice: gasPrice,
		GasLimit: gasLimit,
		TxType:   types.Invoke,
		Nonce:    uint64(time.Now().UnixNano()/1e6),
		Payload:  invokePayload,
		Sigs:     make([]*types.Sig, 0, 0),
	}
	return tx, crossState.SeqId, nil
}

func TransferTx(gasPrice, gasLimit uint64, asset, from, to string, amount uint64) (*types.Transaction, error) {
	fromAddr, err := common.AddressFromBase58(from)
	if err != nil {
		return nil, fmt.Errorf("from address:%s invalid:%s", from, err)
	}
	toAddr, err := common.AddressFromBase58(to)
	if err != nil {
		return nil, fmt.Errorf("To address:%s invalid:%s", to, err)
	}
	var sts []*mbc.State
	sts = append(sts, &mbc.State{
		From:  fromAddr,
		To:    toAddr,
		Value: amount,
	})
	var version byte
	var contractAddr common.Address
	switch strings.ToLower(asset) {
	case ASSET_MBC:
		version = VERSION_CONTRACT_MBC
		contractAddr = utils.MbcContractAddress
	case ASSET_MBG:
		version = VERSION_CONTRACT_MBG
		contractAddr = utils.MbgContractAddress
	default:
		return nil, fmt.Errorf("Unsupport asset:%s", asset)
	}
	invokeCode, err := httpcom.BuildNativeInvokeCode(contractAddr, version, CONTRACT_TRANSFER, []interface{}{sts})
	if err != nil {
		return nil, fmt.Errorf("build invoke code error:%s", err)
	}
	invokePayload := &payload.InvokeCode{
		Code: invokeCode,
	}
	tx := &types.Transaction{
		GasPrice: gasPrice,
		GasLimit: gasLimit,
		TxType:   types.Invoke,
		Nonce:    uint64(time.Now().UnixNano()/1e6),
		Payload:  invokePayload,
		Sigs:     make([]*types.Sig, 0, 0),
	}
	return tx, nil
}

func TransferFromTx(gasPrice, gasLimit uint64, asset, sender, from, to string, amount uint64) (*types.Transaction, error) {
	senderAddr, err := common.AddressFromBase58(sender)
	if err != nil {
		return nil, fmt.Errorf("sender address:%s invalid:%s", to, err)
	}
	fromAddr, err := common.AddressFromBase58(from)
	if err != nil {
		return nil, fmt.Errorf("from address:%s invalid:%s", from, err)
	}
	toAddr, err := common.AddressFromBase58(to)
	if err != nil {
		return nil, fmt.Errorf("To address:%s invalid:%s", to, err)
	}
	transferFrom := &mbc.TransferFrom{
		Sender: senderAddr,
		From:   fromAddr,
		To:     toAddr,
		Value:  amount,
	}
	var version byte
	var contractAddr common.Address
	switch strings.ToLower(asset) {
	case ASSET_MBC:
		version = VERSION_CONTRACT_MBC
		contractAddr = utils.MbcContractAddress
	case ASSET_MBG:
		version = VERSION_CONTRACT_MBG
		contractAddr = utils.MbgContractAddress
	default:
		return nil, fmt.Errorf("Unsupport asset:%s", asset)
	}
	invokeCode, err := httpcom.BuildNativeInvokeCode(contractAddr, version, CONTRACT_TRANSFER_FROM, []interface{}{transferFrom})
	if err != nil {
		return nil, fmt.Errorf("build invoke code error:%s", err)
	}
	invokePayload := &payload.InvokeCode{
		Code: invokeCode,
	}
	tx := &types.Transaction{
		GasPrice: gasPrice,
		GasLimit: gasLimit,
		TxType:   types.Invoke,
		Nonce:    uint64(time.Now().UnixNano()/1e6),
		Payload:  invokePayload,
		Sigs:     make([]*types.Sig, 0, 0),
	}
	return tx, nil
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

//SendRawTransaction send a transaction to mixbee network, and return hash of the transaction
func SendRawTransaction(tx *types.Transaction) (string, error) {
	var buffer bytes.Buffer
	err := tx.Serialize(&buffer)
	if err != nil {
		return "", fmt.Errorf("Serialize error:%s", err)
	}
	txData := hex.EncodeToString(buffer.Bytes())
	data, err := sendRpcRequest("sendrawtransaction", []interface{}{txData})
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

//GetSmartContractEvent return smart contract event execute by invoke transaction by hex string code
func GetSmartContractEvent(txHash string) (*rpccommon.ExecuteNotify, error) {
	data, err := sendRpcRequest("getsmartcodeevent", []interface{}{txHash})
	if err != nil {
		return nil, fmt.Errorf("sendRpcRequest error:%s", err)
	}
	notifies := &rpccommon.ExecuteNotify{}
	err = json.Unmarshal(data, &notifies)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal SmartContactEvent:%s error:%s", data, err)
	}
	return notifies, nil
}

func GetSmartContractEventInfo(txHash string) ([]byte, error) {
	return sendRpcRequest("getsmartcodeevent", []interface{}{txHash})
}

func GetRawTransaction(txHash string) ([]byte, error) {
	return sendRpcRequest("getrawtransaction", []interface{}{txHash, 1})
}

func GetBlock(hashOrHeight interface{}) ([]byte, error) {
	return sendRpcRequest("getblock", []interface{}{hashOrHeight, 1})
}

func GetBlockData(hashOrHeight interface{}) ([]byte, error) {
	data, err := sendRpcRequest("getblock", []interface{}{hashOrHeight})
	if err != nil {
		return nil, err
	}
	hexStr := ""
	err = json.Unmarshal(data, &hexStr)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal error:%s", err)
	}
	blockData, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, fmt.Errorf("hex.DecodeString error:%s", err)
	}
	return blockData, nil
}

func GetBlockCount() (uint32, error) {
	data, err := sendRpcRequest("getblockcount", []interface{}{})
	if err != nil {
		return 0, err
	}
	num := uint32(0)
	err = json.Unmarshal(data, &num)
	if err != nil {
		return 0, fmt.Errorf("json.Unmarshal:%s error:%s", data, err)
	}
	return num, nil
}

func DeployContract(
	gasPrice,
	gasLimit uint64,
	signer *account.Account,
	needStorage bool,
	code,
	cname,
	cversion,
	cauthor,
	cemail,
	cdesc string) (string, error) {

	c, err := hex.DecodeString(code)
	if err != nil {
		return "", fmt.Errorf("hex.DecodeString error:%s", err)
	}
	tx := NewDeployCodeTransaction(gasPrice, gasLimit, c, needStorage, cname, cversion, cauthor, cemail, cdesc)

	err = SignTransaction(signer, tx)
	if err != nil {
		return "", err
	}
	txHash, err := SendRawTransaction(tx)
	if err != nil {
		return "", fmt.Errorf("SendRawTransaction error:%s", err)
	}
	return txHash, nil
}

func PrepareDeployContract(
	needStorage bool,
	code,
	cname,
	cversion,
	cauthor,
	cemail,
	cdesc string) (*cstates.PreExecResult, error) {
	c, err := hex.DecodeString(code)
	if err != nil {
		return nil, fmt.Errorf("hex.DecodeString error:%s", err)
	}
	tx := NewDeployCodeTransaction(0, 0, c, needStorage, cname, cversion, cauthor, cemail, cdesc)
	var buffer bytes.Buffer
	err = tx.Serialize(&buffer)
	if err != nil {
		return nil, fmt.Errorf("Serialize error:%s", err)
	}
	txData := hex.EncodeToString(buffer.Bytes())
	data, err := sendRpcRequest("sendrawtransaction", []interface{}{txData, 1})
	if err != nil {
		return nil, err
	}
	preResult := &cstates.PreExecResult{}
	err = json.Unmarshal(data, &preResult)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal PreExecResult:%s error:%s", data, err)
	}
	return preResult, nil
}

func InvokeNativeContract(
	gasPrice,
	gasLimit uint64,
	signer *account.Account,
	contractAddress common.Address,
	version byte,
	method string,
	params []interface{},
) (string, error) {
	tx, err := httpcom.NewNativeInvokeTransaction(gasPrice, gasLimit, contractAddress, version, method, params)
	if err != nil {
		return "", err
	}
	return InvokeSmartContract(signer, tx)
}

//Invoke wasm smart contract
//methodName is wasm contract action name
//paramType  is Json or Raw format
//version should be greater than 0 (0 is reserved for test)
func InvokeWasmVMContract(
	gasPrice,
	gasLimit uint64,
	siger *account.Account,
	cversion byte, //version of contract
	contractAddress common.Address,
	method string,
	paramType wasmvm.ParamType,
	params []interface{}) (string, error) {

	invokeCode, err := BuildWasmVMInvokeCode(contractAddress, method, paramType, cversion, params)
	if err != nil {
		return "", err
	}
	tx, err := httpcom.NewSmartContractTransaction(gasPrice, gasLimit, invokeCode)
	if err != nil {
		return "", err
	}
	return InvokeSmartContract(siger, tx)
}

//Invoke neo vm smart contract. if isPreExec is true, the invoke will not really execute
func InvokeNeoVMContract(
	gasPrice,
	gasLimit uint64,
	signer *account.Account,
	smartcodeAddress common.Address,
	params []interface{}) (string, error) {
	tx, err := httpcom.NewNeovmInvokeTransaction(gasPrice, gasLimit, smartcodeAddress, params)
	if err != nil {
		return "", err
	}
	return InvokeSmartContract(signer, tx)
}

//InvokeSmartContract is low level method to invoke contact.
func InvokeSmartContract(signer *account.Account, tx *types.Transaction) (string, error) {
	err := SignTransaction(signer, tx)
	if err != nil {
		return "", fmt.Errorf("SignTransaction error:%s", err)
	}
	txHash, err := SendRawTransaction(tx)
	if err != nil {
		return "", fmt.Errorf("SendTransaction error:%s", err)
	}
	return txHash, nil
}

func PrepareInvokeNeoVMContract(
	contractAddress common.Address,
	params []interface{},
) (*cstates.PreExecResult, error) {
	tx, err := httpcom.NewNeovmInvokeTransaction(0, 0, contractAddress, params)
	if err != nil {
		return nil, err
	}
	var buffer bytes.Buffer
	err = tx.Serialize(&buffer)
	if err != nil {
		return nil, fmt.Errorf("Serialize error:%s", err)
	}
	txData := hex.EncodeToString(buffer.Bytes())
	data, err := sendRpcRequest("sendrawtransaction", []interface{}{txData, 1})
	if err != nil {
		return nil, err
	}
	preResult := &cstates.PreExecResult{}
	err = json.Unmarshal(data, &preResult)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal PreExecResult:%s error:%s", data, err)
	}
	return preResult, nil
}

func PrepareInvokeCodeNeoVMContract(code []byte) (*cstates.PreExecResult, error) {
	tx, err := httpcom.NewSmartContractTransaction(0, 0, code)
	if err != nil {
		return nil, err
	}
	var buffer bytes.Buffer
	err = tx.Serialize(&buffer)
	if err != nil {
		return nil, fmt.Errorf("Serialize error:%s", err)
	}
	txData := hex.EncodeToString(buffer.Bytes())
	data, err := sendRpcRequest("sendrawtransaction", []interface{}{txData, 1})
	if err != nil {
		return nil, err
	}
	preResult := &cstates.PreExecResult{}
	err = json.Unmarshal(data, &preResult)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal PreExecResult:%s error:%s", data, err)
	}
	return preResult, nil
}

func PrepareInvokeNativeContract(
	contractAddress common.Address,
	version byte,
	method string,
	params []interface{}) (*cstates.PreExecResult, error) {
	tx, err := httpcom.NewNativeInvokeTransaction(0, 0, contractAddress, version, method, params)
	if err != nil {
		return nil, err
	}
	var buffer bytes.Buffer
	err = tx.Serialize(&buffer)
	if err != nil {
		return nil, fmt.Errorf("Serialize error:%s", err)
	}
	txData := hex.EncodeToString(buffer.Bytes())
	data, err := sendRpcRequest("sendrawtransaction", []interface{}{txData, 1})
	if err != nil {
		return nil, err
	}
	preResult := &cstates.PreExecResult{}
	err = json.Unmarshal(data, &preResult)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal PreExecResult:%s error:%s", data, err)
	}
	return preResult, nil
}

//NewDeployCodeTransaction return a smart contract deploy transaction instance
func NewDeployCodeTransaction(gasPrice, gasLimit uint64, code []byte, needStorage bool,
	cname, cversion, cauthor, cemail, cdesc string) *types.Transaction {

	deployPayload := &payload.DeployCode{
		Code:        code,
		NeedStorage: needStorage,
		Name:        cname,
		Version:     cversion,
		Author:      cauthor,
		Email:       cemail,
		Description: cdesc,
	}
	tx := &types.Transaction{
		Version:  VERSION_TRANSACTION,
		TxType:   types.Deploy,
		Nonce:    uint64(time.Now().UnixNano()/1e6),
		Payload:  deployPayload,
		GasPrice: gasPrice,
		GasLimit: gasLimit,
		Sigs:     make([]*types.Sig, 0, 0),
	}
	return tx
}

//for wasm vm
//build param bytes for wasm contract
func buildWasmContractParam(params []interface{}, paramType wasmvm.ParamType) ([]byte, error) {
	switch paramType {
	case wasmvm.Json:
		args := make([]exec.Param, len(params))

		for i, param := range params {
			switch param.(type) {
			case string:
				arg := exec.Param{Ptype: "string", Pval: param.(string)}
				args[i] = arg
			case int:
				arg := exec.Param{Ptype: "int", Pval: strconv.Itoa(param.(int))}
				args[i] = arg
			case int64:
				arg := exec.Param{Ptype: "int64", Pval: strconv.FormatInt(param.(int64), 10)}
				args[i] = arg
			case []int:
				bf := bytes.NewBuffer(nil)
				array := param.([]int)
				for i, tmp := range array {
					bf.WriteString(strconv.Itoa(tmp))
					if i != len(array)-1 {
						bf.WriteString(",")
					}
				}
				arg := exec.Param{Ptype: "int_array", Pval: bf.String()}
				args[i] = arg
			case []int64:
				bf := bytes.NewBuffer(nil)
				array := param.([]int64)
				for i, tmp := range array {
					bf.WriteString(strconv.FormatInt(tmp, 10))
					if i != len(array)-1 {
						bf.WriteString(",")
					}
				}
				arg := exec.Param{Ptype: "int_array", Pval: bf.String()}
				args[i] = arg
			default:
				return nil, fmt.Errorf("not a supported type :%v\n", param)
			}
		}

		bs, err := json.Marshal(exec.Args{args})
		if err != nil {
			return nil, err
		}
		return bs, nil
	case wasmvm.Raw:
		bf := bytes.NewBuffer(nil)
		for _, param := range params {
			switch param.(type) {
			case string:
				tmp := bytes.NewBuffer(nil)
				serialization.WriteString(tmp, param.(string))
				bf.Write(tmp.Bytes())

			case int:
				tmpBytes := make([]byte, 4)
				binary.LittleEndian.PutUint32(tmpBytes, uint32(param.(int)))
				bf.Write(tmpBytes)

			case int64:
				tmpBytes := make([]byte, 8)
				binary.LittleEndian.PutUint64(tmpBytes, uint64(param.(int64)))
				bf.Write(tmpBytes)

			default:
				return nil, fmt.Errorf("not a supported type :%v\n", param)
			}
		}
		return bf.Bytes(), nil
	default:
		return nil, fmt.Errorf("unsupported type")
	}
}

//BuildWasmVMInvokeCode return wasn vm invoke code
func BuildWasmVMInvokeCode(smartcodeAddress common.Address, methodName string, paramType wasmvm.ParamType, version byte, params []interface{}) ([]byte, error) {
	contract := &cstates.Contract{}
	contract.Address = smartcodeAddress
	contract.Method = methodName
	contract.Version = version

	argbytes, err := buildWasmContractParam(params, paramType)

	if err != nil {
		return nil, fmt.Errorf("build wasm contract param failed:%s", err)
	}
	contract.Args = argbytes
	bf := bytes.NewBuffer(nil)
	contract.Serialize(bf)
	return bf.Bytes(), nil
}

//ParseNeoVMContractReturnTypeBool return bool value of smart contract execute code.
func ParseNeoVMContractReturnTypeBool(hexStr string) (bool, error) {
	return hexStr == "01", nil
}

//ParseNeoVMContractReturnTypeInteger return integer value of smart contract execute code.
func ParseNeoVMContractReturnTypeInteger(hexStr string) (int64, error) {
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return 0, fmt.Errorf("hex.DecodeString error:%s", err)
	}
	return common.BigIntFromNeoBytes(data).Int64(), nil
}

//ParseNeoVMContractReturnTypeByteArray return []byte value of smart contract execute code.
func ParseNeoVMContractReturnTypeByteArray(hexStr string) (string, error) {
	return hexStr, nil
}

//ParseNeoVMContractReturnTypeString return string value of smart contract execute code.
func ParseNeoVMContractReturnTypeString(hexStr string) (string, error) {
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return "", fmt.Errorf("hex.DecodeString:%s error:%s", hexStr, err)
	}
	return string(data), nil
}
