package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/mixbee/mixbee/cmd/utils"
	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/common/config"
	"github.com/mixbee/mixbee/common/log"
	"github.com/mixbee/mixbee/smartcontract/service/native/governance"
	"github.com/urfave/cli"
)

func SetMixbeeConfig(ctx *cli.Context) (*config.MixbeeConfig, error) {
	cfg := config.DefConfig
	netWorkId := ctx.GlobalInt(utils.GetFlagName(utils.NetworkIdFlag))
	switch netWorkId {
	case config.NETWORK_ID_MAIN_NET: // 主网，1
		cfg.Genesis = config.MainNetConfig
	case config.NETWORK_ID_POLARIS_NET:
		cfg.Genesis = config.PolarisConfig
	}

	err := setGenesis(ctx, cfg.Genesis)
	if err != nil {
		return nil, fmt.Errorf("setGenesis error:%s", err)
	}
	setCommonConfig(ctx, cfg.Common)
	setConsensusConfig(ctx, cfg.Consensus)
	setP2PNodeConfig(ctx, cfg.P2PNode)
	setRpcConfig(ctx, cfg.Rpc)
	setRestfulConfig(ctx, cfg.Restful)
	setWebSocketConfig(ctx, cfg.Ws)
	if cfg.Genesis.ConsensusType == config.CONSENSUS_TYPE_SOLO {
		cfg.Ws.EnableHttpWs = true
		cfg.Restful.EnableHttpRestful = true
		cfg.Consensus.EnableConsensus = true
		cfg.P2PNode.NetworkId = config.NETWORK_ID_SOLO_NET // solo模式，network 模式为3
		cfg.P2PNode.NetworkName = config.GetNetworkName(cfg.P2PNode.NetworkId)
		cfg.P2PNode.NetworkMagic = config.GetNetworkMagic(cfg.P2PNode.NetworkId)
	}
	if netWorkId > 0 {
		cfg.P2PNode.NetworkId = uint32(netWorkId)
	}

	err = setCrossChainConfig(ctx, cfg.CrossChain)
	if err != nil {
		return nil, fmt.Errorf("setCrossChain error:%s", err)
	}

	return cfg, nil
}

func setCrossChainConfig(ctx *cli.Context, verifyConfig *config.CrossChainVerifyConfig) error {

	if ctx.GlobalBool(utils.GetFlagName(utils.EnableCrossChainVerifyFlag)) == false && ctx.GlobalBool(utils.GetFlagName(utils.EnableCrossChainInteractiveFlag)) == false {
		return nil
	}

	verifyConfig.EnableCrossChainVerify = ctx.GlobalBool(utils.GetFlagName(utils.EnableCrossChainVerifyFlag))
	verifyConfig.EnableCrossChainInteractive = ctx.GlobalBool(utils.GetFlagName(utils.EnableCrossChainInteractiveFlag))
	if ctx.GlobalBool(utils.GetFlagName(utils.EnableCrossChainInteractiveFlag)) {
		node := ctx.GlobalStringSlice(utils.GetFlagName(utils.CrossChainVerifyNode))
		if len(node) == 0 {
			return fmt.Errorf("if set enablecrosschaininter is true,crosschainnode must be set ")
		}
		verifyConfig.MainVerifyNode = node
	}

	return nil
}

func setGenesis(ctx *cli.Context, cfg *config.GenesisConfig) error {
	if ctx.GlobalBool(utils.GetFlagName(utils.EnableTestModeFlag)) {
		// 如果为 testmode 模式，直接启动solo共识
		cfg.ConsensusType = config.CONSENSUS_TYPE_SOLO
		cfg.SOLO.GenBlockTime = ctx.Uint(utils.GetFlagName(utils.TestModeGenBlockTimeFlag))
		if cfg.SOLO.GenBlockTime <= 1 {
			cfg.SOLO.GenBlockTime = config.DEFAULT_GEN_BLOCK_TIME
		}
		return nil
	}

	if !ctx.IsSet(utils.GetFlagName(utils.ConfigFlag)) {
		return nil
	}

	genesisFile := ctx.GlobalString(utils.GetFlagName(utils.ConfigFlag))
	if !common.FileExisted(genesisFile) {
		return nil
	}

	log.Infof("Load genesis config:%s", genesisFile)
	data, err := ioutil.ReadFile(genesisFile)
	if err != nil {
		return fmt.Errorf("ioutil.ReadFile:%s error:%s", genesisFile, err)
	}
	// Remove the UTF-8 Byte Order Mark
	data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))

	cfg.Reset()
	err = json.Unmarshal(data, cfg)
	if err != nil {
		return fmt.Errorf("json.Unmarshal GenesisConfig:%s error:%s", data, err)
	}
	switch cfg.ConsensusType {
	case config.CONSENSUS_TYPE_DBFT: // dbft
		if len(cfg.DBFT.Bookkeepers) < config.DBFT_MIN_NODE_NUM {
			return fmt.Errorf("DBFT consensus at least need %d bookkeepers in config", config.DBFT_MIN_NODE_NUM)
		}
		if cfg.DBFT.GenBlockTime <= 0 {
			cfg.DBFT.GenBlockTime = config.DEFAULT_GEN_BLOCK_TIME
		}
	case config.CONSENSUS_TYPE_VBFT: // vbft
		err = governance.CheckVBFTConfig(cfg.VBFT)
		if err != nil {
			return fmt.Errorf("VBFT config error %v", err)
		}
		if len(cfg.VBFT.Peers) < config.VBFT_MIN_NODE_NUM {
			return fmt.Errorf("VBFT consensus at least need %d peers in config", config.VBFT_MIN_NODE_NUM)
		}
	default:
		return fmt.Errorf("Unknow consensus:%s", cfg.ConsensusType)
	}

	return nil
}

func setCommonConfig(ctx *cli.Context, cfg *config.CommonConfig) {
	cfg.LogLevel = ctx.GlobalUint(utils.GetFlagName(utils.LogLevelFlag))
	cfg.EnableEventLog = !ctx.GlobalBool(utils.GetFlagName(utils.DisableEventLogFlag))
	cfg.GasLimit = ctx.GlobalUint64(utils.GetFlagName(utils.GasLimitFlag))
	cfg.GasPrice = ctx.GlobalUint64(utils.GetFlagName(utils.GasPriceFlag))
	cfg.DataDir = ctx.GlobalString(utils.GetFlagName(utils.DataDirFlag))
}

func setConsensusConfig(ctx *cli.Context, cfg *config.ConsensusConfig) {
	cfg.EnableConsensus = ctx.GlobalBool(utils.GetFlagName(utils.EnableConsensusFlag))
	cfg.MaxTxInBlock = ctx.GlobalUint(utils.GetFlagName(utils.MaxTxInBlockFlag))
}

func setP2PNodeConfig(ctx *cli.Context, cfg *config.P2PNodeConfig) {
	cfg.NetworkId = uint32(ctx.GlobalUint(utils.GetFlagName(utils.NetworkIdFlag)))
	cfg.NetworkMagic = config.GetNetworkMagic(cfg.NetworkId)
	cfg.NetworkName = config.GetNetworkName(cfg.NetworkId)
	cfg.NodePort = ctx.GlobalUint(utils.GetFlagName(utils.NodePortFlag))
	cfg.NodeConsensusPort = ctx.GlobalUint(utils.GetFlagName(utils.ConsensusPortFlag))
	cfg.DualPortSupport = ctx.GlobalBool(utils.GetFlagName(utils.DualPortSupportFlag))
	cfg.ReservedPeersOnly = ctx.GlobalBool(utils.GetFlagName(utils.ReservedPeersOnlyFlag))
	cfg.MaxConnInBound = ctx.GlobalUint(utils.GetFlagName(utils.MaxConnInBoundFlag))
	cfg.MaxConnOutBound = ctx.GlobalUint(utils.GetFlagName(utils.MaxConnOutBoundFlag))
	cfg.MaxConnInBoundForSingleIP = ctx.GlobalUint(utils.GetFlagName(utils.MaxConnInBoundForSingleIPFlag))

	rsvfile := ctx.GlobalString(utils.GetFlagName(utils.ReservedPeersFileFlag))
	if cfg.ReservedPeersOnly {
		if !common.FileExisted(rsvfile) {
			log.Infof("file %s not exist\n", rsvfile)
			return
		}
		peers, err := ioutil.ReadFile(rsvfile)
		if err != nil {
			log.Errorf("ioutil.ReadFile:%s error:%s", rsvfile, err)
			return
		}
		peers = bytes.TrimPrefix(peers, []byte("\xef\xbb\xbf"))

		err = json.Unmarshal(peers, &cfg.ReservedCfg)
		if err != nil {
			log.Errorf("json.Unmarshal reserved peers:%s error:%s", peers, err)
			return
		}
		for i := 0; i < len(cfg.ReservedCfg.ReservedPeers); i++ {
			log.Info("reserved addr: " + cfg.ReservedCfg.ReservedPeers[i])
		}
		for i := 0; i < len(cfg.ReservedCfg.MaskPeers); i++ {
			log.Info("mask addr: " + cfg.ReservedCfg.MaskPeers[i])
		}
	}

}

func setRpcConfig(ctx *cli.Context, cfg *config.RpcConfig) {
	cfg.EnableHttpJsonRpc = !ctx.Bool(utils.GetFlagName(utils.RPCDisabledFlag))
	cfg.HttpJsonPort = ctx.GlobalUint(utils.GetFlagName(utils.RPCPortFlag))
	cfg.HttpLocalPort = ctx.GlobalUint(utils.GetFlagName(utils.RPCLocalProtFlag))
}

func setRestfulConfig(ctx *cli.Context, cfg *config.RestfulConfig) {
	cfg.EnableHttpRestful = ctx.GlobalBool(utils.GetFlagName(utils.RestfulEnableFlag))
	cfg.HttpRestPort = ctx.GlobalUint(utils.GetFlagName(utils.RestfulPortFlag))
}

func setWebSocketConfig(ctx *cli.Context, cfg *config.WebSocketConfig) {
	cfg.EnableHttpWs = ctx.GlobalBool(utils.GetFlagName(utils.WsEnabledFlag))
	cfg.HttpWsPort = ctx.GlobalUint(utils.GetFlagName(utils.WsPortFlag))
}

func SetRpcPort(ctx *cli.Context) {
	if ctx.IsSet(utils.GetFlagName(utils.RPCPortFlag)) {
		config.DefConfig.Rpc.HttpJsonPort = ctx.Uint(utils.GetFlagName(utils.RPCPortFlag))
	}
}
