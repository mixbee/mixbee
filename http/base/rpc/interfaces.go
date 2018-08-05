

package rpc

import (
	"bytes"
	"encoding/hex"
	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/common/config"
	"github.com/mixbee/mixbee/common/log"
	"github.com/mixbee/mixbee/core/payload"
	scom "github.com/mixbee/mixbee/core/store/common"
	"github.com/mixbee/mixbee/core/types"
	ontErrors "github.com/mixbee/mixbee/errors"
	bactor "github.com/mixbee/mixbee/http/base/actor"
	bcomn "github.com/mixbee/mixbee/http/base/common"
	berr "github.com/mixbee/mixbee/http/base/error"
	"github.com/mixbee/mixbee/smartcontract/service/native/utils"
	"strconv"
	"github.com/mixbee/mixbee/crosschain"
)

func GetGenerateBlockTime(params []interface{}) map[string]interface{} {
	var genBlockTime interface{}
	if config.DefConfig.Genesis.ConsensusType == config.CONSENSUS_TYPE_DBFT {
		genBlockTime = config.DefConfig.Genesis.DBFT.GenBlockTime
	} else if config.DefConfig.Genesis.ConsensusType == config.CONSENSUS_TYPE_SOLO {
		genBlockTime = config.DefConfig.Genesis.SOLO.GenBlockTime
	} else {
		genBlockTime = nil
	}
	return responseSuccess(genBlockTime)
}

func GetBestBlockHash(params []interface{}) map[string]interface{} {
	hash := bactor.CurrentBlockHash()
	return responseSuccess(hash.ToHexString())
}

// Input JSON string examples for getblock method as following:
//   {"jsonrpc": "2.0", "method": "getblock", "params": [1], "id": 0}
//   {"jsonrpc": "2.0", "method": "getblock", "params": ["aabbcc.."], "id": 0}
func GetBlock(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return responsePack(berr.INVALID_PARAMS, nil)
	}
	var err error
	var hash common.Uint256
	switch (params[0]).(type) {
	// block height
	case float64:
		index := uint32(params[0].(float64))
		hash = bactor.GetBlockHashFromStore(index)
		if hash == common.UINT256_EMPTY {
			return responsePack(berr.INVALID_PARAMS, "")
		}
		// block hash
	case string:
		str := params[0].(string)
		hash, err = common.Uint256FromHexString(str)
		if err != nil {
			return responsePack(berr.INVALID_PARAMS, "")
		}
	default:
		return responsePack(berr.INVALID_PARAMS, "")
	}
	block, err := bactor.GetBlockFromStore(hash)
	if err != nil {
		return responsePack(berr.UNKNOWN_BLOCK, "unknown block")
	}
	if len(params) >= 2 {
		switch (params[1]).(type) {
		case float64:
			json := uint32(params[1].(float64))
			if json == 1 {
				return responseSuccess(bcomn.GetBlockInfo(block))
			}
		default:
			return responsePack(berr.INVALID_PARAMS, "")
		}
	}
	w := bytes.NewBuffer(nil)
	block.Serialize(w)
	return responseSuccess(common.ToHexString(w.Bytes()))
}

func GetBlockCount(params []interface{}) map[string]interface{} {
	height := bactor.GetCurrentBlockHeight()
	return responseSuccess(height + 1)
}

// A JSON example for getblockhash method as following:
//   {"jsonrpc": "2.0", "method": "getblockhash", "params": [1], "id": 0}
func GetBlockHash(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return responsePack(berr.INVALID_PARAMS, nil)
	}
	switch params[0].(type) {
	case float64:
		height := uint32(params[0].(float64))
		hash := bactor.GetBlockHashFromStore(height)
		if hash == common.UINT256_EMPTY {
			return responsePack(berr.UNKNOWN_BLOCK, "")
		}
		return responseSuccess(hash.ToHexString())
	default:
		return responsePack(berr.INVALID_PARAMS, "")
	}
}

func GetConnectionCount(params []interface{}) map[string]interface{} {
	count, err := bactor.GetConnectionCnt()
	if err != nil {
		log.Errorf("GetConnectionCount error:%s", err)
		return responsePack(berr.INTERNAL_ERROR, false)
	}
	return responseSuccess(count)
}

func GetRawMemPool(params []interface{}) map[string]interface{} {
	txs := []*bcomn.Transactions{}
	txpool := bactor.GetTxsFromPool(false)
	for _, t := range txpool {
		txs = append(txs, bcomn.TransArryByteToHexString(t))
	}
	if len(txs) == 0 {
		return responsePack(berr.INVALID_PARAMS, nil)
	}
	return responseSuccess(txs)
}

func GetMemPoolTxCount(params []interface{}) map[string]interface{} {
	count, err := bactor.GetTxnCount()
	if err != nil {
		return responsePack(berr.INTERNAL_ERROR, nil)
	}
	return responseSuccess(count)
}

func GetMemPoolTxState(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return responsePack(berr.INVALID_PARAMS, nil)
	}
	switch params[0].(type) {
	case string:
		str := params[0].(string)
		hash, err := common.Uint256FromHexString(str)
		if err != nil {
			return responsePack(berr.INVALID_PARAMS, "")
		}
		txEntry, err := bactor.GetTxFromPool(hash)
		if err != nil {
			return responsePack(berr.UNKNOWN_TRANSACTION, "unknown transaction")
		}
		attrs := []bcomn.TXNAttrInfo{}
		for _, t := range txEntry.Attrs {
			attrs = append(attrs, bcomn.TXNAttrInfo{t.Height, int(t.Type), int(t.ErrCode)})
		}
		info := bcomn.TXNEntryInfo{attrs}
		return responseSuccess(info)
	default:
		return responsePack(berr.INVALID_PARAMS, "")
	}
}

// A JSON example for getrawtransaction method as following:
//   {"jsonrpc": "2.0", "method": "getrawtransaction", "params": ["transactioin hash in hex"], "id": 0}
func GetRawTransaction(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return responsePack(berr.INVALID_PARAMS, nil)
	}
	var tx *types.Transaction
	switch params[0].(type) {
	case string:
		str := params[0].(string)
		hash, err := common.Uint256FromHexString(str)
		if err != nil {
			return responsePack(berr.INVALID_PARAMS, "")
		}
		t, err := bactor.GetTransaction(hash)
		if err != nil {
			return responsePack(berr.UNKNOWN_TRANSACTION, "unknown transaction")
		}
		tx = t
	default:
		return responsePack(berr.INVALID_PARAMS, "")
	}

	if len(params) >= 2 {
		switch (params[1]).(type) {
		case float64:
			json := uint32(params[1].(float64))
			if json == 1 {
				return responseSuccess(bcomn.TransArryByteToHexString(tx))
			}
		default:
			return responsePack(berr.INVALID_PARAMS, "")
		}
	}
	w := bytes.NewBuffer(nil)
	tx.Serialize(w)
	return responseSuccess(common.ToHexString(w.Bytes()))
}

//   {"jsonrpc": "2.0", "method": "getstorage", "params": ["code hash", "key"], "id": 0}
func GetStorage(params []interface{}) map[string]interface{} {
	if len(params) < 2 {
		return responsePack(berr.INVALID_PARAMS, nil)
	}

	var address common.Address
	var key []byte
	switch params[0].(type) {
	case string:
		str := params[0].(string)
		var err error
		if len(str) == common.ADDR_LEN*2 {
			address, err = common.AddressFromHexString(str)
		} else {
			address, err = common.AddressFromBase58(str)
		}
		if err != nil {
			return responsePack(berr.INVALID_PARAMS, "")
		}
	default:
		return responsePack(berr.INVALID_PARAMS, "")
	}

	switch params[1].(type) {
	case string:
		str := params[1].(string)
		hex, err := hex.DecodeString(str)
		if err != nil {
			return responsePack(berr.INVALID_PARAMS, "")
		}
		key = hex
	default:
		return responsePack(berr.INVALID_PARAMS, "")
	}
	value, err := bactor.GetStorageItem(address, key)
	if err != nil {
		if err == scom.ErrNotFound {
			return responseSuccess(nil)
		}
		return responsePack(berr.INVALID_PARAMS, "")
	}
	return responseSuccess(common.ToHexString(value))
}

// A JSON example for sendrawtransaction method as following:
//   {"jsonrpc": "2.0", "method": "sendrawtransaction", "params": ["raw transactioin in hex"], "id": 0}
func SendRawTransaction(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return responsePack(berr.INVALID_PARAMS, nil)
	}
	var hash common.Uint256
	switch params[0].(type) {
	case string:
		str := params[0].(string)
		hex, err := common.HexToBytes(str)
		if err != nil {
			return responsePack(berr.INVALID_PARAMS, "")
		}
		var txn types.Transaction
		if err := txn.Deserialize(bytes.NewReader(hex)); err != nil {
			return responsePack(berr.INVALID_TRANSACTION, "")
		}
		if len(params) > 1 {
			preExec, ok := params[1].(float64)
			if txn.TxType == types.Invoke || txn.TxType == types.Deploy {
				if ok && preExec == 1 {
					result, err := bactor.PreExecuteContract(&txn)
					if err != nil {
						log.Infof("PreExec: ", err)
						return responsePack(berr.SMARTCODE_ERROR, "")
					}
					return responseSuccess(result)
				}
			}
		}
		hash = txn.Hash()
		if errCode := bcomn.VerifyAndSendTx(&txn); errCode != ontErrors.ErrNoError {
			return responseSuccess(errCode.Error())
		}
	default:
		return responsePack(berr.INVALID_PARAMS, "")
	}
	return responseSuccess(hash.ToHexString())
}

func GetNodeVersion(params []interface{}) map[string]interface{} {
	return responseSuccess(config.Version)
}

func GetContractState(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return responsePack(berr.INVALID_PARAMS, nil)
	}
	var contract *payload.DeployCode
	switch params[0].(type) {
	case string:
		str := params[0].(string)
		var address common.Address
		var err error
		if len(str) == (common.ADDR_LEN * 2) {
			address, err = common.AddressFromHexString(str)
		} else {
			address, err = common.AddressFromBase58(str)
		}
		if err != nil {
			return responsePack(berr.INVALID_PARAMS, "")
		}
		c, err := bactor.GetContractStateFromStore(address)
		if err != nil {
			return responsePack(berr.UNKNOWN_CONTRACT, "unknow contract")
		}
		contract = c
	default:
		return responsePack(berr.INVALID_PARAMS, "")
	}
	if len(params) >= 2 {
		switch (params[1]).(type) {
		case float64:
			json := uint32(params[1].(float64))
			if json == 1 {
				return responseSuccess(bcomn.TransPayloadToHex(contract))
			}
		default:
			return responsePack(berr.INVALID_PARAMS, "")
		}
	}
	w := bytes.NewBuffer(nil)
	contract.Serialize(w)
	return responseSuccess(common.ToHexString(w.Bytes()))
}

func GetSmartCodeEvent(params []interface{}) map[string]interface{} {
	if !config.DefConfig.Common.EnableEventLog {
		return responsePack(berr.INVALID_METHOD, "")
	}
	if len(params) < 1 {
		return responsePack(berr.INVALID_PARAMS, nil)
	}

	switch (params[0]).(type) {
	// block height
	case float64:
		height := uint32(params[0].(float64))
		eventInfos, err := bactor.GetEventNotifyByHeight(height)
		if err != nil {
			if err == scom.ErrNotFound {
				return responseSuccess(nil)
			}
			return responsePack(berr.INTERNAL_ERROR, "")
		}
		eInfos := make([]*bcomn.ExecuteNotify, 0, len(eventInfos))
		for _, eventInfo := range eventInfos {
			_, notify := bcomn.GetExecuteNotify(eventInfo)
			eInfos = append(eInfos, &notify)
		}
		return responseSuccess(eInfos)
		//txhash
	case string:
		str := params[0].(string)
		hash, err := common.Uint256FromHexString(str)
		if err != nil {
			return responsePack(berr.INVALID_PARAMS, "")
		}
		eventInfo, err := bactor.GetEventNotifyByTxHash(hash)
		if err != nil {
			if scom.ErrNotFound == err {
				return responseSuccess(nil)
			}
			return responsePack(berr.INTERNAL_ERROR, "")
		}
		_, notify := bcomn.GetExecuteNotify(eventInfo)
		return responseSuccess(notify)
	default:
		return responsePack(berr.INVALID_PARAMS, "")
	}
	return responsePack(berr.INVALID_PARAMS, "")
}

func GetBlockHeightByTxHash(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return responsePack(berr.INVALID_PARAMS, nil)
	}

	switch (params[0]).(type) {
	// tx hash
	case string:
		str := params[0].(string)
		hash, err := common.Uint256FromHexString(str)
		if err != nil {
			return responsePack(berr.INVALID_PARAMS, "")
		}
		height, _, err := bactor.GetTxnWithHeightByTxHash(hash)
		if err != nil {
			return responsePack(berr.INVALID_PARAMS, "")
		}
		return responseSuccess(height)
	default:
		return responsePack(berr.INVALID_PARAMS, "")
	}
	return responsePack(berr.INVALID_PARAMS, "")
}

func GetBalance(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return responsePack(berr.INVALID_PARAMS, "")
	}
	addrBase58, ok := params[0].(string)
	if !ok {
		return responsePack(berr.INVALID_PARAMS, "")
	}
	address, err := common.AddressFromBase58(addrBase58)
	if err != nil {
		return responsePack(berr.INVALID_PARAMS, "")
	}
	rsp, err := bcomn.GetBalance(address)
	if err != nil {
		return responsePack(berr.INVALID_PARAMS, "")
	}
	return responseSuccess(rsp)
}

func GetAllowance(params []interface{}) map[string]interface{} {
	if len(params) < 3 {
		return responsePack(berr.INVALID_PARAMS, "")
	}
	asset, ok := params[0].(string)
	if !ok {
		return responsePack(berr.INVALID_PARAMS, "")
	}
	fromAddrStr, ok := params[1].(string)
	if !ok {
		return responsePack(berr.INVALID_PARAMS, "")
	}
	fromAddr, err := common.AddressFromBase58(fromAddrStr)
	if err != nil {
		return responsePack(berr.INVALID_PARAMS, "")
	}
	toAddrStr, ok := params[2].(string)
	if !ok {
		return responsePack(berr.INVALID_PARAMS, "")
	}
	toAddr, err := common.AddressFromBase58(toAddrStr)
	if err != nil {
		return responsePack(berr.INVALID_PARAMS, "")
	}
	rsp, err := bcomn.GetAllowance(asset, fromAddr, toAddr)
	if err != nil {
		return responsePack(berr.INVALID_PARAMS, "")
	}
	return responseSuccess(rsp)
}

func GetMerkleProof(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return responsePack(berr.INVALID_PARAMS, "")
	}
	str, ok := params[0].(string)
	if !ok {
		return responsePack(berr.INVALID_PARAMS, "")
	}
	hash, err := common.Uint256FromHexString(str)
	if err != nil {
		return responsePack(berr.INVALID_PARAMS, "")
	}
	height, _, err := bactor.GetTxnWithHeightByTxHash(hash)
	if err != nil {
		return responsePack(berr.INVALID_PARAMS, "")
	}
	header, err := bactor.GetHeaderByHeight(height)
	if err != nil {
		return responsePack(berr.INVALID_PARAMS, "")
	}

	curHeight := bactor.GetCurrentBlockHeight()
	curHeader, err := bactor.GetHeaderByHeight(curHeight)
	if err != nil {
		return responsePack(berr.INVALID_PARAMS, "")
	}
	proof, err := bactor.GetMerkleProof(uint32(height), uint32(curHeight))
	if err != nil {
		return responsePack(berr.INTERNAL_ERROR, "")
	}
	var hashes []string
	for _, v := range proof {
		hashes = append(hashes, v.ToHexString())
	}
	return responseSuccess(bcomn.MerkleProof{"MerkleProof", header.TransactionsRoot.ToHexString(), height,
		curHeader.BlockRoot.ToHexString(), curHeight, hashes})
}

func GetBlockTxsByHeight(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return responsePack(berr.INVALID_PARAMS, nil)
	}
	switch params[0].(type) {
	case float64:
		height := uint32(params[0].(float64))
		hash := bactor.GetBlockHashFromStore(height)
		if hash == common.UINT256_EMPTY {
			return responsePack(berr.INVALID_PARAMS, "")
		}
		block, err := bactor.GetBlockFromStore(hash)
		if err != nil {
			return responsePack(berr.UNKNOWN_BLOCK, "")
		}
		return responseSuccess(bcomn.GetBlockTransactions(block))
	default:
		return responsePack(berr.INVALID_PARAMS, "")
	}
}

func GetGasPrice(params []interface{}) map[string]interface{} {
	result, err := bcomn.GetGasPrice()
	if err != nil {
		return responsePack(berr.INTERNAL_ERROR, "")
	}
	return responseSuccess(result)
}

func GetUnboundOng(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return responsePack(berr.INVALID_PARAMS, "")
	}
	str, ok := params[0].(string)
	if !ok {
		return responsePack(berr.INVALID_PARAMS, "")
	}
	toAddr, err := common.AddressFromBase58(str)
	if err != nil {
		return responsePack(berr.INVALID_PARAMS, "")
	}
	fromAddr := utils.OntContractAddress
	rsp, err := bcomn.GetAllowance("ong", fromAddr, toAddr)
	if err != nil {
		return responsePack(berr.INVALID_PARAMS, "")
	}
	return responseSuccess(rsp)
}

// A JSON example for registerSubChainNode method as following:
//   {"jsonrpc": "2.0", "method": "registerSubChainNode", "params": ["netWorkId","host:port"], "id": 0}
func RegisterSubChainNode(params []interface{}) map[string]interface{} {

	if !config.DefConfig.CrossChain.EnableCrossChainVerify {
		return responsePack(berr.INVALID_PARAMS, "this node not support cross chain verify")
	}
	if len(params) < 2 {
		return responsePack(berr.INVALID_PARAMS, "")
	}
	netWorkId, ok := params[0].(string)
	if !ok {
		return responsePack(berr.INVALID_PARAMS, "")
	}
	id,err := strconv.Atoi(netWorkId)
	if err != nil {
		responsePack(berr.INVALID_PARAMS,"networkId is a uint32")
	}
	nid := uint32(id)

	host, ok := params[1].(string)
	if !ok {
		return responsePack(berr.INVALID_PARAMS, "")
	}
	log.Infof("RegisterSubChainNode nid=%d,path=%s",nid,host)
	info,ok := config.DefConfig.CrossChain.SubChainNode[nid]
	if !ok {
			info = []string{}
			config.DefConfig.CrossChain.SubChainNode[nid] = info
	}
	info = append(info,host)
	config.DefConfig.CrossChain.SubChainNode[nid] = info

	return responseSuccess("success")
}

// A JSON example for pushCrossChainTxInfo method as following:
//   {"jsonrpc": "2.0", "method": "pushCrossChainTxInfo", "params": ["addrA","addrB","aAmount","bAmount","aNetId","bNetId","txHash","seqId",timestamp,nonce], "id": 0}
func PushCrossChainTxInfo(params []interface{}) map[string]interface{} {

	if !config.DefConfig.CrossChain.EnableCrossChainVerify {
		return responsePack(berr.INVALID_PARAMS, "this node not support cross chain verify")
	}

	if len(params) < 10 {
		return responsePack(berr.INVALID_PARAMS, "")
	}
	log.Infof("PushCrossChainTxInfo %#v\n",params)
	err := crosschain.CtxServer.PushCtxToPool(params)
	if err != nil {
		return responsePack(berr.INVALID_PARAMS, err.Error())
	}

	return responseSuccess("success")
}
