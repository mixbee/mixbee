

package jsonrpc

import (
	"net/http"
	"strconv"
	"fmt"
	cfg "github.com/mixbee/mixbee/common/config"
	"github.com/mixbee/mixbee/common/log"
	"github.com/mixbee/mixbee/http/base/rpc"
)

func StartRPCServer() error {
	log.Debug()
	http.HandleFunc("/", rpc.Handle)

	rpc.HandleFunc("getgenerateblocktime", rpc.GetGenerateBlockTime)
	rpc.HandleFunc("getbestblockhash", rpc.GetBestBlockHash)
	rpc.HandleFunc("getblock", rpc.GetBlock)
	rpc.HandleFunc("getblockcount", rpc.GetBlockCount)
	rpc.HandleFunc("getblockhash", rpc.GetBlockHash)
	rpc.HandleFunc("getconnectioncount", rpc.GetConnectionCount)

	rpc.HandleFunc("getrawtransaction", rpc.GetRawTransaction)
	rpc.HandleFunc("sendrawtransaction", rpc.SendRawTransaction)
	rpc.HandleFunc("getstorage", rpc.GetStorage)
	rpc.HandleFunc("getversion", rpc.GetNodeVersion)

	rpc.HandleFunc("getcontractstate", rpc.GetContractState)
	rpc.HandleFunc("getmempooltxcount", rpc.GetMemPoolTxCount)
	rpc.HandleFunc("getmempooltxstate", rpc.GetMemPoolTxState)
	rpc.HandleFunc("getsmartcodeevent", rpc.GetSmartCodeEvent)
	rpc.HandleFunc("getblockheightbytxhash", rpc.GetBlockHeightByTxHash)

	rpc.HandleFunc("getbalance", rpc.GetBalance)
	rpc.HandleFunc("getallowance", rpc.GetAllowance)
	rpc.HandleFunc("getmerkleproof", rpc.GetMerkleProof)
	rpc.HandleFunc("getblocktxsbyheight", rpc.GetBlockTxsByHeight)
	rpc.HandleFunc("getgasprice", rpc.GetGasPrice)
	rpc.HandleFunc("getunboundmbg", rpc.GetUnboundMbg)

	//cross chain
	rpc.HandleFunc("registerSubChainNode", rpc.RegisterSubChainNode)
	rpc.HandleFunc("pushCrossChainTxInfo", rpc.PushCrossChainTxInfo)
	rpc.HandleFunc("getAllVerifyNodeInfo",rpc.GetAllCrossChainVerifyNodes)

	// 查询mix test key
	rpc.HandleFunc("getkey", rpc.GetKey)

	//查询 cross chain 信息
	rpc.HandleFunc("crossQuery", rpc.CrossChainQuery)
	rpc.HandleFunc("crossHistory", rpc.CrossChainHistory)

	err := http.ListenAndServe(":"+strconv.Itoa(int(cfg.DefConfig.Rpc.HttpJsonPort)), nil)
	if err != nil {
		return fmt.Errorf("ListenAndServe error:%s", err)
	}
	return nil
}
